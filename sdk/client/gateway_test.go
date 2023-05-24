package client

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/types"
	"testing"
)

func TestHttpClient_GatewayGetUserConnection(t *testing.T) {
	client, err := NewHttpClient(&Config{
		Endpoints:       []string{"http://127.0.0.1:34500"},
		EnableEncrypted: true,
		AesKey:          "xx",
		AesIv:           "xx",
		//ContentType:     "protobuf",
		ContentType: "json",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	gatewayGetUserConnectionResp, err := client.GatewayGetUserConnection(&types.GatewayGetUserConnectionReq{
		Header:     &types.RequestHeader{},
		UserIds:    nil,
		Platforms:  nil,
		InstallIds: nil,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(gatewayGetUserConnectionResp))
}
