// Code generated by goctl. DO NOT EDIT.
// Source: conversation.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/logic/friendservice"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type FriendServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedFriendServiceServer
}

func NewFriendServiceServer(svcCtx *svc.ServiceContext) *FriendServiceServer {
	return &FriendServiceServer{
		svcCtx: svcCtx,
	}
}

// FriendApply 添加好友
func (s *FriendServiceServer) FriendApply(ctx context.Context, in *pb.FriendApplyReq) (*pb.FriendApplyResp, error) {
	l := friendservicelogic.NewFriendApplyLogic(ctx, s.svcCtx)
	return l.FriendApply(in)
}

// FriendApplyHandle 处理好友申请
func (s *FriendServiceServer) FriendApplyHandle(ctx context.Context, in *pb.FriendApplyHandleReq) (*pb.FriendApplyHandleResp, error) {
	l := friendservicelogic.NewFriendApplyHandleLogic(ctx, s.svcCtx)
	return l.FriendApplyHandle(in)
}

// ListFriendApply 列出好友申请
func (s *FriendServiceServer) ListFriendApply(ctx context.Context, in *pb.ListFriendApplyReq) (*pb.ListFriendApplyResp, error) {
	l := friendservicelogic.NewListFriendApplyLogic(ctx, s.svcCtx)
	return l.ListFriendApply(in)
}

// CountFriend 统计好友数量
func (s *FriendServiceServer) CountFriend(ctx context.Context, in *pb.CountFriendReq) (*pb.CountFriendResp, error) {
	l := friendservicelogic.NewCountFriendLogic(ctx, s.svcCtx)
	return l.CountFriend(in)
}
