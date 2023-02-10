package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeepAliveLogic {
	return &KeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// KeepAlive
func (l *KeepAliveLogic) KeepAlive(in *pb.KeepAliveReq) (*pb.KeepAliveResp, error) {
	// 获取用户订阅的所有conv
	getAllConv, err := NewAfterConnectLogic(l.ctx, l.svcCtx).GetAllConv(&pb.AfterConnectReq{
		ConnParam: &pb.ConnParam{
			UserId: in.CommonReq.UserId,
		},
		ConnectedAt: "",
	})
	if err != nil {
		l.Errorf("get all conv error: %v", err)
		return &pb.KeepAliveResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, convId := range getAllConv {
		xtrace.StartFuncSpan(l.ctx, "GetUserNoticeData", func(ctx context.Context) {
			_, err := NewGetUserNoticeDataLogic(l.ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				UserId:    in.CommonReq.UserId,
				ConvId:    convId,
				DeviceId:  utils.AnyPtr(in.CommonReq.DeviceId),
			})
			if err != nil {
				l.Errorf("KeepAlive failed, err: %v", err)
			}
		})
	}
	return &pb.KeepAliveResp{}, nil
}
