package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushNoticeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushNoticeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushNoticeDataLogic {
	return &PushNoticeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushNoticeData 推送通知数据
func (l *PushNoticeDataLogic) PushNoticeData(in *pb.PushNoticeDataReq) (*pb.PushNoticeDataResp, error) {
	notice := &noticemodel.Notice{}
	err := l.svcCtx.Mysql().Model(notice).Where("noticeId = ?", in.NoticeId).First(notice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.PushNoticeDataResp{}, nil
		}
		l.Errorf("find notice error: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if notice.IsBroadcast {
		return l.pushBroadcastNoticeData(in, notice)
	} else if notice.UserId != "" {
		return l.pushUserNoticeData(in, notice, notice.UserId)
	}
	// 删除垃圾数据
	l.svcCtx.Mysql().Model(notice).Where("noticeId = ?", in.NoticeId).Delete(notice)
	return &pb.PushNoticeDataResp{}, nil
}

func (l *PushNoticeDataLogic) pushBroadcastNoticeData(in *pb.PushNoticeDataReq, notice *noticemodel.Notice) (*pb.PushNoticeDataResp, error) {
	var (
		getNoticeConvAllSubscribersResp *pb.GetNoticeConvAllSubscribersResp
		err                             error
	)
	xtrace.StartFuncSpan(l.ctx, "", func(ctx context.Context) {
		getNoticeConvAllSubscribersResp, err = NewGetNoticeConvAllSubscribersLogic(ctx, l.svcCtx).GetNoticeConvAllSubscribers(&pb.GetNoticeConvAllSubscribersReq{
			CommonReq: in.CommonReq,
			ConvId:    notice.ConvId,
		})
	})
	if err != nil {
		l.Errorf("pushBroadcastNoticeData GetNoticeConvAllSubscribers err: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var tmpfs [8][]func() error
	for i, userId := range getNoticeConvAllSubscribersResp.UserIds {
		tmpfs[i%8] = append(tmpfs[i%8], func() error {
			_, err := l.pushUserNoticeData(in, notice, userId)
			return err
		})
	}
	var fs []func() error
	for _, tmpf := range tmpfs {
		fs = append(fs, func() error {
			for _, f := range tmpf {
				if err := f(); err != nil {
					return err
				}
			}
			return err
		})
	}
	err = mr.Finish(fs...)
	if err != nil {
		l.Errorf("pushBroadcastNoticeData err: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.PushNoticeDataResp{}, nil
}

func (l *PushNoticeDataLogic) pushUserNoticeData(in *pb.PushNoticeDataReq, notice *noticemodel.Notice, userId string) (*pb.PushNoticeDataResp, error) {
	var resp *pb.GetUserNoticeDataResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetUserNoticeData", func(ctx context.Context) {
		resp, err = NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
			CommonReq: in.CommonReq,
			UserId:    userId,
		})
	})
	if err != nil {
		l.Errorf("pushUserNoticeData GetUserNoticeData err: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if resp.CommonResp != nil && resp.CommonResp.Failed() {
		l.Errorf("pushUserNoticeData GetUserNoticeData resp: %v", utils.AnyToString(resp))
		return &pb.PushNoticeDataResp{CommonResp: resp.CommonResp}, nil
	}
	return &pb.PushNoticeDataResp{}, nil
}
