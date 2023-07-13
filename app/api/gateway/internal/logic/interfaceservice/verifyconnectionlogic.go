package interfaceservicelogic

import (
	"context"
	"crypto/elliptic"
	"errors"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

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
func (l *VerifyConnectionLogic) VerifyConnection(in *peerpb.VerifyConnectionReq) (*peerpb.VerifyConnectionResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.VerifyConnectionResp{}, nil
}

func (l *VerifyConnectionLogic) VerifyConnection_(connection *connectionmanager.Connection, in *peerpb.VerifyConnectionReq) (*peerpb.VerifyConnectionResp, error) {
	var err error
	in.PublicKey, err = l.svcCtx.RsaInstance.Decrypt(in.PublicKey)
	if err != nil {
		return &peerpb.VerifyConnectionResp{}, errors.New(peerpb.RsaDecryptError)
	}
	ecdh := utils.NewECDH(elliptic.P256())
	publicKey, ok := ecdh.UnmarshalHex(string(in.PublicKey))
	if !ok {
		return &peerpb.VerifyConnectionResp{}, errors.New(peerpb.PublicKeyError)
	}
	connection.PublicKeyLock.Lock()
	connection.ClientPublicKey = publicKey
	connection.PublicKeyLock.Unlock()
	oldHeader := connection.GetHeader()
	connection.HeaderLock.Lock()
	connection.Header = &peerpb.RequestHeader{
		AppId:       in.Header.AppId,
		UserId:      "",
		ClientIp:    oldHeader.ClientIp,
		InstallId:   in.Header.InstallId,
		Platform:    in.Header.Platform,
		DeviceModel: in.Header.DeviceModel,
		OsVersion:   in.Header.OsVersion,
		AppVersion:  in.Header.AppVersion,
		Extra:       in.Header.Extra,
	}
	connection.HeaderLock.Unlock()
	return &peerpb.VerifyConnectionResp{
		PublicKey: []byte(ecdh.MarshalHex(connection.ServerPublicKey)),
	}, nil
}
