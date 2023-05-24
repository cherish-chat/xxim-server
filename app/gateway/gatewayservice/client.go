package gatewayservice

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

func MustNewGatewayService(config zrpc.RpcClientConf) GatewayService {
	client, err := zrpc.NewClient(config)
	if err != nil {
		logx.Errorf("zrpc.NewClient error: %v", err)
		os.Exit(1)
	}
	return NewGatewayService(client)
}
