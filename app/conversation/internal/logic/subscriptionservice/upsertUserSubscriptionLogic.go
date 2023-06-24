package subscriptionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertUserSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertUserSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertUserSubscriptionLogic {
	return &UpsertUserSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpsertUserSubscription 更新用户订阅的订阅号
func (l *UpsertUserSubscriptionLogic) UpsertUserSubscription(in *pb.UpsertUserSubscriptionReq) (*pb.UpsertUserSubscriptionResp, error) {
	userSubscription := &subscriptionmodel.UserSubscription{
		SubscriptionId: in.UserSubscription.SubscriptionId,
		Subscriber:     in.UserSubscription.Subscriber,
		SubscribeTime:  primitive.DateTime(in.UserSubscription.SubscribeTime),
		ExtraMap:       utils.Map.SS2SA(in.UserSubscription.ExtraMap),
	}
	_, err := l.svcCtx.UserSubscriptionCollection.Upsert(l.ctx, bson.M{
		"subscriptionId": userSubscription.SubscriptionId,
		"subscriber":     userSubscription.Subscriber,
	}, userSubscription, opts.UpsertOptions{
		ReplaceOptions: options.Replace().SetUpsert(true),
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return &pb.UpsertUserSubscriptionResp{}, err
	}
	return &pb.UpsertUserSubscriptionResp{}, nil
}
