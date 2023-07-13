package callbackservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeConnectLogic {
	return &UserBeforeConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserBeforeConnectLogic) UserBeforeConnect(in *peerpb.UserBeforeConnectReq) (*peerpb.UserBeforeConnectResp, error) {
	if in.Token != "" {
		tokenObject, verifyTokenErr := l.svcCtx.Jwt.VerifyToken(context.Background(), in.Token, usermodel.GetJwtUniqueKey(in.Header))
		if verifyTokenErr != nil {
			l.Errorf("verifyTokenErr: %v", verifyTokenErr)
			var resp *peerpb.UserBeforeConnectResp
			switch verifyTokenErr {
			case utils.TokenExpiredError:
				resp = &peerpb.UserBeforeConnectResp{
					Header: peerpb.NewAuthError(peerpb.AuthErrorTypeExpired, ""),
				}
			case utils.TokenReplaceError:
				ssm := utils.Map.SSMFromString(tokenObject.Extra)
				deviceModel := ssm.Get("deviceModel")
				resp = &peerpb.UserBeforeConnectResp{
					Header: peerpb.NewAuthError(peerpb.AuthErrorTypeReplace, deviceModel),
				}
			default:
				resp = &peerpb.UserBeforeConnectResp{
					Header: peerpb.NewAuthError(peerpb.AuthErrorTypeInvalid, verifyTokenErr.Error()),
				}
			}
			return resp, nil
		}
		l.Debugf("tokenObject: %+v", tokenObject)
		// 验证权限
		return &peerpb.UserBeforeConnectResp{
			Header:  peerpb.NewOkHeader(),
			UserId:  tokenObject.UserId,
			Success: true,
		}, nil
	} else {
		// 验证权限
		l.Infof("in.Header: %+v", in.Header)
		return &peerpb.UserBeforeConnectResp{
			Header:  peerpb.NewAuthError(peerpb.AuthErrorTypeInvalid, ""),
			Success: false,
		}, nil
	}
}
