package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChannelAfterKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChannelAfterKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChannelAfterKeepAliveLogic {
	return &ChannelAfterKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChannelAfterKeepAliveLogic) ChannelAfterKeepAlive(in *peerpb.ChannelAfterKeepAliveReq) (*peerpb.ChannelAfterKeepAliveResp, error) {
	listMyChannelsResp, err := NewListMyChannelsLogic(context.Background(), l.svcCtx).ListMyChannels(&peerpb.ListMyChannelsReq{
		Header: in.Header,
		Cursor: 0,
		Limit:  2000,
		Filter: &peerpb.ListMyChannelsReq_Filter{},
		Option: &peerpb.ListMyChannelsReq_Option{},
	})
	if err != nil {
		l.Errorf("list my channels error: %v", err)
		return nil, err
	}
	l.Debugf("list my channels: %v", listMyChannelsResp)
	var subscriptionIds = channelmodel.UserSubscribedSystemChannelIds()
	for _, channel := range listMyChannelsResp.MyChannelList {
		subscriptionIds = append(subscriptionIds, channel.ChannelId)
	}
	// 2. 批量更新用户的频道订阅时间
	// filter: {memberUserId: "xxx", channelId: {$in: ["xxx", "xxx"]}}
	// set: {subscribeTime: "xxx"}
	bulk := l.svcCtx.ChannelSubscribeCacheCollection.Bulk()
	for _, channelId := range subscriptionIds {
		bulk.Upsert(bson.M{
			"memberUserId": in.Header.UserId,
			"channelId":    channelId,
		}, &channelmodel.ChannelSubscribeCache{
			ChannelId:     channelId,
			MemberUserId:  in.Header.UserId,
			SubscribeTime: primitive.NewDateTimeFromTime(time.Now()),
		})
	}
	_, err = bulk.Run(context.Background())
	if err != nil {
		l.Errorf("bulk.Run err: %v", err)
		return &peerpb.ChannelAfterKeepAliveResp{}, err
	}
	return &peerpb.ChannelAfterKeepAliveResp{}, nil
}
