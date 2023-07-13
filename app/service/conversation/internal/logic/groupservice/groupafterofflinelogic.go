package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAfterOfflineLogic {
	return &GroupAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupAfterOfflineLogic) GroupAfterOffline(in *peerpb.GroupAfterOfflineReq) (*peerpb.GroupAfterOfflineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GroupAfterOfflineResp{}, nil
}
