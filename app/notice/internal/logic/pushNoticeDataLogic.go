package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/propagation"
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
	err := l.svcCtx.Mysql().Model(notice).Where("noticeId = ? AND convId = ?", in.NoticeId, in.ConvId).First(notice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.PushNoticeDataResp{}, nil
		}
		l.Errorf("find notice error: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if notice.IsBroadcast {
		var resp *pb.PushNoticeDataResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "pushBroadcastNoticeData", func(ctx context.Context) {
			resp, err = l.pushBroadcastNoticeData(in, notice)
		})
		return resp, err
	} else if notice.UserId != "" {
		var resp *pb.PushNoticeDataResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "pushUserNoticeData", func(ctx context.Context) {
			resp, err = l.pushUserNoticeData(in, notice, notice.UserId)
		})
		return resp, err
	}
	// 删除垃圾数据
	l.svcCtx.Mysql().Model(notice).Where("noticeId = ? AND convId = ?", in.NoticeId, in.ConvId).Delete(notice)
	return &pb.PushNoticeDataResp{}, nil
}

func (l *PushNoticeDataLogic) pushBroadcastNoticeData(in *pb.PushNoticeDataReq, notice *noticemodel.Notice) (*pb.PushNoticeDataResp, error) {
	var (
		getNoticeConvAllSubscribersResp *pb.GetNoticeConvAllSubscribersResp
		err                             error
	)
	xtrace.StartFuncSpan(l.ctx, "GetNoticeConvAllSubscribers", func(ctx context.Context) {
		getNoticeConvAllSubscribersResp, err = NewGetNoticeConvAllSubscribersLogic(ctx, l.svcCtx).GetNoticeConvAllSubscribers(&pb.GetNoticeConvAllSubscribersReq{
			CommonReq: in.CommonReq,
			ConvId:    notice.ConvId,
		})
	})
	if err != nil {
		l.Errorf("pushBroadcastNoticeData GetNoticeConvAllSubscribers err: %v", err)
		return &pb.PushNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(getNoticeConvAllSubscribersResp.UserIds) == 0 {
		l.Infof("pushBroadcastNoticeData no user need push")
		return &pb.PushNoticeDataResp{}, nil
	}
	var fs []func() error
	for _, userId := range getNoticeConvAllSubscribersResp.UserIds {
		copyUserId := userId
		fs = append(fs, func() error {
			_, err := l.pushUserNoticeData(in, notice, copyUserId)
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
	}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
		"userId": userId,
	}))
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
