package tests

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
)

var ctx = context.Background()

func imService() imservice.ImService {
	return imservice.NewImService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{
			"127.0.0.1:24000",
		},
	}))
}

// TestSendMsgLogic1 定时消息
func TestSendMsgLogic1(t *testing.T) {
	selfId := "TestSendMsgLogic1"
	imService := imService()
	resp, err := imService.SendMsg(ctx, &imservice.SendMsgReq{
		SelfId: selfId,
		MsgDataList: []*pb.MsgData{{
			ClientMsgId: utils.GenId(),
			ServerMsgId: "",
			SenderId:    selfId,
			SenderInfo: utils.Any2Bytes(utils.H{
				"nickname": "TestSendMsgLogic1",
				"avatar":   "https://www.baidu.com",
			}),
			ContentType: pb.MsgData_Text,
			Content:     utils.Any2Bytes("TestSendMsgLogic1"),
			ClientTime:  utils.GetNowMilli(),
			ServerTime:  0,
			Seq:         0,
			OfflinePush: &pb.MsgData_OfflinePush{
				Enable:        true,
				Title:         selfId,
				Desc:          "TestSendMsgLogic1",
				Ex:            "",
				IOSPushSound:  "",
				IOSBadgeCount: false,
				AtUserIds:     nil,
			},
			MsgOptions: &pb.MsgData_MsgOptions{
				Storage:    true,
				UpdateConv: true,
				UnreadOpt:  1,
				Rewrite:    false,
			},
			Ex: nil,
			ConvList: []*pb.MsgData_Conv{{
				Id: "TestSendMsgLogic1",
				Info: utils.Any2Bytes(utils.H{
					"type": 1, // 1:单聊 2:群聊
				}),
			}},
			ExcludeUIds: nil,
		}},
		SendAt:      utils.Int64Ptr(utils.GetNowMilli() + 63*1000),
		Platform:    "Test",
		AppVersion:  "v0.0.1",
		DeviceModel: "Mac",
		Ips:         "111.111.111.111",
	})
	if err != nil {
		t.Fatalf("send msg failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestSendMsgLogic2 定时消息
func TestSendMsgLogic2(t *testing.T) {
	selfId := "TestSendMsgLogic2"
	imService := imService()
	resp, err := imService.SendMsg(ctx, &imservice.SendMsgReq{
		SelfId: selfId,
		MsgDataList: []*pb.MsgData{{
			ClientMsgId: utils.GenId(),
			ServerMsgId: "",
			SenderId:    selfId,
			SenderInfo: utils.Any2Bytes(utils.H{
				"nickname": "TestSendMsgLogic2",
				"avatar":   "https://www.baidu.com",
			}),
			ContentType: pb.MsgData_Text,
			Content:     utils.Any2Bytes("TestSendMsgLogic2"),
			ClientTime:  utils.GetNowMilli(),
			ServerTime:  0,
			Seq:         0,
			OfflinePush: &pb.MsgData_OfflinePush{
				Enable:        true,
				Title:         selfId,
				Desc:          "TestSendMsgLogic2",
				Ex:            "",
				IOSPushSound:  "",
				IOSBadgeCount: false,
				AtUserIds:     nil,
			},
			MsgOptions: &pb.MsgData_MsgOptions{
				Storage:    true,
				UpdateConv: true,
				UnreadOpt:  1,
				Rewrite:    false,
			},
			Ex: nil,
			ConvList: []*pb.MsgData_Conv{{
				Id: "TestSendMsgLogic2",
				Info: utils.Any2Bytes(utils.H{
					"type": 1, // 1:单聊 2:群聊
				}),
			}},
			ExcludeUIds: nil,
		}},
		SendAt:      nil,
		Platform:    "Test",
		AppVersion:  "v0.0.1",
		DeviceModel: "Mac",
		Ips:         "111.111.111.111",
	})
	if err != nil {
		t.Fatalf("send msg failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}
