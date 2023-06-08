// Code generated by goctl. DO NOT EDIT!
// Source: dispatch.proto

package dispatchservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BeforeConnectReq           = pb.BeforeConnectReq
	BeforeConnectResp          = pb.BeforeConnectResp
	DispatchOnlineCallbackReq  = pb.DispatchOnlineCallbackReq
	DispatchOnlineCallbackResp = pb.DispatchOnlineCallbackResp

	DispatchService interface {
		// DispatchOnlineCallback 上线回调
		DispatchOnlineCallback(ctx context.Context, in *DispatchOnlineCallbackReq, opts ...grpc.CallOption) (*DispatchOnlineCallbackResp, error)
		// BeforeConnect 服务端连接前的回调
		BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error)
	}

	defaultDispatchService struct {
		cli zrpc.Client
	}
)

func NewDispatchService(cli zrpc.Client) DispatchService {
	return &defaultDispatchService{
		cli: cli,
	}
}

// DispatchOnlineCallback 上线回调
func (m *defaultDispatchService) DispatchOnlineCallback(ctx context.Context, in *DispatchOnlineCallbackReq, opts ...grpc.CallOption) (*DispatchOnlineCallbackResp, error) {
	client := pb.NewDispatchServiceClient(m.cli.Conn())
	return client.DispatchOnlineCallback(ctx, in, opts...)
}

// BeforeConnect 服务端连接前的回调
func (m *defaultDispatchService) BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error) {
	client := pb.NewDispatchServiceClient(m.cli.Conn())
	return client.BeforeConnect(ctx, in, opts...)
}