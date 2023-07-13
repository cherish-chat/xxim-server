// Code generated by goctl. DO NOT EDIT.
// Source: conversation.peer.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/logic/groupservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"
)

type GroupServiceServer struct {
	svcCtx *svc.ServiceContext
	peerpb.UnimplementedGroupServiceServer
}

func NewGroupServiceServer(svcCtx *svc.ServiceContext) *GroupServiceServer {
	return &GroupServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *GroupServiceServer) GroupAfterKeepAlive(ctx context.Context, in *peerpb.GroupAfterKeepAliveReq) (*peerpb.GroupAfterKeepAliveResp, error) {
	l := groupservicelogic.NewGroupAfterKeepAliveLogic(ctx, s.svcCtx)
	return l.GroupAfterKeepAlive(in)
}

func (s *GroupServiceServer) GroupAfterOnline(ctx context.Context, in *peerpb.GroupAfterOnlineReq) (*peerpb.GroupAfterOnlineResp, error) {
	l := groupservicelogic.NewGroupAfterOnlineLogic(ctx, s.svcCtx)
	return l.GroupAfterOnline(in)
}

func (s *GroupServiceServer) GroupAfterOffline(ctx context.Context, in *peerpb.GroupAfterOfflineReq) (*peerpb.GroupAfterOfflineResp, error) {
	l := groupservicelogic.NewGroupAfterOfflineLogic(ctx, s.svcCtx)
	return l.GroupAfterOffline(in)
}

func (s *GroupServiceServer) GroupCreate(ctx context.Context, in *peerpb.GroupCreateReq) (*peerpb.GroupCreateResp, error) {
	l := groupservicelogic.NewGroupCreateLogic(ctx, s.svcCtx)
	return l.GroupCreate(in)
}

func (s *GroupServiceServer) CountJoinGroup(ctx context.Context, in *peerpb.CountJoinGroupReq) (*peerpb.CountJoinGroupResp, error) {
	l := groupservicelogic.NewCountJoinGroupLogic(ctx, s.svcCtx)
	return l.CountJoinGroup(in)
}

func (s *GroupServiceServer) CountCreateGroup(ctx context.Context, in *peerpb.CountCreateGroupReq) (*peerpb.CountCreateGroupResp, error) {
	l := groupservicelogic.NewCountCreateGroupLogic(ctx, s.svcCtx)
	return l.CountCreateGroup(in)
}

func (s *GroupServiceServer) ListGroupSubscribers(ctx context.Context, in *peerpb.ListGroupSubscribersReq) (*peerpb.ListGroupSubscribersResp, error) {
	l := groupservicelogic.NewListGroupSubscribersLogic(ctx, s.svcCtx)
	return l.ListGroupSubscribers(in)
}
