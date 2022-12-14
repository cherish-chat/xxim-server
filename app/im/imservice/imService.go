// Code generated by goctl. DO NOT EDIT!
// Source: im.proto

package imservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BatchGetUserLatestConnReq  = pb.BatchGetUserLatestConnReq
	BatchGetUserLatestConnResp = pb.BatchGetUserLatestConnResp
	BeforeConnectReq           = pb.BeforeConnectReq
	BeforeConnectResp          = pb.BeforeConnectResp
	GetAppSystemConfigReq      = pb.GetAppSystemConfigReq
	GetAppSystemConfigResp     = pb.GetAppSystemConfigResp
	GetUserLatestConnReq       = pb.GetUserLatestConnReq
	GetUserLatestConnResp      = pb.GetUserLatestConnResp
	ImMQBody                   = pb.ImMQBody
	MsgNotifyOpt               = pb.MsgNotifyOpt

	ImService interface {
		BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error)
		AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		KickUserConn(ctx context.Context, in *KickUserConnReq, opts ...grpc.CallOption) (*KickUserConnResp, error)
		GetUserConn(ctx context.Context, in *GetUserConnReq, opts ...grpc.CallOption) (*GetUserConnResp, error)
		GetUserLatestConn(ctx context.Context, in *GetUserLatestConnReq, opts ...grpc.CallOption) (*GetUserLatestConnResp, error)
		BatchGetUserLatestConn(ctx context.Context, in *BatchGetUserLatestConnReq, opts ...grpc.CallOption) (*BatchGetUserLatestConnResp, error)
		SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error)
		GetAppSystemConfig(ctx context.Context, in *GetAppSystemConfigReq, opts ...grpc.CallOption) (*GetAppSystemConfigResp, error)
	}

	defaultImService struct {
		cli zrpc.Client
	}
)

func NewImService(cli zrpc.Client) ImService {
	return &defaultImService{
		cli: cli,
	}
}

func (m *defaultImService) BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.BeforeConnect(ctx, in, opts...)
}

func (m *defaultImService) AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.AfterConnect(ctx, in, opts...)
}

func (m *defaultImService) AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.AfterDisconnect(ctx, in, opts...)
}

func (m *defaultImService) KickUserConn(ctx context.Context, in *KickUserConnReq, opts ...grpc.CallOption) (*KickUserConnResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.KickUserConn(ctx, in, opts...)
}

func (m *defaultImService) GetUserConn(ctx context.Context, in *GetUserConnReq, opts ...grpc.CallOption) (*GetUserConnResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.GetUserConn(ctx, in, opts...)
}

func (m *defaultImService) GetUserLatestConn(ctx context.Context, in *GetUserLatestConnReq, opts ...grpc.CallOption) (*GetUserLatestConnResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.GetUserLatestConn(ctx, in, opts...)
}

func (m *defaultImService) BatchGetUserLatestConn(ctx context.Context, in *BatchGetUserLatestConnReq, opts ...grpc.CallOption) (*BatchGetUserLatestConnResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.BatchGetUserLatestConn(ctx, in, opts...)
}

func (m *defaultImService) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.SendMsg(ctx, in, opts...)
}

func (m *defaultImService) GetAppSystemConfig(ctx context.Context, in *GetAppSystemConfigReq, opts ...grpc.CallOption) (*GetAppSystemConfigResp, error) {
	client := pb.NewImServiceClient(m.cli.Conn())
	return client.GetAppSystemConfig(ctx, in, opts...)
}
