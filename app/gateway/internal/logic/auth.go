package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

type AuthLogic struct {
	svcCtx  *svc.ServiceContext
	ctx     context.Context
	request *http.Request
	logx.Logger
}

func NewAuthLogic(r *http.Request, svcCtx *svc.ServiceContext) *AuthLogic {
	ctx := r.Context()
	return &AuthLogic{svcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx), request: r}
}

func (l *AuthLogic) Auth(in *pb.Requester) *pb.CommonResp {
	inputToken := in.Token
	if inputToken == "" {
		// 判断接口是否需要登录
		if strings.Contains(l.request.URL.Path, "white") {
			return &pb.CommonResp{}
		} else {
			return pb.NewAuthErrorResp("请先登录")
		}
	}
	// 验证token
	code, msg := xjwt.VerifyToken(l.ctx, l.svcCtx.Redis(), in.Id, inputToken, xjwt.WithPlatform(in.Platform), xjwt.WithDeviceId(in.DeviceId))
	switch code {
	case xjwt.VerifyTokenCodeOK:
		return &pb.CommonResp{}
	case xjwt.VerifyTokenCodeInternalError:
		return pb.NewInternalErrorResp(msg)
	case xjwt.VerifyTokenCodeError, xjwt.VerifyTokenCodeExpire, xjwt.VerifyTokenCodeBaned, xjwt.VerifyTokenCodeReplace:
		if strings.Contains(l.request.URL.Path, "white") {
			return &pb.CommonResp{}
		}
		return pb.NewAlertErrorResp("下线通知", msg)
	default:
		return pb.NewAuthErrorResp(msg)
	}
}
