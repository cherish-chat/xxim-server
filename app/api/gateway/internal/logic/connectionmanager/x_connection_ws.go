package connectionmanager

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/cherish-chat/xxim-proto/peerpb"
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

func (w *websocketWrapper) CloseConnection() {
	_ = w.conn.Close(websocket.StatusGoingAway, "")
}

// NewWebsocketConnect 创建websocket连接
func NewWebsocketConnect(
	ctx context.Context,
	header *peerpb.RequestHeader,
	conn *websocket.Conn,
) *Connection {
	ecdh := utils.NewECDH(elliptic.P256())
	privateKey, publicKey, _ := ecdh.GenerateKey(rand.Reader)

	connection := &Connection{
		ctx:              ctx,
		Header:           header,
		HeaderLock:       sync.RWMutex{},
		ServerPrivateKey: privateKey,
		ServerPublicKey:  publicKey,
		ClientPublicKey:  nil,
		PublicKeyLock:    sync.RWMutex{},
		Connection:       &websocketWrapper{conn: conn},
		ConnectedTime:    time.Now(),
	}

	return connection
}
