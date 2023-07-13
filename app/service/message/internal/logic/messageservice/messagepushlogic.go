package messageservicelogic

import (
	"context"
	"time"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

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
func (l *MessagePushLogic) MessagePush(in *peerpb.MessagePushReq) (*peerpb.MessagePushResp, error) {
	for _, message := range in.Message {
		switch message.ConversationType {
		case peerpb.ConversationType_Single:
			l.pushSingleMessage(in, message)
		case peerpb.ConversationType_Group:
			l.pushGroupMessage(in, message)
		case peerpb.ConversationType_Channel:
			l.pushChannelMessage(in, message)
		}
	}
	return &peerpb.MessagePushResp{}, nil
}
func (l *MessagePushLogic) pushSingleMessage(in *peerpb.MessagePushReq, message *peerpb.Message) {
	fromId := message.GetSender().GetId()
	toId := peerpb.GetSingleChatOtherId(message.ConversationId, message.GetSender().GetId())
	var userIds []string
	userIds = append(userIds, fromId, toId)
	l.pushMessageToUserIds(in, message, userIds)
}

func (l *MessagePushLogic) pushGroupMessage(in *peerpb.MessagePushReq, message *peerpb.Message) {
	groupId := message.ConversationId
	listGroupSubscribersResp, err := l.svcCtx.GroupService.ListGroupSubscribers(context.Background(), &peerpb.ListGroupSubscribersReq{
		Header:  in.Header,
		GroupId: groupId,
		Cursor:  0,
		Limit:   0,
		Filter: &peerpb.ListGroupSubscribersReq_Filter{
			SubscribeTimeGte: uint32(time.Now().UnixMilli() - 1000*60*5), // 5分钟内在线的用户
		},
		Option: &peerpb.ListGroupSubscribersReq_Option{},
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

func (l *MessagePushLogic) pushChannelMessage(in *peerpb.MessagePushReq, message *peerpb.Message) {
	subscriptionId := message.ConversationId
	listChannelSubscribersResp, err := l.svcCtx.ChannelService.ListChannelSubscribers(context.Background(), &peerpb.ListChannelSubscribersReq{
		ChannelId: subscriptionId,
		Cursor:    0,
		Limit:     0,
		Filter: &peerpb.ListChannelSubscribersReq_Filter{
			SubscribeTimeGte: uint32(time.Now().UnixMilli() - 1000*60*5), // 5分钟内在线的用户
		},
		Option: &peerpb.ListChannelSubscribersReq_Option{},
	})
	if err != nil {
		l.Errorf("get subscription subscribers error: %v", err)
		return
	}
	var userIds []string
	for _, subscriber := range listChannelSubscribersResp.SubscriberList {
		userIds = append(userIds, subscriber.UserId)
	}
	l.pushMessageToUserIds(in, message, userIds)
}

func (l *MessagePushLogic) pushMessageToUserIds(in *peerpb.MessagePushReq, message *peerpb.Message, userIds []string) {
	if len(userIds) > 0 {
		//推送消息
		gatewayWriteDataToWsResp, err := l.svcCtx.InternalService.GatewayWriteData(context.Background(), &peerpb.GatewayWriteDataReq{
			Header: in.Header,
			Filter: &peerpb.GatewayWriteDataReq_Filter{
				UserIds: userIds,
			},
			Content: &peerpb.GatewayWriteDataContent{
				DataType: peerpb.GatewayWriteDataType_PushMessage,
				Response: nil,
				Message:  message,
				Notice:   nil,
			},
			AllPods: true,
		})
		if err != nil {
			l.Errorf("gateway write data to ws error: %v", err)
			return
		}
		_ = gatewayWriteDataToWsResp
	}
}
