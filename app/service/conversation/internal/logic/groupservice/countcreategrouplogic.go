package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountCreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountCreateGroupLogic {
	return &CountCreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CountCreateGroupLogic) CountCreateGroup(in *peerpb.CountCreateGroupReq) (*peerpb.CountCreateGroupResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.CountCreateGroupResp{}, nil
}
