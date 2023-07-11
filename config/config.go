package config

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
)

type Config struct {
	Etcd   EtcdConfig
	Log    LogConfig
	Mode   string `json:",default=pro,options=dev|test|rt|pre|pro"`
	Jaeger JaegerConfig

	Gateway struct {
		Mode            string   `json:",default=p2p,options=p2p"`
		SignalingServer string   `json:",optional"` // xxx.xxx.xxx:xx
		AppId           string   `json:",optional"`
		AppSecret       string   `json:",optional"`
		StunUrls        []string `json:",optional"`
	}
}

func RpcPort() string {
	// 获取空闲端口
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return fmt.Sprintf("0.0.0.0:%d", l.Addr().(*net.TCPAddr).Port)
		}
	}
	logx.Errorf("get rpc port failed: %v", err)
	panic(err)
}
