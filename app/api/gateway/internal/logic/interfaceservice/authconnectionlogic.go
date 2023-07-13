package interfaceservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	internalservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthConnectionLogic {
	return &AuthConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AuthConnection 验证连接用户id和token
func (l *AuthConnectionLogic) AuthConnection(in *peerpb.AuthConnectionReq) (*peerpb.AuthConnectionResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.AuthConnectionResp{}, nil
}

func (l *AuthConnectionLogic) AuthConnection_(connection *internalservicelogic.Connection, in *peerpb.AuthConnectionReq) (*peerpb.AuthConnectionResp, error) {
	if in.UserId == "" || in.Token == "" {
		//取消验证
		oldUserId := connection.GetHeader().UserId
		connection.HeaderLock.Lock()
		connection.Header = &peerpb.RequestHeader{
			AppId:       connection.Header.AppId,
			UserId:      "",
			ClientIp:    connection.Header.ClientIp,
			InstallId:   connection.Header.InstallId,
			Platform:    connection.Header.Platform,
			DeviceModel: connection.Header.DeviceModel,
			OsVersion:   connection.Header.OsVersion,
			AppVersion:  connection.Header.AppVersion,
			Extra:       connection.Header.Extra,
		}
		connection.HeaderLock.Unlock()
		if oldUserId != "" {
			_, _ = l.svcCtx.CallbackService.UserAfterOffline(context.Background(), &peerpb.UserAfterOfflineReq{
				UserId: oldUserId,
			})
		}
		return &peerpb.AuthConnectionResp{
			Header:  peerpb.NewOkHeader(),
			Success: true,
		}, nil
	}
	userBeforeConnectResp, err := l.svcCtx.CallbackService.UserBeforeConnect(context.Background(), &peerpb.UserBeforeConnectReq{
		Header: connection.GetHeader(),
		UserId: in.UserId,
		Token:  in.Token,
	})
	if err != nil {
		l.Errorf("UserBeforeConnect err: %v", err)
		return &peerpb.AuthConnectionResp{
			Header:  peerpb.NewAuthError(peerpb.AuthErrorTypeInvalid, peerpb.ServerError),
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	if !userBeforeConnectResp.Success {
		return &peerpb.AuthConnectionResp{
			Header: userBeforeConnectResp.Header,
			Error:  "auth failed",
		}, nil
	}
	// 验证通过
	connection.HeaderLock.Lock()
	connection.Header = &peerpb.RequestHeader{
		AppId:       connection.Header.AppId,
		UserId:      userBeforeConnectResp.UserId,
		ClientIp:    connection.Header.ClientIp,
		InstallId:   connection.Header.InstallId,
		Platform:    connection.Header.Platform,
		DeviceModel: connection.Header.DeviceModel,
		OsVersion:   connection.Header.OsVersion,
		AppVersion:  connection.Header.AppVersion,
		Extra:       connection.Header.Extra,
	}
	connection.HeaderLock.Unlock()
	return &peerpb.AuthConnectionResp{
		Header:  peerpb.NewOkHeader(),
		Success: true,
	}, nil
}
