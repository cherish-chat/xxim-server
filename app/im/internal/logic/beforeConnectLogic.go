package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xjwt"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BeforeConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBeforeConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BeforeConnectLogic {
	return &BeforeConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BeforeConnectLogic) BeforeConnect(in *pb.BeforeConnectReq) (*pb.BeforeConnectResp, error) {
	// 验证token
	code, msg := xjwt.VerifyToken(l.ctx, l.svcCtx.Redis(), in.ConnParam.UserId, in.ConnParam.Token, xjwt.WithPlatform(in.ConnParam.Platform), xjwt.WithDeviceId(in.ConnParam.DeviceId))
	switch code {
	case xjwt.VerifyTokenCodeInternalError, xjwt.VerifyTokenCodeError, xjwt.VerifyTokenCodeExpire, xjwt.VerifyTokenCodeBaned, xjwt.VerifyTokenCodeReplace:
		return &pb.BeforeConnectResp{Msg: msg}, nil
	default:
		return &pb.BeforeConnectResp{}, nil
	}
}
