package subscriptionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubscriptionAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionAfterOfflineLogic {
	return &SubscriptionAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubscriptionAfterOffline 订阅号在做用户下线后的操作
func (l *SubscriptionAfterOfflineLogic) SubscriptionAfterOffline(in *pb.SubscriptionAfterOfflineReq) (*pb.SubscriptionAfterOfflineResp, error) {
	//1. 使用订阅号发一条通知，告诉他的订阅者，他上线了
	{
		_, err := l.svcCtx.NoticeService.NoticeSend(l.ctx, &pb.NoticeSendReq{
			Header: in.Header,
			Notice: &pb.Notice{
				NoticeId:         utils.Snowflake.String(),
				ConversationId:   subscriptionmodel.UserDefaultSubscriptionId(in.Header.UserId),
				ConversationType: pb.ConversationType_Subscription,
				Content:          utils.Json.MarshalToString(&pb.NoticeContentOnlineStatus{UserId: in.Header.UserId, Online: false}),
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
	return &pb.SubscriptionAfterOfflineResp{}, nil
}
