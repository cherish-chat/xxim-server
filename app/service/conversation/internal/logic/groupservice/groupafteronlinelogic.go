package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAfterOnlineLogic {
	return &GroupAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupAfterOnlineLogic) GroupAfterOnline(in *peerpb.GroupAfterOnlineReq) (*peerpb.GroupAfterOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GroupAfterOnlineResp{}, nil
}
