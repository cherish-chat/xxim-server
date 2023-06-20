package client

import (
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"testing"
	"time"
)

// ListNotice 列出通知
func TestHttpClient_ListNotice(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	friendApplyResp, err := client.ListNotice(&pb.ListNoticeReq{
		ConversationId:   subscriptionmodel.ConversationIdFriendNotification,
		ConversationType: pb.ConversationType_Subscription,
		UpdateTimeGt:     1687105599171,
		Limit:            100,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(friendApplyResp))
}
