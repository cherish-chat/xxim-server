package callbackservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"

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

// UserBeforeConnect 用户连接前的回调
func (l *UserBeforeConnectLogic) UserBeforeConnect(in *pb.UserBeforeConnectReq) (*pb.UserBeforeConnectResp, error) {
	if in.Token != "" {
		tokenObject, verifyTokenErr := l.svcCtx.Jwt.VerifyToken(l.ctx, in.Token, in.Header.GetJwtUniqueKey())
		if verifyTokenErr != nil {
			l.Errorf("verifyTokenErr: %v", verifyTokenErr)
			var resp *pb.UserBeforeConnectResp
			switch verifyTokenErr {
			case utils.TokenExpiredError:
				resp = &pb.UserBeforeConnectResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeExpired, ""),
				}
			case utils.TokenReplaceError:
				ssm := utils.Map.SSMFromString(tokenObject.Extra)
				deviceModel := ssm.Get("deviceModel")
				resp = &pb.UserBeforeConnectResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeReplace, deviceModel),
				}
			default:
				resp = &pb.UserBeforeConnectResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeInvalid, ""),
				}
			}
			return resp, nil
		}
		l.Debugf("tokenObject: %+v", tokenObject)
		// 验证权限
		return &pb.UserBeforeConnectResp{
			Header:  i18n.NewOkHeader(),
			UserId:  tokenObject.UserId,
			Success: true,
		}, nil
	} else {
		// 验证权限
		l.Infof("in.Header: %+v", in.Header)
		return &pb.UserBeforeConnectResp{
			Header:  i18n.NewAuthError(pb.AuthErrorTypeInvalid, ""),
			Success: false,
		}, nil
	}
}
