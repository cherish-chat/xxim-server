package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
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
	Account: AccountConfig{
		AuthType: AuthType_Password,
		Password: &AccountConfigPassword{
			Username: "xx",
			Password: "xx",
		},
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

// GatewayWriteDataToWs 向用户连接写入数据
func TestHttpClient_GatewayWriteDataToWs(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	gatewayWriteDataToWsResp, err := client.GatewayWriteDataToWs(&pb.GatewayWriteDataToWsReq{
		Header: &pb.RequestHeader{},
		Filter: &pb.GatewayGetConnectionFilter{
			UserIds: []string{
				"",
			},
		},
		Data: utils.Json.MarshalToBytes(&pb.GatewayApiResponse{
			Header: &pb.ResponseHeader{
				Code:       pb.ResponseCode_SUCCESS,
				ActionType: 0,
				ActionData: "",
				Extra:      "sdfghjklkjhgfds",
			},
			RequestId: "1234567890987654321",
			Path:      "aaa/aaaa/aaaaa",
			Body:      []byte("1234567890"),
		}),
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayWriteDataToWsResp))
	time.Sleep(1 * time.Second)
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
