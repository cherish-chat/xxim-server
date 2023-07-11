package types

import (
	"context"
)

type LongConnection interface {
	//SendMessage 发送消息
	SendMessage(ctx context.Context, message []byte) error
	//CloseConnection 关闭连接
	CloseConnection()
}
