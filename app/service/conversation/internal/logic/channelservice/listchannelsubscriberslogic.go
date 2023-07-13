package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListChannelSubscribersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListChannelSubscribersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListChannelSubscribersLogic {
	return &ListChannelSubscribersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListChannelSubscribersLogic) ListChannelSubscribers(in *peerpb.ListChannelSubscribersReq) (*peerpb.ListChannelSubscribersResp, error) {
	filter := bson.M{
		"channelId": in.ChannelId,
	}
	// filter
	{
		if in.GetFilter().GetSubscribeTimeGte() > 0 {
			filter["subscribeTime"] = bson.M{
				"$gte": primitive.DateTime(in.GetFilter().GetSubscribeTimeGte()),
			}
		}
	}
	queryI := l.svcCtx.ChannelSubscribeCacheCollection.Find(context.Background(), filter)
	if in.Limit > 0 {
		queryI = queryI.Limit(int64(in.Limit))
	}
	if in.Cursor > 0 {
		queryI = queryI.Skip(int64(in.Cursor))
	}
	var result []*channelmodel.ChannelSubscribeCache
	err := queryI.All(&result)
	if err != nil {
		l.Errorf("find subscription subscribe error: %v", err)
		return &peerpb.ListChannelSubscribersResp{}, err
	}
	var resp = &peerpb.ListChannelSubscribersResp{
		SubscriberList: make([]*peerpb.ListChannelSubscribersResp_Subscriber, 0),
	}
	for _, item := range result {
		resp.SubscriberList = append(resp.SubscriberList, &peerpb.ListChannelSubscribersResp_Subscriber{
			UserId:        item.MemberUserId,
			SubscribeTime: uint32(item.SubscribeTime),
		})
	}
	return resp, nil
}
