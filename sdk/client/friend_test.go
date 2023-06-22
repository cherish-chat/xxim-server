package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"testing"
	"time"
)

// FriendApply 申请好友
func TestHttpClient_FriendApply(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	friendApplyResp, err := client.FriendApply(&pb.FriendApplyReq{
		ToUserId: "2",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(friendApplyResp))
}

// ListFriendApply 列出好友申请
func TestHttpClient_ListFriendApply(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	listFriendApplyResp, err := client.ListFriendApply(&pb.ListFriendApplyReq{
		Cursor: 0,
		Limit:  10,
		Filter: &pb.ListFriendApplyReq_Filter{
			Status: utils.AnyPtr(pb.FriendApplyStatus_Applying),
		},
		Option: &pb.ListFriendApplyReq_Option{
			IncludeApplyByMe: true,
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(listFriendApplyResp))
}

// FriendApplyHandle 处理好友申请
func TestHttpClient_FriendApplyHandle(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	friendApplyHandleResp, err := client.FriendApplyHandle(&pb.FriendApplyHandleReq{
		ApplyId:      "4d35f233b59136d7280656f677e43fe5",
		Agree:        true,
		FirstMessage: nil,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(friendApplyHandleResp))
}
