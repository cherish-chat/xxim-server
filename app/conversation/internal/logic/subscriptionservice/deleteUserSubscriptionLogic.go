package subscriptionservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserSubscriptionLogic {
	return &DeleteUserSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteUserSubscription 删除用户订阅的订阅号
func (l *DeleteUserSubscriptionLogic) DeleteUserSubscription(in *pb.DeleteUserSubscriptionReq) (*pb.DeleteUserSubscriptionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteUserSubscriptionResp{}, nil
}
