package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlushShieldWordTireTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlushShieldWordTireTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlushShieldWordTireTreeLogic {
	return &FlushShieldWordTireTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FlushShieldWordTireTree 刷新屏蔽词
func (l *FlushShieldWordTireTreeLogic) FlushShieldWordTireTree(in *pb.FlushShieldWordTireTreeReq) (*pb.FlushShieldWordTireTreeResp, error) {
	ShieldWordTrieTreeInstance.Flush()
	return &pb.FlushShieldWordTireTreeResp{}, nil
}
