package groupservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGroupSubscribersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGroupSubscribersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupSubscribersLogic {
	return &ListGroupSubscribersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListGroupSubscribers 列出群组订阅者
func (l *ListGroupSubscribersLogic) ListGroupSubscribers(in *pb.ListGroupSubscribersReq) (*pb.ListGroupSubscribersResp, error) {
	filter := bson.M{
		"groupId": in.GroupId,
	}
	// filter
	{
		if in.GetFilter().GetSubscribeTimeGte() >= 0 {
			filter["subscribeTime"] = bson.M{
				"$gte": primitive.DateTime(in.GetFilter().GetSubscribeTimeGte()),
			}
		}
	}
	queryI := l.svcCtx.GroupSubscribeCollection.Find(l.ctx, filter)
	if in.Limit > 0 {
		queryI = queryI.Limit(in.Limit)
	}
	if in.Cursor > 0 {
		queryI = queryI.Skip(in.Cursor)
	}
	var result []*groupmodel.GroupSubscribe
	err := queryI.All(&result)
	if err != nil {
		l.Errorf("find group subscribe error: %v", err)
		return &pb.ListGroupSubscribersResp{}, err
	}
	var resp = &pb.ListGroupSubscribersResp{
		SubscriberList: make([]*pb.ListGroupSubscribersResp_Subscriber, 0),
	}
	for _, item := range result {
		resp.SubscriberList = append(resp.SubscriberList, &pb.ListGroupSubscribersResp_Subscriber{
			UserId:        item.MemberUserId,
			SubscribeTime: int64(item.SubscribeTime),
		})
	}
	return resp, nil
}
