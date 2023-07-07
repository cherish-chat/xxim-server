package gatewayservicelogic

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type websocketWrapper struct {
	conn *websocket.Conn
}

func (w *websocketWrapper) SendMessage(ctx context.Context, message []byte) error {
	return w.conn.Write(ctx, websocket.MessageBinary, message)
}

func (w *websocketWrapper) CloseConnection(code pb.WebsocketCustomCloseCode, reason string) {
	_ = w.conn.Close(websocket.StatusCode(code), reason)
}

// NewWebsocketConnect 创建websocket连接
func NewWebsocketConnect(
	ctx context.Context,
	header *pb.RequestHeader,
	conn *websocket.Conn,
) *Connection {
	ecdh := utils.NewECDH(elliptic.P256())
	privateKey, publicKey, _ := ecdh.GenerateKey(rand.Reader)

	connection := &Connection{
		ctx:              ctx,
		header:           header,
		headerLock:       sync.RWMutex{},
		ServerPrivateKey: privateKey,
		ServerPublicKey:  publicKey,
		ClientPublicKey:  nil,
		PublicKeyLock:    sync.RWMutex{},
		Connection:       &websocketWrapper{conn: conn},
		ConnectedTime:    time.Now(),
	}

	return connection
}
