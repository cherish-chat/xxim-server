// Code generated by goctl. DO NOT EDIT!
// Source: dispatch.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/dispatch/internal/logic"
	"github.com/cherish-chat/xxim-server/app/dispatch/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type DispatchServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedDispatchServiceServer
}

func NewDispatchServiceServer(svcCtx *svc.ServiceContext) *DispatchServiceServer {
	return &DispatchServiceServer{
		svcCtx: svcCtx,
	}
}

// DispatchOnlineCallback 上线回调
func (s *DispatchServiceServer) DispatchOnlineCallback(ctx context.Context, in *pb.DispatchOnlineCallbackReq) (*pb.DispatchOnlineCallbackResp, error) {
	l := logic.NewDispatchOnlineCallbackLogic(ctx, s.svcCtx)
	return l.DispatchOnlineCallback(in)
}

// BeforeConnect 服务端连接前的回调
func (s *DispatchServiceServer) BeforeConnect(ctx context.Context, in *pb.BeforeConnectReq) (*pb.BeforeConnectResp, error) {
	l := logic.NewBeforeConnectLogic(ctx, s.svcCtx)
	return l.BeforeConnect(in)
}