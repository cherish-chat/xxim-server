package noticeservicelogic

import (
	"context"
	"time"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticePushLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticePushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticePushLogic {
	return &NoticePushLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// NoticePush 推送消息
func (l *NoticePushLogic) NoticePush(in *peerpb.NoticeSendReq) (*peerpb.NoticeSendResp, error) {
	for _, message := range in.Notices {
		switch message.ConversationType {
		case peerpb.ConversationType_Single:
			l.pushSingleNotice(in, message)
		case peerpb.ConversationType_Group:
			l.pushGroupNotice(in, message)
		case peerpb.ConversationType_Channel:
			l.pushChannelNotice(in, message)
		}
	}
	return &peerpb.NoticeSendResp{}, nil
}
func (l *NoticePushLogic) pushSingleNotice(in *peerpb.NoticeSendReq, message *peerpb.Message) {
	fromId := message.GetSender().GetId()
	toId := peerpb.GetSingleChatOtherId(message.ConversationId, message.GetSender().GetId())
	var userIds []string
	userIds = append(userIds, fromId, toId)
	l.pushNoticeToUserIds(in, message, userIds)
}

func (l *NoticePushLogic) pushGroupNotice(in *peerpb.NoticeSendReq, message *peerpb.Message) {
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
	l.pushNoticeToUserIds(in, message, userIds)

}

func (l *NoticePushLogic) pushChannelNotice(in *peerpb.NoticeSendReq, message *peerpb.Message) {
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
	l.pushNoticeToUserIds(in, message, userIds)
}

func (l *NoticePushLogic) pushNoticeToUserIds(in *peerpb.NoticeSendReq, message *peerpb.Message, userIds []string) {
	if len(userIds) > 0 {
		//推送消息
		gatewayWriteDataToWsResp, err := l.svcCtx.InternalService.GatewayWriteData(context.Background(), &peerpb.GatewayWriteDataReq{
			Header: in.Header,
			Filter: &peerpb.GatewayWriteDataReq_Filter{
				UserIds: userIds,
			},
			Content: &peerpb.GatewayWriteDataContent{
				DataType: peerpb.GatewayWriteDataType_PushNotice,
				Response: nil,
				Notice:   message,
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
