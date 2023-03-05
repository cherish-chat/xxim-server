package types

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"nhooyr.io/websocket"
)

type EventHandler interface {
	// BeforeClose is called before the connection is closed.
	BeforeClose(code websocket.StatusCode, reason string)
	// AfterClose is called after the connection is closed.
	AfterClose(code websocket.StatusCode, reason string)
	// BeforeReConnect is called before the connection is reconnected.
	BeforeReConnect()
	// AfterReConnect is called after the connection is reconnected.
	AfterReConnect()
	// BeforeConnect is called before the connection is connected.
	BeforeConnect()
	// AfterConnect is called after the connection is connected.
	AfterConnect()
	// OnMessage is called when a message is received.
	OnMessage(typ websocket.MessageType, message []byte)
	// OnPushMsgDataList is called when a message is received.
	OnPushMsgDataList(body *pb.PushBody)
	// OnPushNoticeData is called when a message is received.
	OnPushNoticeData(noticeData *pb.NoticeData) bool
	// OnPushResponseBody is called when a message is received.
	OnPushResponseBody(body *pb.PushBody)
	// OnTimer is called when a timer is triggered.
	OnTimer()
}
