package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/deploy/tentcent/config"
	"github.com/cherish-chat/xxim-server/deploy/tentcent/tdmq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/deploy.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// tdmq
	{
		mgr := tdmq.NewMgr(&tdmq.Config{
			Namespace:   c.Tencent.TDMQ.Namespace,
			ClusterName: c.Tencent.TDMQ.ClusterName,
			SecretId:    c.Tencent.SecretId,
			SecretKey:   c.Tencent.SecretKey,
			Region:      c.Tencent.Region,
		})
		{
			err := mgr.CreateTopic("xxmsg", c.Tencent.TDMQ.Topics.Msg.Partition)
			if err != nil {
				logx.Errorf("create topic xxmsg failed: %v", err)
				panic(err)
			}
			err = mgr.CreateSubscription("xxmsg", "msg")
			if err != nil {
				logx.Errorf("create subscription xxmsg failed: %v", err)
				panic(err)
			}
		}
	}
}
