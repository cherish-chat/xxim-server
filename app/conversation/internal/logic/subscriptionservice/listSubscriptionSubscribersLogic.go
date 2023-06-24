package subscriptionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSubscriptionSubscribersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSubscriptionSubscribersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSubscriptionSubscribersLogic {
	return &ListSubscriptionSubscribersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListSubscriptionSubscribers 列出订阅号订阅者
func (l *ListSubscriptionSubscribersLogic) ListSubscriptionSubscribers(in *pb.ListSubscriptionSubscribersReq) (*pb.ListSubscriptionSubscribersResp, error) {
	filter := bson.M{
		"subscriptionId": in.SubscriptionId,
	}
	// filter
	{
		if in.GetFilter().GetSubscribeTimeGte() >= 0 {
			filter["subscribeTime"] = bson.M{
				"$gte": primitive.DateTime(in.GetFilter().GetSubscribeTimeGte()),
			}
		}
	}
	queryI := l.svcCtx.SubscriptionSubscribeCollection.Find(l.ctx, filter)
	if in.Limit > 0 {
		queryI = queryI.Limit(in.Limit)
	}
	if in.Cursor > 0 {
		queryI = queryI.Skip(in.Cursor)
	}
	var result []*subscriptionmodel.SubscriptionSubscribe
	err := queryI.All(&result)
	if err != nil {
		l.Errorf("find subscription subscribe error: %v", err)
		return &pb.ListSubscriptionSubscribersResp{}, err
	}
	var resp = &pb.ListSubscriptionSubscribersResp{
		SubscriberList: make([]*pb.ListSubscriptionSubscribersResp_Subscriber, 0),
	}
	for _, item := range result {
		resp.SubscriberList = append(resp.SubscriberList, &pb.ListSubscriptionSubscribersResp_Subscriber{
			UserId:        item.MemberUserId,
			SubscribeTime: int64(item.SubscribeTime),
		})
	}
	return resp, nil
}
