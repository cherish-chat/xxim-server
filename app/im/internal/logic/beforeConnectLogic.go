package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"time"

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
	// 查询user是否被禁用
	detail, err := l.svcCtx.UserService().GetUserModelDetail(l.ctx, &pb.GetUserModelDetailReq{Id: in.ConnParam.UserId})
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.BeforeConnectResp{Msg: "连接失败"}, err
	}
	if detail.UserModel.UnblockTime > 0 {
		// 和当前时间比较
		if detail.UserModel.UnblockTime > time.Now().UnixMilli() {
			l.Errorf("用户被禁用: %v", err)
			return &pb.BeforeConnectResp{Msg: "您的账号违反了相关规定，已被禁用"}, err
		}
	}
	// 验证token
	code, msg := xjwt.VerifyToken(l.ctx, l.svcCtx.Redis(), in.ConnParam.UserId, in.ConnParam.Token, xjwt.WithPlatform(in.ConnParam.Platform), xjwt.WithDeviceId(in.ConnParam.DeviceId))
	switch code {
	case xjwt.VerifyTokenCodeInternalError, xjwt.VerifyTokenCodeError, xjwt.VerifyTokenCodeExpire, xjwt.VerifyTokenCodeBaned, xjwt.VerifyTokenCodeReplace:
		return &pb.BeforeConnectResp{Msg: msg}, nil
	default:
		return &pb.BeforeConnectResp{}, nil
	}
}
