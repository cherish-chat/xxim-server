package subscriptionservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionSubscribeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubscriptionSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionSubscribeLogic {
	return &SubscriptionSubscribeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubscriptionSubscribe 订阅号订阅
func (l *SubscriptionSubscribeLogic) SubscriptionSubscribe(in *pb.SubscriptionSubscribeReq) (*pb.SubscriptionSubscribeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SubscriptionSubscribeResp{}, nil
}
