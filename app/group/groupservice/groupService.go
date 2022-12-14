// Code generated by goctl. DO NOT EDIT!
// Source: group.proto

package groupservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BanAllGroupMemberReq                           = pb.BanAllGroupMemberReq
	BanAllGroupMemberResp                          = pb.BanAllGroupMemberResp
	BanGroupMemberReq                              = pb.BanGroupMemberReq
	BanGroupMemberResp                             = pb.BanGroupMemberResp
	CreateGroupNoticeReq                           = pb.CreateGroupNoticeReq
	CreateGroupNoticeResp                          = pb.CreateGroupNoticeResp
	CreateGroupReq                                 = pb.CreateGroupReq
	CreateGroupResp                                = pb.CreateGroupResp
	DeleteGroupNoticeReq                           = pb.DeleteGroupNoticeReq
	DeleteGroupNoticeResp                          = pb.DeleteGroupNoticeResp
	DismissGroupReq                                = pb.DismissGroupReq
	DismissGroupResp                               = pb.DismissGroupResp
	EditGroupInfoReq                               = pb.EditGroupInfoReq
	EditGroupInfoResp                              = pb.EditGroupInfoResp
	EditGroupNoticeReq                             = pb.EditGroupNoticeReq
	EditGroupNoticeResp                            = pb.EditGroupNoticeResp
	GetGroupHomeReq                                = pb.GetGroupHomeReq
	GetGroupHomeResp                               = pb.GetGroupHomeResp
	GetGroupHomeResp_MemberStatistics              = pb.GetGroupHomeResp_MemberStatistics
	GetGroupMemberInfoReq                          = pb.GetGroupMemberInfoReq
	GetGroupMemberInfoResp                         = pb.GetGroupMemberInfoResp
	GetGroupMemberListReq                          = pb.GetGroupMemberListReq
	GetGroupMemberListReq_GetGroupMemberListFilter = pb.GetGroupMemberListReq_GetGroupMemberListFilter
	GetGroupMemberListReq_GetGroupMemberListOpt    = pb.GetGroupMemberListReq_GetGroupMemberListOpt
	GetGroupMemberListResp                         = pb.GetGroupMemberListResp
	GetGroupMemberListResp_GroupMember             = pb.GetGroupMemberListResp_GroupMember
	GetGroupNoticeListReq                          = pb.GetGroupNoticeListReq
	GetGroupNoticeListResp                         = pb.GetGroupNoticeListResp
	GetGroupSettingReq                             = pb.GetGroupSettingReq
	GetGroupSettingResp                            = pb.GetGroupSettingResp
	GetMyGroupListReq                              = pb.GetMyGroupListReq
	GetMyGroupListReq_Filter                       = pb.GetMyGroupListReq_Filter
	GetMyGroupListResp                             = pb.GetMyGroupListResp
	GetMyGroupListResp_Group                       = pb.GetMyGroupListResp_Group
	GroupMemberInfo                                = pb.GroupMemberInfo
	GroupNotice                                    = pb.GroupNotice
	GroupSetting                                   = pb.GroupSetting
	GroupSetting_JoinGroupOpt                      = pb.GroupSetting_JoinGroupOpt
	GroupSetting_MemberPermission                  = pb.GroupSetting_MemberPermission
	InviteFriendToGroupReq                         = pb.InviteFriendToGroupReq
	InviteFriendToGroupResp                        = pb.InviteFriendToGroupResp
	KickGroupMemberReq                             = pb.KickGroupMemberReq
	KickGroupMemberResp                            = pb.KickGroupMemberResp
	QuitGroupReq                                   = pb.QuitGroupReq
	QuitGroupResp                                  = pb.QuitGroupResp
	SetGroupMemberInfoReq                          = pb.SetGroupMemberInfoReq
	SetGroupMemberInfoResp                         = pb.SetGroupMemberInfoResp
	SetGroupMemberRoleReq                          = pb.SetGroupMemberRoleReq
	SetGroupMemberRoleResp                         = pb.SetGroupMemberRoleResp
	SetGroupMsgNotifyTypeReq                       = pb.SetGroupMsgNotifyTypeReq
	SetGroupMsgNotifyTypeResp                      = pb.SetGroupMsgNotifyTypeResp
	SetGroupSettingReq                             = pb.SetGroupSettingReq
	SetGroupSettingResp                            = pb.SetGroupSettingResp
	TransferGroupOwnerReq                          = pb.TransferGroupOwnerReq
	TransferGroupOwnerResp                         = pb.TransferGroupOwnerResp
	UnbanAllGroupMemberReq                         = pb.UnbanAllGroupMemberReq
	UnbanAllGroupMemberResp                        = pb.UnbanAllGroupMemberResp
	UnbanGroupMemberReq                            = pb.UnbanGroupMemberReq
	UnbanGroupMemberResp                           = pb.UnbanGroupMemberResp

	GroupService interface {
		// CreateGroup ????????????
		CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateGroupResp, error)
		// GetGroupHome ??????????????????
		GetGroupHome(ctx context.Context, in *GetGroupHomeReq, opts ...grpc.CallOption) (*GetGroupHomeResp, error)
		// InviteFriendToGroup ????????????????????????
		InviteFriendToGroup(ctx context.Context, in *InviteFriendToGroupReq, opts ...grpc.CallOption) (*InviteFriendToGroupResp, error)
		// CreateGroupNotice ???????????????
		CreateGroupNotice(ctx context.Context, in *CreateGroupNoticeReq, opts ...grpc.CallOption) (*CreateGroupNoticeResp, error)
		// DeleteGroupNotice ???????????????
		DeleteGroupNotice(ctx context.Context, in *DeleteGroupNoticeReq, opts ...grpc.CallOption) (*DeleteGroupNoticeResp, error)
		// GetGroupNoticeList ?????????????????????
		GetGroupNoticeList(ctx context.Context, in *GetGroupNoticeListReq, opts ...grpc.CallOption) (*GetGroupNoticeListResp, error)
		// SetGroupMemberInfo ?????????????????????
		SetGroupMemberInfo(ctx context.Context, in *SetGroupMemberInfoReq, opts ...grpc.CallOption) (*SetGroupMemberInfoResp, error)
		// GetGroupMemberInfo ?????????????????????
		GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoReq, opts ...grpc.CallOption) (*GetGroupMemberInfoResp, error)
		// EditGroupInfo ???????????????
		EditGroupInfo(ctx context.Context, in *EditGroupInfoReq, opts ...grpc.CallOption) (*EditGroupInfoResp, error)
		// SetGroupSetting ???????????????
		SetGroupSetting(ctx context.Context, in *SetGroupSettingReq, opts ...grpc.CallOption) (*SetGroupSettingResp, error)
		// GetGroupSetting ???????????????
		GetGroupSetting(ctx context.Context, in *GetGroupSettingReq, opts ...grpc.CallOption) (*GetGroupSettingResp, error)
		// TransferGroupOwner ????????????
		TransferGroupOwner(ctx context.Context, in *TransferGroupOwnerReq, opts ...grpc.CallOption) (*TransferGroupOwnerResp, error)
		// SetGroupMemberRole ?????????????????????
		SetGroupMemberRole(ctx context.Context, in *SetGroupMemberRoleReq, opts ...grpc.CallOption) (*SetGroupMemberRoleResp, error)
		// KickGroupMember ???????????????
		KickGroupMember(ctx context.Context, in *KickGroupMemberReq, opts ...grpc.CallOption) (*KickGroupMemberResp, error)
		// QuitGroup ????????????
		QuitGroup(ctx context.Context, in *QuitGroupReq, opts ...grpc.CallOption) (*QuitGroupResp, error)
		// BanGroupMember ???????????????
		BanGroupMember(ctx context.Context, in *BanGroupMemberReq, opts ...grpc.CallOption) (*BanGroupMemberResp, error)
		// BanAllGroupMember ?????????????????????
		BanAllGroupMember(ctx context.Context, in *BanAllGroupMemberReq, opts ...grpc.CallOption) (*BanAllGroupMemberResp, error)
		// UnbanGroupMember ?????????????????????
		UnbanGroupMember(ctx context.Context, in *UnbanGroupMemberReq, opts ...grpc.CallOption) (*UnbanGroupMemberResp, error)
		// UnbanAllGroupMember ???????????????????????????
		UnbanAllGroupMember(ctx context.Context, in *UnbanAllGroupMemberReq, opts ...grpc.CallOption) (*UnbanAllGroupMemberResp, error)
		// GetGroupMemberList ?????????????????????
		GetGroupMemberList(ctx context.Context, in *GetGroupMemberListReq, opts ...grpc.CallOption) (*GetGroupMemberListResp, error)
		// DismissGroup ????????????
		DismissGroup(ctx context.Context, in *DismissGroupReq, opts ...grpc.CallOption) (*DismissGroupResp, error)
		// SetGroupMsgNotifyType ???????????????????????????
		SetGroupMsgNotifyType(ctx context.Context, in *SetGroupMsgNotifyTypeReq, opts ...grpc.CallOption) (*SetGroupMsgNotifyTypeResp, error)
		// GetMyGroupList ????????????????????????
		GetMyGroupList(ctx context.Context, in *GetMyGroupListReq, opts ...grpc.CallOption) (*GetMyGroupListResp, error)
	}

	defaultGroupService struct {
		cli zrpc.Client
	}
)

