package mgmtservice

import (
	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
)

func MustLoadConfig(rpcAddr string, name string, dest interface{}) MgmtService {
	// 创建一个客户端
	client := NewMgmtService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{rpcAddr},
	}))
	resp, err := client.GetServerConfig(context.Background(), &GetServerConfigReq{ServerName: name})
	if err != nil {
		log.Fatalf("获取服务配置失败: %s", err)
	}
	// 打印配置信息
	log.Printf("配置信息: \n%s", resp.Config)
	if err := conf.LoadFromJsonBytes(resp.Config, dest); err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	return client
}
