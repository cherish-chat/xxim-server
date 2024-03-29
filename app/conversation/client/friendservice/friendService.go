// Code generated by goctl. DO NOT EDIT.
// Source: conversation.proto

package friendservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ConversationSettingUpdateReq                            = pb.ConversationSettingUpdateReq
	ConversationSettingUpdateResp                           = pb.ConversationSettingUpdateResp
	CountCreateGroupReq                                     = pb.CountCreateGroupReq
	CountCreateGroupResp                                    = pb.CountCreateGroupResp
	CountFriendReq                                          = pb.CountFriendReq
	CountFriendResp                                         = pb.CountFriendResp
	CountJoinGroupReq                                       = pb.CountJoinGroupReq
	CountJoinGroupResp                                      = pb.CountJoinGroupResp
	DeleteUserSubscriptionReq                               = pb.DeleteUserSubscriptionReq
	DeleteUserSubscriptionResp                              = pb.DeleteUserSubscriptionResp
	FriendApplyHandleReq                                    = pb.FriendApplyHandleReq
	FriendApplyHandleResp                                   = pb.FriendApplyHandleResp
	FriendApplyReq                                          = pb.FriendApplyReq
	FriendApplyResp                                         = pb.FriendApplyResp
	GroupCreateReq                                          = pb.GroupCreateReq
	GroupCreateResp                                         = pb.GroupCreateResp
	GroupSubscribeReq                                       = pb.GroupSubscribeReq
	GroupSubscribeResp                                      = pb.GroupSubscribeResp
	ListFriendApplyReq                                      = pb.ListFriendApplyReq
	ListFriendApplyReq_Filter                               = pb.ListFriendApplyReq_Filter
	ListFriendApplyReq_Option                               = pb.ListFriendApplyReq_Option
	ListFriendApplyResp                                     = pb.ListFriendApplyResp
	ListFriendApplyResp_FriendApply                         = pb.ListFriendApplyResp_FriendApply
	ListGroupSubscribersReq                                 = pb.ListGroupSubscribersReq
	ListGroupSubscribersReq_Filter                          = pb.ListGroupSubscribersReq_Filter
	ListGroupSubscribersReq_Option                          = pb.ListGroupSubscribersReq_Option
	ListGroupSubscribersResp                                = pb.ListGroupSubscribersResp
	ListGroupSubscribersResp_Subscriber                     = pb.ListGroupSubscribersResp_Subscriber
	ListJoinedConversationsReq                              = pb.ListJoinedConversationsReq
	ListJoinedConversationsReq_Filter                       = pb.ListJoinedConversationsReq_Filter
	ListJoinedConversationsReq_Filter_SettingKV             = pb.ListJoinedConversationsReq_Filter_SettingKV
	ListJoinedConversationsReq_Option                       = pb.ListJoinedConversationsReq_Option
	ListJoinedConversationsResp                             = pb.ListJoinedConversationsResp
	ListJoinedConversationsResp_Conversation                = pb.ListJoinedConversationsResp_Conversation
	ListJoinedConversationsResp_Conversation_SelfMemberInfo = pb.ListJoinedConversationsResp_Conversation_SelfMemberInfo
	ListSubscriptionSubscribersReq                          = pb.ListSubscriptionSubscribersReq
	ListSubscriptionSubscribersReq_Filter                   = pb.ListSubscriptionSubscribersReq_Filter
	ListSubscriptionSubscribersReq_Option                   = pb.ListSubscriptionSubscribersReq_Option
	ListSubscriptionSubscribersResp                         = pb.ListSubscriptionSubscribersResp
	ListSubscriptionSubscribersResp_Subscriber              = pb.ListSubscriptionSubscribersResp_Subscriber
	SubscriptionAfterOfflineReq                             = pb.SubscriptionAfterOfflineReq
	SubscriptionAfterOfflineResp                            = pb.SubscriptionAfterOfflineResp
	SubscriptionAfterOnlineReq                              = pb.SubscriptionAfterOnlineReq
	SubscriptionAfterOnlineResp                             = pb.SubscriptionAfterOnlineResp
	SubscriptionSubscribeReq                                = pb.SubscriptionSubscribeReq
	SubscriptionSubscribeResp                               = pb.SubscriptionSubscribeResp
	UpsertUserSubscriptionReq                               = pb.UpsertUserSubscriptionReq
	UpsertUserSubscriptionResp                              = pb.UpsertUserSubscriptionResp
	UserSubscription                                        = pb.UserSubscription

	FriendService interface {
		// FriendApply 添加好友
		FriendApply(ctx context.Context, in *FriendApplyReq, opts ...grpc.CallOption) (*FriendApplyResp, error)
		// FriendApplyHandle 处理好友申请
		FriendApplyHandle(ctx context.Context, in *FriendApplyHandleReq, opts ...grpc.CallOption) (*FriendApplyHandleResp, error)
		// ListFriendApply 列出好友申请
		ListFriendApply(ctx context.Context, in *ListFriendApplyReq, opts ...grpc.CallOption) (*ListFriendApplyResp, error)
		// CountFriend 统计好友数量
		CountFriend(ctx context.Context, in *CountFriendReq, opts ...grpc.CallOption) (*CountFriendResp, error)
	}

	defaultFriendService struct {
		cli zrpc.Client
	}
)

func NewFriendService(cli zrpc.Client) FriendService {
	return &defaultFriendService{
		cli: cli,
	}
}

// FriendApply 添加好友
func (m *defaultFriendService) FriendApply(ctx context.Context, in *FriendApplyReq, opts ...grpc.CallOption) (*FriendApplyResp, error) {
	client := pb.NewFriendServiceClient(m.cli.Conn())
	return client.FriendApply(ctx, in, opts...)
}

// FriendApplyHandle 处理好友申请
func (m *defaultFriendService) FriendApplyHandle(ctx context.Context, in *FriendApplyHandleReq, opts ...grpc.CallOption) (*FriendApplyHandleResp, error) {
	client := pb.NewFriendServiceClient(m.cli.Conn())
	return client.FriendApplyHandle(ctx, in, opts...)
}

// ListFriendApply 列出好友申请
func (m *defaultFriendService) ListFriendApply(ctx context.Context, in *ListFriendApplyReq, opts ...grpc.CallOption) (*ListFriendApplyResp, error) {
	client := pb.NewFriendServiceClient(m.cli.Conn())
	return client.ListFriendApply(ctx, in, opts...)
}

// CountFriend 统计好友数量
func (m *defaultFriendService) CountFriend(ctx context.Context, in *CountFriendReq, opts ...grpc.CallOption) (*CountFriendResp, error) {
	client := pb.NewFriendServiceClient(m.cli.Conn())
	return client.CountFriend(ctx, in, opts...)
}
