package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessagePushLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessagePushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagePushLogic {
	return &MessagePushLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MessagePush 推送消息
func (l *MessagePushLogic) MessagePush(in *pb.MessagePushReq) (*pb.MessagePushResp, error) {
	// todo: add your logic here and delete this line
	l.Debugf("MessagePush 推送消息: %+v", in)
	for _, message := range in.Message {
		switch message.ConversationType {
		case pb.ConversationType_Single:
			l.pushSingleMessage(in, message)
		case pb.ConversationType_Group:
			l.pushGroupMessage(in, message)
		case pb.ConversationType_Subscription:
			l.pushSubscriptionMessage(in, message)
		}
	}
	return &pb.MessagePushResp{}, nil
}

func (l *MessagePushLogic) pushSingleMessage(in *pb.MessagePushReq, message *pb.Message) {
	fromId := message.GetSender().GetId()
	toId := pb.GetSingleChatOtherId(message.ConversationId, message.GetSender().GetId())
	var userIds []string
	userIds = append(userIds, fromId, toId)
	l.pushMessageToUserIds(in, message, userIds)
}

func (l *MessagePushLogic) pushGroupMessage(in *pb.MessagePushReq, message *pb.Message) {
	groupId := message.ConversationId
	listGroupSubscribersResp, err := l.svcCtx.GroupService.ListGroupSubscribers(l.ctx, &pb.ListGroupSubscribersReq{
		Header:  in.Header,
		GroupId: groupId,
		Cursor:  0,
		Limit:   0,
		Filter: &pb.ListGroupSubscribersReq_Filter{
			SubscribeTimeGte: time.Now().UnixMilli() - 1000*60*5, // 5分钟内在线的用户
		},
		Option: &pb.ListGroupSubscribersReq_Option{},
	})
	if err != nil {
		l.Errorf("get group subscribers error: %v", err)
		return
	}
	var userIds []string
	for _, subscriber := range listGroupSubscribersResp.SubscriberList {
		userIds = append(userIds, subscriber.UserId)
	}
	l.pushMessageToUserIds(in, message, userIds)

}

func (l *MessagePushLogic) pushSubscriptionMessage(in *pb.MessagePushReq, message *pb.Message) {
	subscriptionId := message.ConversationId
	listSubscriptionSubscribersResp, err := l.svcCtx.SubscriptionService.ListSubscriptionSubscribers(l.ctx, &pb.ListSubscriptionSubscribersReq{
		SubscriptionId: subscriptionId,
		Cursor:         0,
		Limit:          0,
		Filter: &pb.ListSubscriptionSubscribersReq_Filter{
			SubscribeTimeGte: time.Now().UnixMilli() - 1000*60*5, // 5分钟内在线的用户
		},
		Option: &pb.ListSubscriptionSubscribersReq_Option{},
	})
	if err != nil {
		l.Errorf("get subscription subscribers error: %v", err)
		return
	}
	var userIds []string
	for _, subscriber := range listSubscriptionSubscribersResp.SubscriberList {
		userIds = append(userIds, subscriber.UserId)
	}
	l.pushMessageToUserIds(in, message, userIds)
}

func (l *MessagePushLogic) pushMessageToUserIds(in *pb.MessagePushReq, message *pb.Message, userIds []string) {
	if len(userIds) > 0 {
		//推送消息
		gatewayWriteDataToWsResp, err := l.svcCtx.GatewayService.GatewayWriteDataToWsWrapper(context.Background(), &pb.GatewayWriteDataToWsWrapperReq{
			Header: in.Header,
			Filter: &pb.GatewayGetConnectionFilter{
				UserIds: userIds,
			},
			Data: &pb.GatewayWriteDataContent{
				DataType: pb.GatewayWriteDataType_PushMessage,
				Response: nil,
				Message:  message,
				Notice:   nil,
			},
		})
		if err != nil {
			l.Errorf("gateway write data to ws error: %v", err)
			return
		}
		_ = gatewayWriteDataToWsResp
	}
}
