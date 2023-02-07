package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNoticeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNoticeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNoticeDataLogic {
	return &GetUserNoticeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserNoticeData 获取用户通知数据
func (l *GetUserNoticeDataLogic) GetUserNoticeData(in *pb.GetUserNoticeDataReq) (*pb.GetUserNoticeDataResp, error) {
	if in.UserId == "" {
		// 查询会话的订阅者
		subscribers, err := l.svcCtx.MsgService().GetConvSubscribers(l.ctx, &pb.GetConvSubscribersReq{
			CommonReq:      in.CommonReq,
			ConvId:         in.ConvId,
			LastActiveTime: nil,
		})
		if err != nil {
			l.Errorf("GetUserNoticeData failed, err: %v", err)
			return nil, err
		}
		for _, userId := range subscribers.UserIdList {
			resp, err := l.getUserNoticeData(&pb.GetUserNoticeDataReq{
				CommonReq: in.GetCommonReq(),
				UserId:    userId,
				ConvId:    in.ConvId,
				DeviceId:  in.DeviceId,
			})
			if err != nil {
				l.Errorf("GetUserNoticeData failed, err: %v", err)
				return resp, err
			}
		}
		return &pb.GetUserNoticeDataResp{}, nil
	} else {
		return l.getUserNoticeData(in)
	}
}

func (l *GetUserNoticeDataLogic) getUserNoticeData(in *pb.GetUserNoticeDataReq) (*pb.GetUserNoticeDataResp, error) {
	var maxSeq int64
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetMaxConvAutoId", func(ctx context.Context) {
		maxSeq, err = noticemodel.GetMaxConvAutoId(l.ctx, l.svcCtx.Mysql(), in.ConvId, 0)
	})
	if in.DeviceId == nil {
		// 查询用户当前在线的设备
		userConn, err := l.svcCtx.ImService().GetUserConn(l.ctx, &pb.GetUserConnReq{
			UserIds: []string{in.UserId},
		})
		if err != nil {
			l.Errorf("GetUserNoticeData failed, err: %v", err)
			return nil, err
		}
		var deviceIds []string
		for _, conn := range userConn.ConnParams {
			deviceIds = append(deviceIds, conn.DeviceId)
		}
		if len(deviceIds) == 0 {
			return &pb.GetUserNoticeDataResp{}, nil
		}
		for _, deviceId := range deviceIds {
			err = l.getDeviceNoticeData(in, maxSeq, deviceId)
			if err != nil {
				l.Errorf("GetUserNoticeData failed, err: %v", err)
				return nil, err
			}
		}
	} else {
		err = l.getDeviceNoticeData(in, maxSeq, *in.DeviceId)
		if err != nil {
			l.Errorf("GetUserNoticeData failed, err: %v", err)
			return &pb.GetUserNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.GetUserNoticeDataResp{}, nil
}

func (l *GetUserNoticeDataLogic) getDeviceNoticeData(in *pb.GetUserNoticeDataReq, maxSeq int64, deviceId string) error {
	var minSeq int64
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetMinConvAutoId", func(ctx context.Context) {
		minSeq, err = noticemodel.GetMinConvAutoId(l.ctx, l.svcCtx.Mysql(), in.ConvId, in.UserId, deviceId)
	})
	if err != nil {
		l.Errorf("GetUserNoticeData failed, err: %v", err)
		return err
	}
	if minSeq >= maxSeq {
		return nil
	}
	notice, err := noticemodel.GetNotice(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.ConvId, in.UserId, deviceId, minSeq, maxSeq)
	if err != nil {
		l.Errorf("GetUserNoticeData failed, err: %v", err)
		return err
	}
	if notice == nil {
		// 查询广播的通知
		notice, err = noticemodel.GetNotice(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.ConvId, "", deviceId, minSeq, maxSeq)
		if err != nil {
			l.Errorf("GetUserNoticeData failed, err: %v", err)
			return err
		}
		if notice != nil {
			// 推送通知
			err = l.pushNotice(in, notice, deviceId)
			if err != nil {
				l.Errorf("GetUserNoticeData failed, err: %v", err)
				return err
			}
		}
	}
	return nil
}

func (l *GetUserNoticeDataLogic) pushNotice(in *pb.GetUserNoticeDataReq, notice *noticemodel.Notice, deviceId string) error {
	noticeData := &pb.NoticeData{
		ConvId:      notice.ConvId,
		NoticeId:    notice.NoticeId,
		CreateTime:  strconv.FormatInt(notice.CreateTime, 10),
		Title:       notice.Title,
		ContentType: notice.ContentType,
		Content:     notice.Content,
		Options: &pb.NoticeData_Options{
			StorageForClient: notice.Options.StorageForClient,
			UpdateConvNotice: notice.Options.UpdateConvNotice,
		},
		Ext: notice.Ext,
	}
	data, _ := proto.Marshal(noticeData)
	_, err := l.svcCtx.ImService().SendMsg(l.ctx, &pb.SendMsgReq{
		GetUserConnReq: &pb.GetUserConnReq{
			UserIds: []string{in.UserId},
			Devices: []string{deviceId},
		},
		Event: pb.PushEvent_PushNoticeData,
		Data:  data,
	})
	if err != nil {
		l.Errorf("GetUserNoticeData failed, err: %v", err)
		return err
	}
	return nil
}
