package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertChannelMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertChannelMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertChannelMemberLogic {
	return &UpsertChannelMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpsertChannelMemberLogic) UpsertChannelMember(in *peerpb.UpsertChannelMemberReq) (*peerpb.UpsertChannelMemberResp, error) {
	userSubscription := &channelmodel.ChannelMember{
		ChannelId:    in.UserChannel.ChannelId,
		MemberUserId: in.UserChannel.UserId,
		JoinTime:     primitive.DateTime(in.UserChannel.SubscribeTime),
	}
	_, err := l.svcCtx.ChannelMemberCollection.Upsert(context.Background(), bson.M{
		"channelId":    userSubscription.ChannelId,
		"memberUserId": userSubscription.MemberUserId,
	}, userSubscription, opts.UpsertOptions{
		ReplaceOptions: options.Replace().SetUpsert(true),
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return &peerpb.UpsertChannelMemberResp{}, err
	}
	return &peerpb.UpsertChannelMemberResp{}, nil
}
