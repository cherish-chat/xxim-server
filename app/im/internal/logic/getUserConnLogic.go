package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"

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
		fs = append(fs, func() error {
			resp, err := pod.GetUserConn(l.ctx, in)
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
