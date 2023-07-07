package gatewayservicelogic

import (
	"context"
	"crypto/elliptic"
	"errors"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyConnectionLogic {
	return &VerifyConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VerifyConnection 验证连接
func (l *VerifyConnectionLogic) VerifyConnection(in *pb.VerifyConnectionReq) (*pb.VerifyConnectionResp, error) {

	return &pb.VerifyConnectionResp{}, nil
}

func (l *VerifyConnectionLogic) VerifyConnection_(connection *Connection, in *pb.VerifyConnectionReq) (*pb.VerifyConnectionResp, error) {
	ecdh := utils.NewECDH(elliptic.P256())
	publicKey, ok := ecdh.Unmarshal(in.PublicKey)
	if !ok {
		return &pb.VerifyConnectionResp{}, errors.New(i18n.PublicKeyError)
	}
	// 计算共享密钥
	secret, _ := ecdh.GenerateSharedSecret(connection.ServerPrivateKey, publicKey)
	connection.PublicKeyLock.Lock()
	connection.ClientPublicKey = publicKey
	connection.SharedSecret = secret
	connection.PublicKeyLock.Unlock()
	return &pb.VerifyConnectionResp{
		PublicKey: ecdh.Marshal(connection.ServerPublicKey),
	}, nil
}
