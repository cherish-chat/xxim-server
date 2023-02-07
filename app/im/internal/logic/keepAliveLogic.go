package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeepAliveLogic {
	return &KeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KeepAliveLogic) KeepAlive(in *pb.KeepAliveReq) (*pb.KeepAliveResp, error) {
	var fs []func() error
	fs = append(fs, func() error {
		var err error
		_, err = l.svcCtx.MsgService().KeepAlive(l.ctx, in)
		return err
	})
	fs = append(fs, func() error {
		var err error
		_, err = l.svcCtx.NoticeService().KeepAlive(l.ctx, in)
		return err
	})
	err := mr.Finish(fs...)
	if err != nil {
		return &pb.KeepAliveResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.KeepAliveResp{CommonResp: pb.NewSuccessResp()}, nil
}
