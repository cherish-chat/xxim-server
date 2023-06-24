package subscriptionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

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
	// 1. 获取用户订阅的订阅号列表
	listJoinedSubscriptionsResp, err := l.svcCtx.ConversationService.ListJoinedConversations(l.ctx, &pb.ListJoinedConversationsReq{
		Header:           in.Header,
		ConversationType: pb.ConversationType_Subscription,
		Cursor:           0,
		Limit:            1000000,
		Filter: &pb.ListJoinedConversationsReq_Filter{
			SettingList: []*pb.ListJoinedConversationsReq_Filter_SettingKV{{
				//是否被屏蔽这个key != "true" or 不存在
				Key:         pb.ConversationSettingKey_IsBlocked, // 是否被屏蔽
				Value:       "true",
				Operator:    pb.ListJoinedConversationsReq_Filter_SettingKV_NotEqual,
				OrNotExists: true,
				OrExists:    false,
			}},
		},
		Option: &pb.ListJoinedConversationsReq_Option{
			IncludeSelfMemberInfo: true,
		},
	})
	if err != nil {
		l.Errorf("listJoinedSubscriptionsResp err: %v", err)
		return &pb.SubscriptionSubscribeResp{}, err
	}
	if len(listJoinedSubscriptionsResp.GetConversationList()) == 0 {
		return &pb.SubscriptionSubscribeResp{}, nil
	}

	var subscriptionIds = subscriptionmodel.UserSubscribedSystemConversationIds()
	for _, subscription := range listJoinedSubscriptionsResp.GetConversationList() {
		subscriptionIds = append(subscriptionIds, subscription.GetConversationId())
	}

	// 2. 批量更新用户的群组订阅时间
	// filter: {memberUserId: "xxx", subscriptionId: {$in: ["xxx", "xxx"]}}
	// set: {subscribeTime: "xxx"}
	bulk := l.svcCtx.SubscriptionSubscribeCollection.Bulk()
	for _, subscriptionId := range subscriptionIds {
		bulk.Upsert(bson.M{
			"memberUserId":   in.Header.UserId,
			"subscriptionId": subscriptionId,
		}, &subscriptionmodel.SubscriptionSubscribe{
			SubscriptionId: subscriptionId,
			MemberUserId:   in.Header.UserId,
			SubscribeTime:  primitive.NewDateTimeFromTime(time.Now()),
		})
	}
	_, err = bulk.Run(l.ctx)
	if err != nil {
		l.Errorf("bulk.Run err: %v", err)
		return &pb.SubscriptionSubscribeResp{}, err
	}
	return &pb.SubscriptionSubscribeResp{}, nil
}
