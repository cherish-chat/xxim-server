package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
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
	ip := in.ConnParam.Ips
	for _, v := range allIpBlackList {
		if v[2] == in.ConnParam.UserId || v[2] == "" {
			if strings.Compare(ip, v[0]) >= 0 && strings.Compare(ip, v[1]) <= 0 {
				return &pb.BeforeConnectResp{}, status.Error(codes.Unauthenticated, "ip被封禁")
			}
		}
	}
	// 查询user是否被禁用
	detail, err := l.svcCtx.UserService().GetUserModelDetail(l.ctx, &pb.GetUserModelDetailReq{Id: in.ConnParam.UserId})
	if err != nil {
		l.Infof("查询用户失败: %v", err)
		return &pb.BeforeConnectResp{Msg: "连接失败"}, status.Error(codes.Unauthenticated, "用户被删除")
	}
	if detail.UserModel.UnblockTime > 0 {
		// 和当前时间比较
		if detail.UserModel.UnblockTime > time.Now().UnixMilli() {
			l.Errorf("用户被禁用: %v", err)
			return &pb.BeforeConnectResp{Msg: "您的账号违反了相关规定，已被禁用"}, status.Error(codes.Unauthenticated, "用户被禁用")
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
