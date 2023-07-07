package types

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type LongConnection interface {
	//发送消息
	SendMessage(ctx context.Context, message []byte) error
	//关闭连接
	CloseConnection(code pb.WebsocketCustomCloseCode, reason string)
}
