// Code generated by goctl. DO NOT EDIT!
// Source: xx.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/logic"
	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type XxServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedXxServiceServer
}

func NewXxServiceServer(svcCtx *svc.ServiceContext) *XxServiceServer {
	return &XxServiceServer{
		svcCtx: svcCtx,
	}
}

// Register 注册用户
func (s *XxServiceServer) Register(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

// GetUser 获取用户信息
func (s *XxServiceServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	l := logic.NewGetUserLogic(ctx, s.svcCtx)
	return l.GetUser(in)
}

// Login 登录
func (s *XxServiceServer) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

// UpdateUser 更新用户信息
func (s *XxServiceServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svcCtx)
	return l.UpdateUser(in)
}

// SearchUser 搜索用户
func (s *XxServiceServer) SearchUser(ctx context.Context, in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	l := logic.NewSearchUserLogic(ctx, s.svcCtx)
	return l.SearchUser(in)
}

// GetFriend 获取好友
func (s *XxServiceServer) GetFriend(ctx context.Context, in *pb.GetFriendReq) (*pb.GetFriendResp, error) {
	l := logic.NewGetFriendLogic(ctx, s.svcCtx)
	return l.GetFriend(in)
}

// CreateFriend 添加好友
func (s *XxServiceServer) CreateFriend(ctx context.Context, in *pb.CreateFriendReq) (*pb.CreateFriendResp, error) {
	l := logic.NewCreateFriendLogic(ctx, s.svcCtx)
	return l.CreateFriend(in)
}

// UpdateFriend 更新好友
func (s *XxServiceServer) UpdateFriend(ctx context.Context, in *pb.UpdateFriendReq) (*pb.UpdateFriendResp, error) {
	l := logic.NewUpdateFriendLogic(ctx, s.svcCtx)
	return l.UpdateFriend(in)
}

// DeleteFriend 删除好友
func (s *XxServiceServer) DeleteFriend(ctx context.Context, in *pb.DeleteFriendReq) (*pb.DeleteFriendResp, error) {
	l := logic.NewDeleteFriendLogic(ctx, s.svcCtx)
	return l.DeleteFriend(in)
}

// CreateGroup 创建群组
func (s *XxServiceServer) CreateGroup(ctx context.Context, in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	l := logic.NewCreateGroupLogic(ctx, s.svcCtx)
	return l.CreateGroup(in)
}

// GetGroup 获取群组
func (s *XxServiceServer) GetGroup(ctx context.Context, in *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	l := logic.NewGetGroupLogic(ctx, s.svcCtx)
	return l.GetGroup(in)
}

// InviteGroup 邀请加入群组
func (s *XxServiceServer) InviteGroup(ctx context.Context, in *pb.InviteGroupReq) (*pb.InviteGroupResp, error) {
	l := logic.NewInviteGroupLogic(ctx, s.svcCtx)
	return l.InviteGroup(in)
}

// UpdateGroup 更新群组
func (s *XxServiceServer) UpdateGroup(ctx context.Context, in *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	l := logic.NewUpdateGroupLogic(ctx, s.svcCtx)
	return l.UpdateGroup(in)
}

// QuitGroupMember 退出群组
func (s *XxServiceServer) QuitGroupMember(ctx context.Context, in *pb.QuitGroupMemberReq) (*pb.QuitGroupMemberResp, error) {
	l := logic.NewQuitGroupMemberLogic(ctx, s.svcCtx)
	return l.QuitGroupMember(in)
}

// KickGroupMember 踢出群组
func (s *XxServiceServer) KickGroupMember(ctx context.Context, in *pb.KickGroupMemberReq) (*pb.KickGroupMemberResp, error) {
	l := logic.NewKickGroupMemberLogic(ctx, s.svcCtx)
	return l.KickGroupMember(in)
}

// ClearGroupMember 清空群组成员
func (s *XxServiceServer) ClearGroupMember(ctx context.Context, in *pb.ClearGroupMemberReq) (*pb.ClearGroupMemberResp, error) {
	l := logic.NewClearGroupMemberLogic(ctx, s.svcCtx)
	return l.ClearGroupMember(in)
}
