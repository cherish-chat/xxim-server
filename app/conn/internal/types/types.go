package types

import (
	"context"
	"time"
)

type (
	IConn interface {
		Close(code int, desc string) error
		Write(ctx context.Context, typ int, msg []byte) error
		Read(ctx context.Context) (int, []byte, error)
	}
	UserConn struct {
		Conn        IConn
		ConnParam   ConnParam
		Ctx         context.Context
		ConnectedAt time.Time
	}
	IServer interface {
		SetBeforeConnect(f func(ctx context.Context, param ConnParam) (int, error))
		SetAddSubscriber(f func(c *UserConn))
		SetDeleteSubscriber(f func(c *UserConn))
		SetOnReceive(f func(ctx context.Context, c *UserConn, typ int, msg []byte))
		Start() error
	}
	MsgBytes struct {
		UserId string
		Msg    []byte
	}
	ConnParam struct {
		UserId      string            // uid
		Token       string            // token
		DeviceId    string            // 设备id
		Platform    string            // ios, android, web, pc, mac, linux, windows
		Ips         string            // ip
		NetworkUsed string            // 4G/5G/WIFI
		Headers     map[string]string // 其他参数
	}
)
