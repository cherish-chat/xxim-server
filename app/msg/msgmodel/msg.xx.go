package msgmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"time"
)

type H map[string]any

func CreateTextMsgToUser(
	sender *pb.UserBaseInfo,
	userId string,
	text string,
	options MsgOptions,
	offlinePush *MsgOfflinePush,
	ext any,
) *Msg {
	return CreateCustomMsgToUser(sender, userId, pb.ContentType_TEXT, H{
		"text": text,
	}, options, offlinePush, ext)
}

func CreateCustomMsgToUser(
	sender *pb.UserBaseInfo,
	userId string,
	contentType ContentType,
	content any,
	options MsgOptions,
	offlinePush *MsgOfflinePush,
	ext any,
) *Msg {
	return &Msg{
		ClientMsgId:    utils.GenId(),
		ClientTime:     time.Now().UnixMilli(),
		Sender:         sender.Id,
		SenderInfo:     "",
		SenderConvInfo: "",
		Receiver:       MsgReceiver{UserId: userId},
		ConvId:         SingleConvId(sender.Id, userId),
		ContentType:    contentType,
		Content:        utils.AnyToBytes(content),
		Options:        options,
		OfflinePush:    offlinePush,
		Ext:            utils.AnyToBytes(ext),
	}
}
