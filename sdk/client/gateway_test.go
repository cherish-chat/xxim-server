package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
	"time"
)

var defaultConfig = &Config{
	Endpoints:   []string{"http://127.0.0.1:34500"},
	Encoding:    utils.AnyPtr(pb.EncodingProto_JSON),
	AppId:       "",
	InstallId:   "",
	Platform:    nil,
	DeviceModel: "",
	OsVersion:   "",
	Language:    nil,
	UserToken:   `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDMxMjA1OTgsImp0aSI6IjUifQ.10EYxB5FZaEiRNpzNLlU1H2MEoPvTAfWYmPtQDJ2hWY`, // 5
	//UserToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDMxMTUzMzMsImp0aSI6IjQifQ.Ki1lKJF5CSA-gzOZYQi4IDvt8CJhxuXNRmofa40gIns`, // 4
	//UserToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDI1NzA5NzAsImp0aSI6IjMifQ.HATz69UJPo6lEL0GF5eSqLkFN-9s1Ej0TJfTUGRX-90`, // 3
	//UserToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDI3NjU4NTQsImp0aSI6IjIifQ.hOXp17_KUXI_HjJRKGFvI6xRiUBKvT_2p4rKKqBTOqM`, // 2
	//UserToken:       `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDI3ODAzMzcsImp0aSI6IjEifQ.oveSP7-bsAe4CUX0JcRmWZFtySG-Egady5tzuRmVJTE`, // 1
	KeepAliveSecond: time.Second * 30,
	LogConf: logx.LogConf{
		ServiceName:         "sdk",
		Mode:                "console",
		Encoding:            "json",
		TimeFormat:          "",
		Path:                "",
		Level:               "info",
		MaxContentLength:    0,
		Compress:            false,
		Stat:                false,
		KeepDays:            0,
		StackCooldownMillis: 0,
		MaxBackups:          0,
		MaxSize:             0,
		Rotation:            "",
	},
}

func getHttpClient(t *testing.T, config *Config) IClient {
	if config == nil {
		config = defaultConfig
	}
	client, err := NewHttpClient(config)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return client
}

func getWsClient(t *testing.T, config *Config) IClient {
	if config == nil {
		config = defaultConfig
	}
	client, err := NewWsClient(config)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return client
}

// GatewayGetUserConnection 获取用户连接
func TestHttpClient_GatewayGetUserConnection(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	gatewayGetUserConnectionResp, err := client.GatewayGetUserConnection(&pb.GatewayGetUserConnectionReq{
		Header: &pb.RequestHeader{},
		UserId: "",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayGetUserConnectionResp))
}

// GatewayBatchGetUserConnection 批量获取用户连接
func TestHttpClient_GatewayBatchGetUserConnection(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	gatewayBatchGetUserConnectionResp, err := client.GatewayBatchGetUserConnection(&pb.GatewayBatchGetUserConnectionReq{
		Header: &pb.RequestHeader{},
		UserIds: []string{
			"",
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayBatchGetUserConnectionResp))
}

// GatewayGetConnectionByFilter 根据条件获取用户连接
func TestHttpClient_GatewayGetConnectionByFilter(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	gatewayGetConnectionByFilterResp, err := client.GatewayGetConnectionByFilter(&pb.GatewayGetConnectionByFilterReq{
		Header: &pb.RequestHeader{},
		Filter: &pb.GatewayGetConnectionFilter{
			UserIds: []string{
				"",
			},
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayGetConnectionByFilterResp))
}

// GatewayKickWs 踢出用户连接
func TestHttpClient_GatewayKickWs(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	gatewayKickWsResp, err := client.GatewayKickWs(&pb.GatewayKickWsReq{
		Header: &pb.RequestHeader{},
		Filter: &pb.GatewayGetConnectionFilter{
			UserIds: []string{
				"",
			},
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayKickWsResp))
}

// GatewayKeepAlive 保持用户连接
func TestHttpClient_GatewayKeepAlive(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(50 * time.Second)
	_ = client
}
