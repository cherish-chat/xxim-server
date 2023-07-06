package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOfflineLogic {
	return &UserAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAfterOffline 用户下线回调
func (l *UserAfterOfflineLogic) UserAfterOffline(in *pb.UserAfterOfflineReq) (*pb.UserAfterOfflineResp, error) {
	//1. 订阅号的逻辑
	{
		_, err := l.svcCtx.SubscriptionService.SubscriptionAfterOffline(l.ctx, &pb.SubscriptionAfterOfflineReq{
			UserId: in.UserId,
		})
		if err != nil {
			l.Errorf("subscription after online error: %v", err)
		}
	}
	return &pb.UserAfterOfflineResp{}, nil
}
