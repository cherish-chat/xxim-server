package callbackservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"regexp"
	"strings"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeRequestLogic {
	return &UserBeforeRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserBeforeRequest 用户请求前的回调
func (l *UserBeforeRequestLogic) UserBeforeRequest(in *pb.UserBeforeRequestReq) (*pb.UserBeforeRequestResp, error) {
	if in.Header.UserToken != "" {
		tokenObject, verifyTokenErr := l.svcCtx.Jwt.VerifyToken(l.ctx, in.Header.UserToken, in.Header.GetJwtUniqueKey())
		if verifyTokenErr != nil {
			l.Errorf("verifyTokenErr: %v", verifyTokenErr)
			var resp *pb.UserBeforeRequestResp
			switch verifyTokenErr {
			case utils.TokenExpiredError:
				resp = &pb.UserBeforeRequestResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeExpired, ""),
				}
			case utils.TokenReplaceError:
				ssm := utils.Map.SSMFromString(tokenObject.Extra)
				deviceModel := ssm.Get("deviceModel")
				resp = &pb.UserBeforeRequestResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeReplace, deviceModel),
				}
			default:
				resp = &pb.UserBeforeRequestResp{
					Header: i18n.NewAuthError(pb.AuthErrorTypeInvalid, ""),
				}
			}
			// 如果是白名单接口，那么就不需要返回错误
			if strings.Contains(in.Path, "/white/") {
				resp.Header = i18n.NewOkHeader()
				l.Debugf("white path: %v", in.Path)
			} else {
				l.Debugf("not white path: %v", in.Path)
			}
			return resp, nil
		}
		l.Debugf("tokenObject: %+v", tokenObject)
		// 验证权限
		if !strings.Contains(in.Path, "/white/") {
			verifyScope := false
			for _, scopeRegex := range tokenObject.Scope {
				// 是否匹配 path 只要有一个匹配就可以
				matched, err := regexp.MatchString(scopeRegex, in.Path)
				if err != nil {
					l.Errorf("regexp.MatchString error: %v", err)
					return &pb.UserBeforeRequestResp{
						Header: i18n.NewAuthError(pb.AuthErrorTypeInvalid, ""),
					}, nil
				}
				if matched {
					verifyScope = true
					break
				}
			}
			if !verifyScope {
				l.Errorf("verifyScope error: %v", verifyScope)
				return &pb.UserBeforeRequestResp{
					Header: i18n.NewForbiddenError(),
				}, nil
			}
		}
		return &pb.UserBeforeRequestResp{
			Header: i18n.NewOkHeader(),
			UserId: tokenObject.UserId,
		}, nil
	}
	// 如果是白名单接口，那么就不需要返回错误
	if strings.Contains(in.Path, "/white/") {
		return &pb.UserBeforeRequestResp{
			Header: i18n.NewOkHeader(),
		}, nil
	}
	return &pb.UserBeforeRequestResp{
		Header: i18n.NewAuthError(pb.AuthErrorTypeInvalid, ""),
	}, nil
}
