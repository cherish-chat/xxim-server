package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
	// 获取用户在线状态
	conn, err := l.svcCtx.ImService().GetUserConn(l.ctx, &pb.GetUserConnReq{
		UserIds: []string{in.UserId},
	})
	if err != nil {
		l.Errorf("get user latest conn error: %v", err)
		return &pb.GetUserNoticeDataResp{}, err
	}
	if len(conn.ConnParams) == 0 {
		// 用户不在线
		return &pb.GetUserNoticeDataResp{}, nil
	}
	var deviceIds []string
	for _, v := range conn.ConnParams {
		deviceIds = append(deviceIds, v.DeviceId)
	}
	// 获取用户订阅的通知
	var getUserNoticeConvIdsResp *pb.GetUserNoticeConvIdsResp
	xtrace.StartFuncSpan(l.ctx, "GetUserNoticeConvIds", func(ctx context.Context) {
		getUserNoticeConvIdsResp, err = NewGetUserNoticeConvIdsLogic(ctx, l.svcCtx).GetUserNoticeConvIds(&pb.GetUserNoticeConvIdsReq{UserId: in.UserId})
	})
	if err != nil {
		l.Errorf("get user notice convIds error: %v", err)
		return &pb.GetUserNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(getUserNoticeConvIdsResp.ConvIds) == 0 {
		return &pb.GetUserNoticeDataResp{}, nil
	}
	// 获取用户设备所有会话最终ack的消息的时间
	// hmget key convId1 convId2 ...
	var fs []func() error
	for _, deviceId := range deviceIds {
		fs = append(fs, func() error {
			values, err := l.svcCtx.Redis().HmgetCtx(l.ctx, rediskey.UserAckRecord(in.UserId, deviceId), getUserNoticeConvIdsResp.ConvIds...)
			if err != nil {
				l.Errorf("hmget user ack record error: %v", err)
				return err
			}
			dest := &noticemodel.Notice{}
			tx := l.svcCtx.Mysql().Model(dest).
				Where("((userId = ? OR isBroadcast = ?) AND pushInvalid = ?)", in.UserId, true, false)
			orBuilder := ""
			args := make([]interface{}, 0)
			for i, convId := range getUserNoticeConvIdsResp.ConvIds {
				if values[i] == "" {
					// 设置为当前时间
					values[i] = utils.AnyToString(time.Now().UnixMilli())
					err = l.svcCtx.Redis().HsetCtx(l.ctx, rediskey.UserAckRecord(in.UserId, deviceId), convId, values[i])
					if err != nil {
						l.Errorf("hset user ack record error: %v", err)
						return err
					}
				}
				orBuilder += "(convId = ? AND createTime > ?) OR "
				args = append(args, convId, utils.AnyToInt64(values[i]))
			}
			orBuilder = orBuilder[:len(orBuilder)-4]
			tx = tx.Where("("+orBuilder+")", args...)
			err = tx.Order("createTime ASC").First(dest).Error
			if err != nil {
				if xorm.RecordNotFound(err) {
					return nil
				}
				l.Errorf("get user notice data error: %v", err)
				return err
			}

			noticeDataListBytes, _ := proto.Marshal(&pb.NoticeDataList{NoticeDataList: []*pb.NoticeData{dest.ToProto()}})
			err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
				if dest.Options.OnlinePushOnce {
					// 设置为失效
					dest.PushInvalid = true
					err := tx.Model(dest).Where("noticeId = ? AND convId = ?", dest.NoticeId, dest.ConvId).Update("pushInvalid", true).Error
					if err != nil {
						l.Errorf("update notice pushInvalid error: %v", err)
						return err
					}
				}
				return nil
			}, func(*gorm.DB) error {
				resp, err := l.svcCtx.ImService().SendMsg(l.ctx, &pb.SendMsgReq{
					GetUserConnReq: &pb.GetUserConnReq{
						UserIds: []string{in.UserId},
						Devices: []string{deviceId},
					},
					Event: pb.PushEvent_PushNoticeDataList,
					Data:  noticeDataListBytes,
				})
				if err != nil {
					l.Errorf("PushNoticeData SendMsg err: %v", err)
				} else {
					l.Infof("PushNoticeData SendMsg resp: %v", utils.AnyToString(resp))
				}
				return err
			})
			if err != nil {
				l.Errorf("PushNoticeData SendMsg error: %v", err)
				return err
			}
			return nil
		})
	}
	err = mr.Finish(fs...)
	if err != nil {
		return &pb.GetUserNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetUserNoticeDataResp{}, nil
}
