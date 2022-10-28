package types

import (
	"context"
	"time"
)

type IConn interface {
	Close(code int, desc string) error
	Write(ctx context.Context, typ int, msg []byte) error
}
type UserConn struct {
	Conn        IConn
	ConnParam   ConnParam
	Ctx         context.Context
	ConnectedAt time.Time
}
type IServer interface {
	SetBeforeConnect(f func(ctx context.Context, param ConnParam) (int, error))
	SetAddSubscriber(f func(c *UserConn))
	SetDeleteSubscriber(f func(c *UserConn))
	Start() error
}
type MsgBytes struct {
	UserId string
	Msg    []byte
}
type ConnParam struct {
	UserId      string         // uid
	Token       string         // token
	DeviceId    string         // 设备id
	Platform    string         // ios, android, web, pc, mac, linux, windows
	Ips         string         // ip
	NetworkUsed string         // 4G/5G/WIFI
	Headers     map[string]any // 其他参数
}