func NewGroupService(cli zrpc.Client) GroupService {
	return &defaultGroupService{
		cli: cli,
	}
}

// CreateGroup ????????????
func (m *defaultGroupService) CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateGroupResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.CreateGroup(ctx, in, opts...)
}

// GetGroupHome ??????????????????
func (m *defaultGroupService) GetGroupHome(ctx context.Context, in *GetGroupHomeReq, opts ...grpc.CallOption) (*GetGroupHomeResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupHome(ctx, in, opts...)
}

// InviteFriendToGroup ????????????????????????
func (m *defaultGroupService) InviteFriendToGroup(ctx context.Context, in *InviteFriendToGroupReq, opts ...grpc.CallOption) (*InviteFriendToGroupResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.InviteFriendToGroup(ctx, in, opts...)
}

// CreateGroupNotice ???????????????
func (m *defaultGroupService) CreateGroupNotice(ctx context.Context, in *CreateGroupNoticeReq, opts ...grpc.CallOption) (*CreateGroupNoticeResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.CreateGroupNotice(ctx, in, opts...)
}

// DeleteGroupNotice ???????????????
func (m *defaultGroupService) DeleteGroupNotice(ctx context.Context, in *DeleteGroupNoticeReq, opts ...grpc.CallOption) (*DeleteGroupNoticeResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.DeleteGroupNotice(ctx, in, opts...)
}

