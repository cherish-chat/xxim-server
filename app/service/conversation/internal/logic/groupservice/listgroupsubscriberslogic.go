package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGroupSubscribersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGroupSubscribersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupSubscribersLogic {
	return &ListGroupSubscribersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListGroupSubscribersLogic) ListGroupSubscribers(in *peerpb.ListGroupSubscribersReq) (*peerpb.ListGroupSubscribersResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.ListGroupSubscribersResp{}, nil
}
