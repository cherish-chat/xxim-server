package subscriptionservicelogic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubscriptionAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionAfterOnlineLogic {
	return &SubscriptionAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubscriptionAfterOnline 订阅号在做用户上线后的操作
func (l *SubscriptionAfterOnlineLogic) SubscriptionAfterOnline(in *pb.SubscriptionAfterOnlineReq) (*pb.SubscriptionAfterOnlineResp, error) {
	//1. 检查用户有没有创建他的默认订阅号
	count, err := l.svcCtx.SubscriptionCollection.Find(l.ctx, bson.M{
		"subscriptionId":   subscriptionmodel.UserDefaultSubscriptionId(in.Header.UserId),
		"subscriptionType": subscriptionmodel.SubscriptionTypeHidden,
	}).Count()
	if err != nil {
		l.Errorf("find subscription error: %v", err)
		return &pb.SubscriptionAfterOnlineResp{}, err
	}
	if count == 0 {
		//1.1. 如果没有创建，创建一个默认的订阅号
		_, err := l.svcCtx.SubscriptionCollection.InsertOne(l.ctx, &subscriptionmodel.Subscription{
			SubscriptionType: subscriptionmodel.SubscriptionTypeHidden,
			SubscriptionId:   subscriptionmodel.UserDefaultSubscriptionId(in.Header.UserId),
			Avatar:           "",
			Nickname:         "",
			Bio:              fmt.Sprintf("系统帮用户[%s]创建的默认订阅号，用于向订阅者推送通知。如：上线下线、发布世界圈、更新资料等", in.Header.UserId),
			Cover:            "",
			ExtraMap: bson.M{
				"defaultForUser": "true",
			},
		})
		if err != nil {
			l.Errorf("insert subscription error: %v", err)
			return &pb.SubscriptionAfterOnlineResp{}, err
		}
	}
	//2. 使用订阅号发一条通知，告诉他的订阅者，他上线了
	{
		_, err := l.svcCtx.NoticeService.NoticeSend(l.ctx, &pb.NoticeSendReq{
			Header: in.Header,
			Notice: &pb.Notice{
				NoticeId:         utils.Snowflake.String(),
				ConversationId:   subscriptionmodel.UserDefaultSubscriptionId(in.Header.UserId),
				ConversationType: pb.ConversationType_Subscription,
				Content:          utils.Json.MarshalToString(&pb.NoticeContentOnlineStatus{UserId: in.Header.UserId, Online: true}),
				ContentType:      pb.NoticeContentType_OnlineStatus,
				UpdateTime:       utils.Time.Now13(),
				Sort:             0,
			},
			UserIds:   nil,
			Broadcast: true,
		})
		if err != nil {
			l.Errorf("notice send error: %v", err)
		}
	}
	return &pb.SubscriptionAfterOnlineResp{}, nil
}
