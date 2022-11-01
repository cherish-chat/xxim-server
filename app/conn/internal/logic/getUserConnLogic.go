package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConnLogic {
	return &GetUserConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConnLogic) GetUserConn(in *pb.GetUserConnReq) (*pb.GetUserConnResp, error) {
	conns := GetConnLogic().GetConnsByFilter(GetConnLogic().BuildSearchUserConnFilter(in))
	var resp []*pb.ConnParam
	for _, conn := range conns {
		resp = append(resp, &pb.ConnParam{
			UserId:      conn.ConnParam.UserId,
			Token:       conn.ConnParam.Token,
			DeviceId:    conn.ConnParam.DeviceId,
			Platform:    conn.ConnParam.Platform,
			Ips:         conn.ConnParam.Ips,
			NetworkUsed: conn.ConnParam.NetworkUsed,
			Headers:     conn.ConnParam.Headers,
			PodIp:       l.svcCtx.PodIp,
		})
	}
	return &pb.GetUserConnResp{ConnParams: resp}, nil
}
