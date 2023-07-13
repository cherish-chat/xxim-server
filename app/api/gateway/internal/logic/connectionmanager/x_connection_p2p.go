package connectionmanager

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/pion/webrtc/v2"
	"sync"
	"time"
)

type p2pWrapper struct {
	dataChannel *webrtc.DataChannel
}

func (w *p2pWrapper) SendMessage(ctx context.Context, message []byte) error {
	return w.dataChannel.Send(message)
}

func (w *p2pWrapper) CloseConnection() {
	_ = w.dataChannel.Close()
}

func NewP2pConnection(
	ctx context.Context,
	header *peerpb.RequestHeader,
	dataChannel *webrtc.DataChannel,
) *Connection {
	ecdh := utils.NewECDH(elliptic.P256())
	privateKey, publicKey, _ := ecdh.GenerateKey(rand.Reader)

	return &Connection{
		ctx:              ctx,
		Header:           header,
		HeaderLock:       sync.RWMutex{},
		ServerPrivateKey: privateKey,
		ServerPublicKey:  publicKey,
		ClientPublicKey:  nil,
		PublicKeyLock:    sync.RWMutex{},
		Connection: &p2pWrapper{
			dataChannel: dataChannel,
		},
		ConnectedTime: time.Now(),
	}
}
