package handler

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"time"
)

type ConnLogic struct {
	svcCtx *svc.ServiceContext
}

func (l *ConnLogic) BeforeClose(code websocket.StatusCode, reason string) {

}

func (l *ConnLogic) AfterClose(code websocket.StatusCode, reason string) {
	// 重连
	time.Sleep(time.Second * 3)
	l.svcCtx.Client().ReConnect()
}

func (l *ConnLogic) BeforeReConnect() {

}

func (l *ConnLogic) AfterReConnect() {

}

func (l *ConnLogic) BeforeConnect() {

}

func (l *ConnLogic) AfterConnect() {
	logx.Infof("connect success")
}

func (l *ConnLogic) OnMessage(typ websocket.MessageType, message []byte) {
	//logx.Debugf("message type: %d, message.len: %d", typ, len(message))
}

func (l *ConnLogic) OnPushMsgDataList(body *pb.PushBody) {
	msgDataList := &pb.MsgDataList{}
	err := proto.Unmarshal(body.Data, msgDataList)
	if err != nil {
		logx.Errorf("unmarshal MsgDataList error: %s", err.Error())
		return
	}
	for _, msgData := range msgDataList.MsgDataList {
		logx.Debugf("msgData: %s", utils.AnyToString(msgData))
	}
}

func (l *ConnLogic) OnPushNoticeData(noticeData *pb.NoticeData) bool {
	logx.Debugf("noticeData: %s", utils.AnyToString(noticeData))
	return true
}

func (l *ConnLogic) OnPushResponseBody(body *pb.PushBody) {

}

func (l *ConnLogic) OnTimer() {
}

func NewEventHandler(svcCtx *svc.ServiceContext) *ConnLogic {
	return &ConnLogic{svcCtx: svcCtx}
}
