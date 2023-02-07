package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckNoticeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAckNoticeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckNoticeDataLogic {
	return &AckNoticeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AckNoticeData 确认通知数据
func (l *AckNoticeDataLogic) AckNoticeData(in *pb.AckNoticeDataReq) (*pb.AckNoticeDataResp, error) {
	defer func() {
		xtrace.StartFuncSpan(l.ctx, "GetUserNoticeData", func(ctx context.Context) {
			_, err := NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
				CommonReq: in.GetCommonReq(),
				UserId:    in.CommonReq.UserId,
				ConvId:    in.ConvId,
				DeviceId:  utils.AnyPtr(in.CommonReq.DeviceId),
			})
			if err != nil {
				l.Errorf("GetUserNoticeData failed, err: %v", err)
			}
		})
	}()
	_, seq, _ := pb.ParseServerNoticeId(in.NoticeId)
	ackRecord := &noticemodel.NoticeAckRecord{
		ConvId:     in.ConvId,
		UserId:     in.CommonReq.UserId,
		DeviceId:   in.CommonReq.DeviceId,
		ConvAutoId: seq,
	}
	err := l.svcCtx.Mysql().Model(ackRecord).
		Where("convId = ? and userId = ? and deviceId = ?",
			in.ConvId, in.CommonReq.UserId, in.CommonReq.DeviceId).
		Updates(map[string]interface{}{
			"convAutoId": seq,
		}).Error
	if err != nil {
		l.Errorf("AckNoticeData failed, err: %v", err)
		return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	err = noticemodel.DelNoticeZSet(l.ctx, l.svcCtx.Redis(), in.ConvId, in.CommonReq.UserId, in.CommonReq.DeviceId, seq)
	if err != nil {
		l.Errorf("AckNoticeData failed, err: %v", err)
		return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AckNoticeDataResp{}, nil
}
