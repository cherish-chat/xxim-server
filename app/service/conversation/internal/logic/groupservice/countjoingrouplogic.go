package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountJoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountJoinGroupLogic {
	return &CountJoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CountJoinGroupLogic) CountJoinGroup(in *peerpb.CountJoinGroupReq) (*peerpb.CountJoinGroupResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.CountJoinGroupResp{}, nil
}
