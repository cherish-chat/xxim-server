package friendservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendApplyLogic {
	return &ListFriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFriendApplyLogic) ListFriendApply(in *peerpb.ListFriendApplyReq) (*peerpb.ListFriendApplyResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.ListFriendApplyResp{}, nil
}
