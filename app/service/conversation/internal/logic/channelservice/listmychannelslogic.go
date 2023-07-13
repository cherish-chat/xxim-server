package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyChannelsLogic {
	return &ListMyChannelsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyChannelsLogic) ListMyChannels(in *peerpb.ListMyChannelsReq) (*peerpb.ListMyChannelsResp, error) {
	{
		if in.Limit == 0 {
			in.Limit = 2000
		}
	}
	filter := bson.M{
		"memberUserId": in.Header.UserId,
	}
	conversationMembers := make([]*channelmodel.ChannelMember, 0)
	queryI := l.svcCtx.ChannelMemberCollection.Find(context.Background(), filter)
	queryI = queryI.Limit(int64(in.Limit))
	queryI = queryI.Skip(int64(in.Cursor))
	queryI = queryI.Sort("joinTime")
	err := queryI.All(&conversationMembers)
	if err != nil {
		l.Errorf("find conversation member error: %v", err)
		return &peerpb.ListMyChannelsResp{}, err
	}
	var resp = &peerpb.ListMyChannelsResp{}
	for _, conversationMember := range conversationMembers {
		resp.MyChannelList = append(resp.MyChannelList, &peerpb.ListMyChannelsResp_MyChannel{
			ChannelId: conversationMember.ChannelId,
			JoinTime:  uint32(conversationMember.JoinTime),
		})
	}
	return resp, nil
}
