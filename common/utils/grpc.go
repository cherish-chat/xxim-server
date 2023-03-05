package utils

import (
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ZrpcConfig struct {
}

var Zrpc ZrpcConfig

func (c ZrpcConfig) MaxCallRecvMsgSize() int {
	return 1024 * 1024 * 100 // 100MB
}

func (c ZrpcConfig) Options() []zrpc.ClientOption {
	return []zrpc.ClientOption{
		zrpc.WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(
			c.MaxCallRecvMsgSize(),
		))),
	}
}
