package gatewayservicelogic

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type p2pWrapper struct {
	dataChannel *webrtc.DataChannel
}

func (w *p2pWrapper) SendMessage(ctx context.Context, message []byte) error {
	return w.dataChannel.Send(message)
}

func (w *p2pWrapper) CloseConnection(code pb.WebsocketCustomCloseCode, reason string) {
	logx.Infof("p2pWrapper CloseConnection code: %d, reason: %s", code, reason)
	_ = w.dataChannel.Close()
}

func NewP2pConnection(
	ctx context.Context,
	header *pb.RequestHeader,
	dataChannel *webrtc.DataChannel,
) *Connection {
	ecdh := utils.NewECDH(elliptic.P256())
	privateKey, publicKey, _ := ecdh.GenerateKey(rand.Reader)

	return &Connection{
		ctx:              ctx,
		header:           header,
		headerLock:       sync.RWMutex{},
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
