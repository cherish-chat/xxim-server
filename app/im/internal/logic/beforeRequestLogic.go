package logic

import (
	"context"
	"strings"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BeforeRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBeforeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BeforeRequestLogic {
	return &BeforeRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BeforeRequestLogic) BeforeRequest(in *pb.BeforeRequestReq) (*pb.BeforeRequestResp, error) {
	if strings.Contains(in.Method, "/white/") {
		return &pb.BeforeRequestResp{
			CommonResp: pb.NewSuccessResp(),
		}, nil
	}
	resp, _ := NewBeforeConnectLogic(l.ctx, l.svcCtx).BeforeConnect(&pb.BeforeConnectReq{ConnParam: &pb.ConnParam{
		UserId:      in.CommonReq.UserId,
		Token:       in.CommonReq.Token,
		DeviceId:    in.CommonReq.DeviceId,
		Platform:    in.CommonReq.Platform,
		Ips:         in.CommonReq.Ip,
		DeviceModel: in.CommonReq.DeviceModel,
		OsVersion:   in.CommonReq.OsVersion,
		AppVersion:  in.CommonReq.AppVersion,
		Language:    in.CommonReq.Language,
	}})
	if resp.Code == 0 {
		return &pb.BeforeRequestResp{
			CommonResp: pb.NewSuccessResp(),
		}, nil
	}
	return &pb.BeforeRequestResp{CommonResp: pb.NewAuthErrorResp(resp.Msg)}, nil
}
