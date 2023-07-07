package gatewayservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/i18n"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthenticationConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationConnectionLogic {
	return &AuthenticationConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AuthenticationConnection 验证连接
func (l *AuthenticationConnectionLogic) AuthenticationConnection(in *pb.AuthenticationConnectionReq) (*pb.AuthenticationConnectionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AuthenticationConnectionResp{}, nil
}

func (l *AuthenticationConnectionLogic) AuthenticationConnection_(connection *Connection, in *pb.AuthenticationConnectionReq) (*pb.AuthenticationConnectionResp, error) {
	userBeforeConnectResp, err := l.svcCtx.CallbackService.UserBeforeConnect(l.ctx, &pb.UserBeforeConnectReq{
		Header: in.Header,
		UserId: in.UserId,
		Token:  in.Token,
	})
	if err != nil {
		l.Errorf("UserBeforeConnect err: %v", err)
		return &pb.AuthenticationConnectionResp{
			Header:  i18n.NewAuthError(pb.AuthErrorTypeInvalid, i18n.ServerError),
			Success: false,
		}, nil
	}
	if !userBeforeConnectResp.Success {
		return &pb.AuthenticationConnectionResp{
			Header: userBeforeConnectResp.Header,
		}, nil
	}
	// 验证通过
	connection.headerLock.Lock()
	connection.header = &pb.RequestHeader{
		AppId:       connection.header.AppId,
		UserId:      userBeforeConnectResp.UserId,
		ClientIp:    connection.header.ClientIp,
		InstallId:   connection.header.InstallId,
		Platform:    connection.header.Platform,
		DeviceModel: connection.header.DeviceModel,
		OsVersion:   connection.header.OsVersion,
		AppVersion:  connection.header.AppVersion,
		Extra:       connection.header.Extra,
	}
	connection.headerLock.Unlock()
	return &pb.AuthenticationConnectionResp{
		Header:  i18n.NewOkHeader(),
		Success: true,
	}, nil
}
