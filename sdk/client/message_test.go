package client

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/mr"
	"testing"
	"time"
)

// MessageBatchSend 批量发送消息
func TestHttpClient_MessageBatchSend(t *testing.T) {
	//client := getHttpClient(t, nil)
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	var messages []*pb.Message
	for i := 0; i < 1; i++ {
		messages = append(messages, generateMessageToFriend(&pb.Message_Sender{
			Id:     "1",
			Name:   "哈哈",
			Avatar: "头像",
			Extra:  "xx",
		}, "2", "你好你好", map[string]string{
			"platformSource": "Test",
		}))
	}
	messageBatchSendResp, err := client.MessageBatchSend(&pb.MessageBatchSendReq{
		Messages:     messages,
		DisableQueue: false,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(messageBatchSendResp))
}

// MessageBatchSend 压力测试
func TestHttpClient_MessageBatchSend_Pressure(t *testing.T) {
	var fs []func()
	goroutineCount := 100
	loopCount := 1000
	client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	message := generateMessageToFriend(&pb.Message_Sender{
		Id:     "1",
		Name:   "哈哈",
		Avatar: "头像",
		Extra:  "xx",
	}, "2", "你好你好", map[string]string{
		"platformSource": "Test",
	})
	for i := 0; i < goroutineCount; i++ {
		fs = append(fs, func() {
			for j := 0; j < loopCount; j++ {
				_, _ = client.MessageBatchSend(&pb.MessageBatchSendReq{
					Messages:     []*pb.Message{message},
					DisableQueue: false,
				})
			}
		})
	}
	t.Logf("开始时间: %s", time.Now().Format("2006-01-02 15:04:05.000"))
	mr.FinishVoid(fs...)
	t.Logf("结束时间: %s", time.Now().Format("2006-01-02 15:04:05.000"))
}

func generateMessageToFriend(
	sender *pb.Message_Sender,
	friendId string,
	text string,
	extraMap map[string]string) *pb.Message {
	convId := pb.GetSingleChatConversationId(sender.Id, friendId)
	content := &pb.MessageContentText{Items: []*pb.MessageContentText_Item{{
		Type:  pb.MessageContentText_Item_TEXT,
		Text:  text,
		Image: nil,
		At:    nil,
	}}}
	contentBytes, _ := json.Marshal(content)
	return &pb.Message{
		Uuid:             utils.Snowflake.String(),
		ConversationId:   convId,
		ConversationType: pb.ConversationType_Single,
		Sender:           sender,
		Content:          contentBytes,
		ContentType:      pb.MessageContentType_Text,
		SendTime:         time.Now().UnixMilli(),
		Option: &pb.Message_Option{
			StorageForServer: true,
			StorageForClient: true,
			NeedDecrypt:      false,
			CountUnread:      true,
		},
		ExtraMap: extraMap,
	}
}