// GetGroupNoticeList ?????????????????????
func (m *defaultGroupService) GetGroupNoticeList(ctx context.Context, in *GetGroupNoticeListReq, opts ...grpc.CallOption) (*GetGroupNoticeListResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupNoticeList(ctx, in, opts...)
}

// SetGroupMemberInfo ?????????????????????
func (m *defaultGroupService) SetGroupMemberInfo(ctx context.Context, in *SetGroupMemberInfoReq, opts ...grpc.CallOption) (*SetGroupMemberInfoResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.SetGroupMemberInfo(ctx, in, opts...)
}

// GetGroupMemberInfo ?????????????????????
func (m *defaultGroupService) GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoReq, opts ...grpc.CallOption) (*GetGroupMemberInfoResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupMemberInfo(ctx, in, opts...)
}

// EditGroupInfo ???????????????
func (m *defaultGroupService) EditGroupInfo(ctx context.Context, in *EditGroupInfoReq, opts ...grpc.CallOption) (*EditGroupInfoResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.EditGroupInfo(ctx, in, opts...)
}

// SetGroupSetting ???????????????
func (m *defaultGroupService) SetGroupSetting(ctx context.Context, in *SetGroupSettingReq, opts ...grpc.CallOption) (*SetGroupSettingResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.SetGroupSetting(ctx, in, opts...)
}

