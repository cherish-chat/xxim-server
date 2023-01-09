package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/propagation"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
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
	var respList []*pb.GetUserConnResp
	fs := make([]func() error, 0)
	for _, pod := range l.svcCtx.ConnPodsMgr.AllConnServices() {
		podValue := *pod
		fs = append(fs, func() error {
			var resp *pb.GetUserConnResp
			var err error
			xtrace.StartFuncSpan(l.ctx, "GetUserConn", func(ctx context.Context) {
				resp, err = podValue.GetUserConn(l.ctx, in)
			}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
				"podIpPort": podValue.PodIpPort,
			}))
			if err == nil {
				respList = append(respList, resp)
			}
			return err
		})
	}
	err := mr.Finish(fs...)
	if err != nil {
		l.Errorf("GetUserConn failed, err: %v", err)
		return &pb.GetUserConnResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.ConnParam
	for _, r := range respList {
		resp = append(resp, r.ConnParams...)
	}
	return &pb.GetUserConnResp{ConnParams: resp}, nil
}
