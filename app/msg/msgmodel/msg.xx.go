package msgmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"time"
)

func CreateTextMsgToUser(
	sender *pb.UserBaseInfo,
	userId string,
	text string,
	options MsgOptions,
	offlinePush *MsgOfflinePush,
	ext any,
) *Msg {
	return CreateCustomMsgToUser(sender, userId, pb.ContentType_TEXT, xorm.M{
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
		SenderId:       sender.Id,
		SenderInfo:     make([]byte, 0),
		SenderConvInfo: make([]byte, 0),
		ConvId:         pb.SingleConvId(sender.Id, userId),
		ContentType:    contentType,
		Content:        utils.AnyToBytes(content),
		Options:        options,
		OfflinePush:    offlinePush,
		Ext:            utils.AnyToBytes(ext),
	}
}

func CreateConvProfileChangeMsg(
	userId string,
	convId string,
	content []byte,
) *Msg {
	return &Msg{
		ClientMsgId:    utils.GenId(),
		ClientTime:     time.Now().UnixMilli(),
		SenderId:       userId,
		SenderInfo:     make([]byte, 0),
		SenderConvInfo: make([]byte, 0),
		ConvId:         convId,
		ContentType:    pb.ContentType_CONV_PROFILE_CHANGE,
		Content:        content,
		Options: MsgOptions{
			OfflinePush:       false,
			StorageForServer:  true,
			StorageForClient:  true,
			UpdateUnreadCount: false,
			NeedDecrypt:       false,
			UpdateConvMsg:     false,
		},
		OfflinePush: nil,
		Ext:         make([]byte, 0),
	}
}
