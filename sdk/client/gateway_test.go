package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"testing"
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

func TestHttpClient_GatewayGetUserConnection(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	gatewayGetUserConnectionResp, err := client.GatewayGetUserConnection(&pb.GatewayGetUserConnectionReq{
		Header: &pb.RequestHeader{},
		UserId: "",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayGetUserConnectionResp))
}
