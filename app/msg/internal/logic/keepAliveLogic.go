package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
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

func (l *KeepAliveLogic) KeepAlive(in *pb.KeepAliveReq) (*pb.KeepAliveResp, error) {
	// 是否进入冷却，如果进入冷却，不做任何处理，否则设置30s的冷却时间
	// 1. 判断是否进入冷却
	lockKey := "lock:keepAlive:" + in.GetCommonReq().GetUserId()
	exist, err := l.svcCtx.Redis().ExistsCtx(context.Background(), lockKey)
	if err != nil {
		l.Errorf("exists lock error: %v", err)
		return &pb.KeepAliveResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if exist {
		// 进入冷却
		return &pb.KeepAliveResp{}, nil
	}
	// 2. 设置冷却时间
	err = l.svcCtx.Redis().SetexCtx(context.Background(), lockKey, "1", 30)
	if err != nil {
		l.Errorf("set lock error: %v", err)
		return &pb.KeepAliveResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 3. 更新订阅关系
	NewFlushUsersSubConvLogic(l.ctx, l.svcCtx).SetUserSubscriptions(in.GetCommonReq().GetUserId(), in.GetConvIdList())
	return &pb.KeepAliveResp{}, nil
}
