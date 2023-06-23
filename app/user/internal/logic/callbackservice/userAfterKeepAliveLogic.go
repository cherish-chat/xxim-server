package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterKeepAliveLogic {
	return &UserAfterKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAfterKeepAlive 用户保活回调
func (l *UserAfterKeepAliveLogic) UserAfterKeepAlive(in *pb.UserAfterKeepAliveReq) (*pb.UserAfterKeepAliveResp, error) {
	// todo: add your logic here and delete this line
	// 1. 订阅群聊消息
	{
		_, err := l.svcCtx.GroupService.GroupSubscribe(l.ctx, &pb.GroupSubscribeReq{
			Header: in.Header,
		})
		if err != nil {
			l.Errorf("subscribe group message error: %v", err)
		}
	}
	// 2. 订阅 订阅号 消息
	{
		_, err := l.svcCtx.SubscriptionService.SubscriptionSubscribe(l.ctx, &pb.SubscriptionSubscribeReq{
			Header: in.Header,
		})
		if err != nil {
			l.Errorf("subscribe subscription message error: %v", err)
		}
	}
	return &pb.UserAfterKeepAliveResp{}, nil
}