// GetGroupSetting ???????????????
func (m *defaultGroupService) GetGroupSetting(ctx context.Context, in *GetGroupSettingReq, opts ...grpc.CallOption) (*GetGroupSettingResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupSetting(ctx, in, opts...)
}

// TransferGroupOwner ????????????
func (m *defaultGroupService) TransferGroupOwner(ctx context.Context, in *TransferGroupOwnerReq, opts ...grpc.CallOption) (*TransferGroupOwnerResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.TransferGroupOwner(ctx, in, opts...)
}

// SetGroupMemberRole ?????????????????????
func (m *defaultGroupService) SetGroupMemberRole(ctx context.Context, in *SetGroupMemberRoleReq, opts ...grpc.CallOption) (*SetGroupMemberRoleResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.SetGroupMemberRole(ctx, in, opts...)
}

// KickGroupMember ???????????????
func (m *defaultGroupService) KickGroupMember(ctx context.Context, in *KickGroupMemberReq, opts ...grpc.CallOption) (*KickGroupMemberResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.KickGroupMember(ctx, in, opts...)
}

// QuitGroup ????????????
func (m *defaultGroupService) QuitGroup(ctx context.Context, in *QuitGroupReq, opts ...grpc.CallOption) (*QuitGroupResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.QuitGroup(ctx, in, opts...)
}

// BanGroupMember ???????????????
func (m *defaultGroupService) BanGroupMember(ctx context.Context, in *BanGroupMemberReq, opts ...grpc.CallOption) (*BanGroupMemberResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.BanGroupMember(ctx, in, opts...)
}

// BanAllGroupMember ?????????????????????
func (m *defaultGroupService) BanAllGroupMember(ctx context.Context, in *BanAllGroupMemberReq, opts ...grpc.CallOption) (*BanAllGroupMemberResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.BanAllGroupMember(ctx, in, opts...)
}

// UnbanGroupMember ?????????????????????
func (m *defaultGroupService) UnbanGroupMember(ctx context.Context, in *UnbanGroupMemberReq, opts ...grpc.CallOption) (*UnbanGroupMemberResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.UnbanGroupMember(ctx, in, opts...)
}

// UnbanAllGroupMember ???????????????????????????
func (m *defaultGroupService) UnbanAllGroupMember(ctx context.Context, in *UnbanAllGroupMemberReq, opts ...grpc.CallOption) (*UnbanAllGroupMemberResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.UnbanAllGroupMember(ctx, in, opts...)
}

// GetGroupMemberList ?????????????????????
func (m *defaultGroupService) GetGroupMemberList(ctx context.Context, in *GetGroupMemberListReq, opts ...grpc.CallOption) (*GetGroupMemberListResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupMemberList(ctx, in, opts...)
}

// DismissGroup ????????????
func (m *defaultGroupService) DismissGroup(ctx context.Context, in *DismissGroupReq, opts ...grpc.CallOption) (*DismissGroupResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.DismissGroup(ctx, in, opts...)
}

// SetGroupMsgNotifyType ???????????????????????????
func (m *defaultGroupService) SetGroupMsgNotifyType(ctx context.Context, in *SetGroupMsgNotifyTypeReq, opts ...grpc.CallOption) (*SetGroupMsgNotifyTypeResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.SetGroupMsgNotifyType(ctx, in, opts...)
}

// GetMyGroupList ????????????????????????
func (m *defaultGroupService) GetMyGroupList(ctx context.Context, in *GetMyGroupListReq, opts ...grpc.CallOption) (*GetMyGroupListResp, error) {
	client := pb.NewGroupServiceClient(m.cli.Conn())
	return client.GetMyGroupList(ctx, in, opts...)
}
