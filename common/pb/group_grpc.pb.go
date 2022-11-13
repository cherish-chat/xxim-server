// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: group.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GroupServiceClient is the client API for GroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupServiceClient interface {
	//CreateGroup 创建群聊
	CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateGroupResp, error)
	//GetGroupHome 获取群聊首页
	GetGroupHome(ctx context.Context, in *GetGroupHomeReq, opts ...grpc.CallOption) (*GetGroupHomeResp, error)
	//InviteFriendToGroup 邀请好友加入群聊
	InviteFriendToGroup(ctx context.Context, in *InviteFriendToGroupReq, opts ...grpc.CallOption) (*InviteFriendToGroupResp, error)
	//CreateGroupNotice 创建群公告
	CreateGroupNotice(ctx context.Context, in *CreateGroupNoticeReq, opts ...grpc.CallOption) (*CreateGroupNoticeResp, error)
	//DeleteGroupNotice 删除群公告
	DeleteGroupNotice(ctx context.Context, in *DeleteGroupNoticeReq, opts ...grpc.CallOption) (*DeleteGroupNoticeResp, error)
	//GetGroupNoticeList 获取群公告列表
	GetGroupNoticeList(ctx context.Context, in *GetGroupNoticeListReq, opts ...grpc.CallOption) (*GetGroupNoticeListResp, error)
	//SetGroupMemberInfo 设置群成员信息
	SetGroupMemberInfo(ctx context.Context, in *SetGroupMemberInfoReq, opts ...grpc.CallOption) (*SetGroupMemberInfoResp, error)
	//GetGroupMemberInfo 获取群成员信息
	GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoReq, opts ...grpc.CallOption) (*GetGroupMemberInfoResp, error)
	//EditGroupInfo 编辑群信息
	EditGroupInfo(ctx context.Context, in *EditGroupInfoReq, opts ...grpc.CallOption) (*EditGroupInfoResp, error)
	//SetGroupSetting 设置群设置
	SetGroupSetting(ctx context.Context, in *SetGroupSettingReq, opts ...grpc.CallOption) (*SetGroupSettingResp, error)
	//GetGroupSetting 获取群设置
	GetGroupSetting(ctx context.Context, in *GetGroupSettingReq, opts ...grpc.CallOption) (*GetGroupSettingResp, error)
	//TransferGroupOwner 转让群主
	TransferGroupOwner(ctx context.Context, in *TransferGroupOwnerReq, opts ...grpc.CallOption) (*TransferGroupOwnerResp, error)
	//SetGroupMemberRole 设置群成员角色
	SetGroupMemberRole(ctx context.Context, in *SetGroupMemberRoleReq, opts ...grpc.CallOption) (*SetGroupMemberRoleResp, error)
	//KickGroupMember 踢出群成员
	KickGroupMember(ctx context.Context, in *KickGroupMemberReq, opts ...grpc.CallOption) (*KickGroupMemberResp, error)
	//QuitGroup 退出群聊
	QuitGroup(ctx context.Context, in *QuitGroupReq, opts ...grpc.CallOption) (*QuitGroupResp, error)
	//BanGroupMember 禁言群成员
	BanGroupMember(ctx context.Context, in *BanGroupMemberReq, opts ...grpc.CallOption) (*BanGroupMemberResp, error)
	//BanAllGroupMember 禁言全部群成员
	BanAllGroupMember(ctx context.Context, in *BanAllGroupMemberReq, opts ...grpc.CallOption) (*BanAllGroupMemberResp, error)
	//UnbanGroupMember 解除禁言群成员
	UnbanGroupMember(ctx context.Context, in *UnbanGroupMemberReq, opts ...grpc.CallOption) (*UnbanGroupMemberResp, error)
	//UnbanAllGroupMember 解除禁言全部群成员
	UnbanAllGroupMember(ctx context.Context, in *UnbanAllGroupMemberReq, opts ...grpc.CallOption) (*UnbanAllGroupMemberResp, error)
	//GetGroupMemberList 获取群成员列表
	GetGroupMemberList(ctx context.Context, in *GetGroupMemberListReq, opts ...grpc.CallOption) (*GetGroupMemberListResp, error)
	//DismissGroup 解散群聊
	DismissGroup(ctx context.Context, in *DismissGroupReq, opts ...grpc.CallOption) (*DismissGroupResp, error)
	//SetGroupMsgNotifyType 设置群消息通知选项
	SetGroupMsgNotifyType(ctx context.Context, in *SetGroupMsgNotifyTypeReq, opts ...grpc.CallOption) (*SetGroupMsgNotifyTypeResp, error)
	//GetMyGroupList 获取我的群聊列表
	GetMyGroupList(ctx context.Context, in *GetMyGroupListReq, opts ...grpc.CallOption) (*GetMyGroupListResp, error)
}

type groupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupServiceClient(cc grpc.ClientConnInterface) GroupServiceClient {
	return &groupServiceClient{cc}
}

func (c *groupServiceClient) CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateGroupResp, error) {
	out := new(CreateGroupResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/CreateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupHome(ctx context.Context, in *GetGroupHomeReq, opts ...grpc.CallOption) (*GetGroupHomeResp, error) {
	out := new(GetGroupHomeResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetGroupHome", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) InviteFriendToGroup(ctx context.Context, in *InviteFriendToGroupReq, opts ...grpc.CallOption) (*InviteFriendToGroupResp, error) {
	out := new(InviteFriendToGroupResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/InviteFriendToGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) CreateGroupNotice(ctx context.Context, in *CreateGroupNoticeReq, opts ...grpc.CallOption) (*CreateGroupNoticeResp, error) {
	out := new(CreateGroupNoticeResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/CreateGroupNotice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) DeleteGroupNotice(ctx context.Context, in *DeleteGroupNoticeReq, opts ...grpc.CallOption) (*DeleteGroupNoticeResp, error) {
	out := new(DeleteGroupNoticeResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/DeleteGroupNotice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupNoticeList(ctx context.Context, in *GetGroupNoticeListReq, opts ...grpc.CallOption) (*GetGroupNoticeListResp, error) {
	out := new(GetGroupNoticeListResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetGroupNoticeList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SetGroupMemberInfo(ctx context.Context, in *SetGroupMemberInfoReq, opts ...grpc.CallOption) (*SetGroupMemberInfoResp, error) {
	out := new(SetGroupMemberInfoResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/SetGroupMemberInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoReq, opts ...grpc.CallOption) (*GetGroupMemberInfoResp, error) {
	out := new(GetGroupMemberInfoResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetGroupMemberInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) EditGroupInfo(ctx context.Context, in *EditGroupInfoReq, opts ...grpc.CallOption) (*EditGroupInfoResp, error) {
	out := new(EditGroupInfoResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/EditGroupInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SetGroupSetting(ctx context.Context, in *SetGroupSettingReq, opts ...grpc.CallOption) (*SetGroupSettingResp, error) {
	out := new(SetGroupSettingResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/SetGroupSetting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupSetting(ctx context.Context, in *GetGroupSettingReq, opts ...grpc.CallOption) (*GetGroupSettingResp, error) {
	out := new(GetGroupSettingResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetGroupSetting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) TransferGroupOwner(ctx context.Context, in *TransferGroupOwnerReq, opts ...grpc.CallOption) (*TransferGroupOwnerResp, error) {
	out := new(TransferGroupOwnerResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/TransferGroupOwner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SetGroupMemberRole(ctx context.Context, in *SetGroupMemberRoleReq, opts ...grpc.CallOption) (*SetGroupMemberRoleResp, error) {
	out := new(SetGroupMemberRoleResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/SetGroupMemberRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) KickGroupMember(ctx context.Context, in *KickGroupMemberReq, opts ...grpc.CallOption) (*KickGroupMemberResp, error) {
	out := new(KickGroupMemberResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/KickGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) QuitGroup(ctx context.Context, in *QuitGroupReq, opts ...grpc.CallOption) (*QuitGroupResp, error) {
	out := new(QuitGroupResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/QuitGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) BanGroupMember(ctx context.Context, in *BanGroupMemberReq, opts ...grpc.CallOption) (*BanGroupMemberResp, error) {
	out := new(BanGroupMemberResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/BanGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) BanAllGroupMember(ctx context.Context, in *BanAllGroupMemberReq, opts ...grpc.CallOption) (*BanAllGroupMemberResp, error) {
	out := new(BanAllGroupMemberResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/BanAllGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) UnbanGroupMember(ctx context.Context, in *UnbanGroupMemberReq, opts ...grpc.CallOption) (*UnbanGroupMemberResp, error) {
	out := new(UnbanGroupMemberResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/UnbanGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) UnbanAllGroupMember(ctx context.Context, in *UnbanAllGroupMemberReq, opts ...grpc.CallOption) (*UnbanAllGroupMemberResp, error) {
	out := new(UnbanAllGroupMemberResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/UnbanAllGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupMemberList(ctx context.Context, in *GetGroupMemberListReq, opts ...grpc.CallOption) (*GetGroupMemberListResp, error) {
	out := new(GetGroupMemberListResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetGroupMemberList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) DismissGroup(ctx context.Context, in *DismissGroupReq, opts ...grpc.CallOption) (*DismissGroupResp, error) {
	out := new(DismissGroupResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/DismissGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SetGroupMsgNotifyType(ctx context.Context, in *SetGroupMsgNotifyTypeReq, opts ...grpc.CallOption) (*SetGroupMsgNotifyTypeResp, error) {
	out := new(SetGroupMsgNotifyTypeResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/SetGroupMsgNotifyType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetMyGroupList(ctx context.Context, in *GetMyGroupListReq, opts ...grpc.CallOption) (*GetMyGroupListResp, error) {
	out := new(GetMyGroupListResp)
	err := c.cc.Invoke(ctx, "/pb.groupService/GetMyGroupList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupServiceServer is the server API for GroupService service.
// All implementations must embed UnimplementedGroupServiceServer
// for forward compatibility
type GroupServiceServer interface {
	//CreateGroup 创建群聊
	CreateGroup(context.Context, *CreateGroupReq) (*CreateGroupResp, error)
	//GetGroupHome 获取群聊首页
	GetGroupHome(context.Context, *GetGroupHomeReq) (*GetGroupHomeResp, error)
	//InviteFriendToGroup 邀请好友加入群聊
	InviteFriendToGroup(context.Context, *InviteFriendToGroupReq) (*InviteFriendToGroupResp, error)
	//CreateGroupNotice 创建群公告
	CreateGroupNotice(context.Context, *CreateGroupNoticeReq) (*CreateGroupNoticeResp, error)
	//DeleteGroupNotice 删除群公告
	DeleteGroupNotice(context.Context, *DeleteGroupNoticeReq) (*DeleteGroupNoticeResp, error)
	//GetGroupNoticeList 获取群公告列表
	GetGroupNoticeList(context.Context, *GetGroupNoticeListReq) (*GetGroupNoticeListResp, error)
	//SetGroupMemberInfo 设置群成员信息
	SetGroupMemberInfo(context.Context, *SetGroupMemberInfoReq) (*SetGroupMemberInfoResp, error)
	//GetGroupMemberInfo 获取群成员信息
	GetGroupMemberInfo(context.Context, *GetGroupMemberInfoReq) (*GetGroupMemberInfoResp, error)
	//EditGroupInfo 编辑群信息
	EditGroupInfo(context.Context, *EditGroupInfoReq) (*EditGroupInfoResp, error)
	//SetGroupSetting 设置群设置
	SetGroupSetting(context.Context, *SetGroupSettingReq) (*SetGroupSettingResp, error)
	//GetGroupSetting 获取群设置
	GetGroupSetting(context.Context, *GetGroupSettingReq) (*GetGroupSettingResp, error)
	//TransferGroupOwner 转让群主
	TransferGroupOwner(context.Context, *TransferGroupOwnerReq) (*TransferGroupOwnerResp, error)
	//SetGroupMemberRole 设置群成员角色
	SetGroupMemberRole(context.Context, *SetGroupMemberRoleReq) (*SetGroupMemberRoleResp, error)
	//KickGroupMember 踢出群成员
	KickGroupMember(context.Context, *KickGroupMemberReq) (*KickGroupMemberResp, error)
	//QuitGroup 退出群聊
	QuitGroup(context.Context, *QuitGroupReq) (*QuitGroupResp, error)
	//BanGroupMember 禁言群成员
	BanGroupMember(context.Context, *BanGroupMemberReq) (*BanGroupMemberResp, error)
	//BanAllGroupMember 禁言全部群成员
	BanAllGroupMember(context.Context, *BanAllGroupMemberReq) (*BanAllGroupMemberResp, error)
	//UnbanGroupMember 解除禁言群成员
	UnbanGroupMember(context.Context, *UnbanGroupMemberReq) (*UnbanGroupMemberResp, error)
	//UnbanAllGroupMember 解除禁言全部群成员
	UnbanAllGroupMember(context.Context, *UnbanAllGroupMemberReq) (*UnbanAllGroupMemberResp, error)
	//GetGroupMemberList 获取群成员列表
	GetGroupMemberList(context.Context, *GetGroupMemberListReq) (*GetGroupMemberListResp, error)
	//DismissGroup 解散群聊
	DismissGroup(context.Context, *DismissGroupReq) (*DismissGroupResp, error)
	//SetGroupMsgNotifyType 设置群消息通知选项
	SetGroupMsgNotifyType(context.Context, *SetGroupMsgNotifyTypeReq) (*SetGroupMsgNotifyTypeResp, error)
	//GetMyGroupList 获取我的群聊列表
	GetMyGroupList(context.Context, *GetMyGroupListReq) (*GetMyGroupListResp, error)
	mustEmbedUnimplementedGroupServiceServer()
}

// UnimplementedGroupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGroupServiceServer struct {
}

func (UnimplementedGroupServiceServer) CreateGroup(context.Context, *CreateGroupReq) (*CreateGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupHome(context.Context, *GetGroupHomeReq) (*GetGroupHomeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupHome not implemented")
}
func (UnimplementedGroupServiceServer) InviteFriendToGroup(context.Context, *InviteFriendToGroupReq) (*InviteFriendToGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteFriendToGroup not implemented")
}
func (UnimplementedGroupServiceServer) CreateGroupNotice(context.Context, *CreateGroupNoticeReq) (*CreateGroupNoticeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroupNotice not implemented")
}
func (UnimplementedGroupServiceServer) DeleteGroupNotice(context.Context, *DeleteGroupNoticeReq) (*DeleteGroupNoticeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroupNotice not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupNoticeList(context.Context, *GetGroupNoticeListReq) (*GetGroupNoticeListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupNoticeList not implemented")
}
func (UnimplementedGroupServiceServer) SetGroupMemberInfo(context.Context, *SetGroupMemberInfoReq) (*SetGroupMemberInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGroupMemberInfo not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupMemberInfo(context.Context, *GetGroupMemberInfoReq) (*GetGroupMemberInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMemberInfo not implemented")
}
func (UnimplementedGroupServiceServer) EditGroupInfo(context.Context, *EditGroupInfoReq) (*EditGroupInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditGroupInfo not implemented")
}
func (UnimplementedGroupServiceServer) SetGroupSetting(context.Context, *SetGroupSettingReq) (*SetGroupSettingResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGroupSetting not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupSetting(context.Context, *GetGroupSettingReq) (*GetGroupSettingResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupSetting not implemented")
}
func (UnimplementedGroupServiceServer) TransferGroupOwner(context.Context, *TransferGroupOwnerReq) (*TransferGroupOwnerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferGroupOwner not implemented")
}
func (UnimplementedGroupServiceServer) SetGroupMemberRole(context.Context, *SetGroupMemberRoleReq) (*SetGroupMemberRoleResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGroupMemberRole not implemented")
}
func (UnimplementedGroupServiceServer) KickGroupMember(context.Context, *KickGroupMemberReq) (*KickGroupMemberResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KickGroupMember not implemented")
}
func (UnimplementedGroupServiceServer) QuitGroup(context.Context, *QuitGroupReq) (*QuitGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuitGroup not implemented")
}
func (UnimplementedGroupServiceServer) BanGroupMember(context.Context, *BanGroupMemberReq) (*BanGroupMemberResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BanGroupMember not implemented")
}
func (UnimplementedGroupServiceServer) BanAllGroupMember(context.Context, *BanAllGroupMemberReq) (*BanAllGroupMemberResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BanAllGroupMember not implemented")
}
func (UnimplementedGroupServiceServer) UnbanGroupMember(context.Context, *UnbanGroupMemberReq) (*UnbanGroupMemberResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbanGroupMember not implemented")
}
func (UnimplementedGroupServiceServer) UnbanAllGroupMember(context.Context, *UnbanAllGroupMemberReq) (*UnbanAllGroupMemberResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbanAllGroupMember not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupMemberList(context.Context, *GetGroupMemberListReq) (*GetGroupMemberListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMemberList not implemented")
}
func (UnimplementedGroupServiceServer) DismissGroup(context.Context, *DismissGroupReq) (*DismissGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DismissGroup not implemented")
}
func (UnimplementedGroupServiceServer) SetGroupMsgNotifyType(context.Context, *SetGroupMsgNotifyTypeReq) (*SetGroupMsgNotifyTypeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGroupMsgNotifyType not implemented")
}
func (UnimplementedGroupServiceServer) GetMyGroupList(context.Context, *GetMyGroupListReq) (*GetMyGroupListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMyGroupList not implemented")
}
func (UnimplementedGroupServiceServer) mustEmbedUnimplementedGroupServiceServer() {}

// UnsafeGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupServiceServer will
// result in compilation errors.
type UnsafeGroupServiceServer interface {
	mustEmbedUnimplementedGroupServiceServer()
}

func RegisterGroupServiceServer(s grpc.ServiceRegistrar, srv GroupServiceServer) {
	s.RegisterService(&GroupService_ServiceDesc, srv)
}

func _GroupService_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/CreateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).CreateGroup(ctx, req.(*CreateGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupHome_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupHomeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupHome(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetGroupHome",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupHome(ctx, req.(*GetGroupHomeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_InviteFriendToGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteFriendToGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).InviteFriendToGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/InviteFriendToGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).InviteFriendToGroup(ctx, req.(*InviteFriendToGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_CreateGroupNotice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupNoticeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).CreateGroupNotice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/CreateGroupNotice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).CreateGroupNotice(ctx, req.(*CreateGroupNoticeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_DeleteGroupNotice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGroupNoticeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).DeleteGroupNotice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/DeleteGroupNotice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).DeleteGroupNotice(ctx, req.(*DeleteGroupNoticeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupNoticeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupNoticeListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupNoticeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetGroupNoticeList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupNoticeList(ctx, req.(*GetGroupNoticeListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SetGroupMemberInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGroupMemberInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SetGroupMemberInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/SetGroupMemberInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SetGroupMemberInfo(ctx, req.(*SetGroupMemberInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupMemberInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMemberInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupMemberInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetGroupMemberInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupMemberInfo(ctx, req.(*GetGroupMemberInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_EditGroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditGroupInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).EditGroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/EditGroupInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).EditGroupInfo(ctx, req.(*EditGroupInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SetGroupSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGroupSettingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SetGroupSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/SetGroupSetting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SetGroupSetting(ctx, req.(*SetGroupSettingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupSettingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetGroupSetting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupSetting(ctx, req.(*GetGroupSettingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_TransferGroupOwner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferGroupOwnerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).TransferGroupOwner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/TransferGroupOwner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).TransferGroupOwner(ctx, req.(*TransferGroupOwnerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SetGroupMemberRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGroupMemberRoleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SetGroupMemberRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/SetGroupMemberRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SetGroupMemberRole(ctx, req.(*SetGroupMemberRoleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_KickGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KickGroupMemberReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).KickGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/KickGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).KickGroupMember(ctx, req.(*KickGroupMemberReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_QuitGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuitGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).QuitGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/QuitGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).QuitGroup(ctx, req.(*QuitGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_BanGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BanGroupMemberReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).BanGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/BanGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).BanGroupMember(ctx, req.(*BanGroupMemberReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_BanAllGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BanAllGroupMemberReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).BanAllGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/BanAllGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).BanAllGroupMember(ctx, req.(*BanAllGroupMemberReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_UnbanGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbanGroupMemberReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).UnbanGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/UnbanGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).UnbanGroupMember(ctx, req.(*UnbanGroupMemberReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_UnbanAllGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbanAllGroupMemberReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).UnbanAllGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/UnbanAllGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).UnbanAllGroupMember(ctx, req.(*UnbanAllGroupMemberReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupMemberList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMemberListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupMemberList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetGroupMemberList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupMemberList(ctx, req.(*GetGroupMemberListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_DismissGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DismissGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).DismissGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/DismissGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).DismissGroup(ctx, req.(*DismissGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SetGroupMsgNotifyType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGroupMsgNotifyTypeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SetGroupMsgNotifyType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/SetGroupMsgNotifyType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SetGroupMsgNotifyType(ctx, req.(*SetGroupMsgNotifyTypeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetMyGroupList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMyGroupListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetMyGroupList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.groupService/GetMyGroupList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetMyGroupList(ctx, req.(*GetMyGroupListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// GroupService_ServiceDesc is the grpc.ServiceDesc for GroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.groupService",
	HandlerType: (*GroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGroup",
			Handler:    _GroupService_CreateGroup_Handler,
		},
		{
			MethodName: "GetGroupHome",
			Handler:    _GroupService_GetGroupHome_Handler,
		},
		{
			MethodName: "InviteFriendToGroup",
			Handler:    _GroupService_InviteFriendToGroup_Handler,
		},
		{
			MethodName: "CreateGroupNotice",
			Handler:    _GroupService_CreateGroupNotice_Handler,
		},
		{
			MethodName: "DeleteGroupNotice",
			Handler:    _GroupService_DeleteGroupNotice_Handler,
		},
		{
			MethodName: "GetGroupNoticeList",
			Handler:    _GroupService_GetGroupNoticeList_Handler,
		},
		{
			MethodName: "SetGroupMemberInfo",
			Handler:    _GroupService_SetGroupMemberInfo_Handler,
		},
		{
			MethodName: "GetGroupMemberInfo",
			Handler:    _GroupService_GetGroupMemberInfo_Handler,
		},
		{
			MethodName: "EditGroupInfo",
			Handler:    _GroupService_EditGroupInfo_Handler,
		},
		{
			MethodName: "SetGroupSetting",
			Handler:    _GroupService_SetGroupSetting_Handler,
		},
		{
			MethodName: "GetGroupSetting",
			Handler:    _GroupService_GetGroupSetting_Handler,
		},
		{
			MethodName: "TransferGroupOwner",
			Handler:    _GroupService_TransferGroupOwner_Handler,
		},
		{
			MethodName: "SetGroupMemberRole",
			Handler:    _GroupService_SetGroupMemberRole_Handler,
		},
		{
			MethodName: "KickGroupMember",
			Handler:    _GroupService_KickGroupMember_Handler,
		},
		{
			MethodName: "QuitGroup",
			Handler:    _GroupService_QuitGroup_Handler,
		},
		{
			MethodName: "BanGroupMember",
			Handler:    _GroupService_BanGroupMember_Handler,
		},
		{
			MethodName: "BanAllGroupMember",
			Handler:    _GroupService_BanAllGroupMember_Handler,
		},
		{
			MethodName: "UnbanGroupMember",
			Handler:    _GroupService_UnbanGroupMember_Handler,
		},
		{
			MethodName: "UnbanAllGroupMember",
			Handler:    _GroupService_UnbanAllGroupMember_Handler,
		},
		{
			MethodName: "GetGroupMemberList",
			Handler:    _GroupService_GetGroupMemberList_Handler,
		},
		{
			MethodName: "DismissGroup",
			Handler:    _GroupService_DismissGroup_Handler,
		},
		{
			MethodName: "SetGroupMsgNotifyType",
			Handler:    _GroupService_SetGroupMsgNotifyType_Handler,
		},
		{
			MethodName: "GetMyGroupList",
			Handler:    _GroupService_GetMyGroupList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "group.proto",
}