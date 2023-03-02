package types

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"time"
)

type (
	IConn interface {
		Close(code int, desc string) error
		Write(ctx context.Context, typ int, msg []byte) error
		Read(ctx context.Context) (int, []byte, error)
	}
	UserConn struct {
		Pointer     string
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
		DeviceModel string            // 设备型号
		OsVersion   string            // 系统版本
		AppVersion  string            // app版本
		Language    string            // 语言
		Platform    string            // ios, android, web, pc, mac, linux, windows
		Ips         string            // ip
		NetworkUsed string            // 4G/5G/WIFI
		Headers     map[string]string // 其他参数
		Timestamp   int64             // 时间戳
		AesKey      *string           // aes key
		AesIv       *string           // aes iv
	}
)

func (c *UserConn) SetConnParams(connParam *pb.ConnParam) {
	c.ConnParam.UserId = connParam.UserId
	c.ConnParam.Token = connParam.Token
	c.ConnParam.DeviceId = connParam.DeviceId
	c.ConnParam.DeviceModel = connParam.DeviceModel
	c.ConnParam.OsVersion = connParam.OsVersion
	c.ConnParam.AppVersion = connParam.AppVersion
	c.ConnParam.Language = connParam.Language
	c.ConnParam.Platform = connParam.Platform
	c.ConnParam.Ips = connParam.Ips
	c.ConnParam.NetworkUsed = connParam.NetworkUsed
	c.ConnParam.Headers = connParam.Headers
	c.ConnParam.AesKey = connParam.AesKey
	c.ConnParam.AesIv = connParam.AesIv
}

func WebsocketStatusCodeAuthFailed(code int) int {
	return 3000
}

func WebsocketStatusCodeRsaFailed() int {
	return 3001
}

func WebsocketStatusCodePlatformFailed() int {
	return 3002
}
