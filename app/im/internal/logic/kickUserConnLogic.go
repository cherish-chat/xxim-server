package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickUserConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickUserConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickUserConnLogic {
	return &KickUserConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickUserConnLogic) KickUserConn(in *pb.KickUserConnReq) (*pb.KickUserConnResp, error) {
	fs := make([]func() error, 0)
	for _, pod := range l.svcCtx.ConnPodsMgr.AllConnServices() {
		fs = append(fs, func() error {
			_, err := pod.KickUserConn(l.ctx, in)
			return err
		})
	}
	err := mr.Finish(fs...)
	if err != nil {
		l.Errorf("KickUserConn failed, err: %v", err)
		return &pb.KickUserConnResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.KickUserConnResp{}, nil
}
