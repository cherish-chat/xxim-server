// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/logic/callbackservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type CallbackServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedCallbackServiceServer
}

func NewCallbackServiceServer(svcCtx *svc.ServiceContext) *CallbackServiceServer {
	return &CallbackServiceServer{
		svcCtx: svcCtx,
	}
}

// UserAfterOnline 用户上线回调
func (s *CallbackServiceServer) UserAfterOnline(ctx context.Context, in *pb.UserAfterOnlineReq) (*pb.UserAfterOnlineResp, error) {
	l := callbackservicelogic.NewUserAfterOnlineLogic(ctx, s.svcCtx)
	return l.UserAfterOnline(in)
}

// UserAfterOffline 用户下线回调
func (s *CallbackServiceServer) UserAfterOffline(ctx context.Context, in *pb.UserAfterOfflineReq) (*pb.UserAfterOfflineResp, error) {
	l := callbackservicelogic.NewUserAfterOfflineLogic(ctx, s.svcCtx)
	return l.UserAfterOffline(in)
}

// UserBeforeConnect 用户连接前的回调
func (s *CallbackServiceServer) UserBeforeConnect(ctx context.Context, in *pb.UserBeforeConnectReq) (*pb.UserBeforeConnectResp, error) {
	l := callbackservicelogic.NewUserBeforeConnectLogic(ctx, s.svcCtx)
	return l.UserBeforeConnect(in)
}

// UserBeforeRequest 用户请求前的回调
func (s *CallbackServiceServer) UserBeforeRequest(ctx context.Context, in *pb.UserBeforeRequestReq) (*pb.UserBeforeRequestResp, error) {
	l := callbackservicelogic.NewUserBeforeRequestLogic(ctx, s.svcCtx)
	return l.UserBeforeRequest(in)
}

// UserAfterKeepAlive 用户保活回调
func (s *CallbackServiceServer) UserAfterKeepAlive(ctx context.Context, in *pb.UserAfterKeepAliveReq) (*pb.UserAfterKeepAliveResp, error) {
	l := callbackservicelogic.NewUserAfterKeepAliveLogic(ctx, s.svcCtx)
	return l.UserAfterKeepAlive(in)
}