package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"testing"
	"time"
)

// GroupCreate 创建群
func TestHttpClient_GroupCreate(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	groupCreateResp, err := client.GroupCreate(&pb.GroupCreateReq{
		Name:       nil,
		Avatar:     nil,
		MemberList: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		//InfoMap: map[string]string{
		//
		//},
		ExtraMap: nil,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(groupCreateResp))
}
