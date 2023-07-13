package sessionservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListJoinedConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListJoinedConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJoinedConversationsLogic {
	return &ListJoinedConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListJoinedConversationsLogic) ListJoinedConversations(in *peerpb.ListJoinedConversationsReq) (*peerpb.ListJoinedConversationsResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.ListJoinedConversationsResp{}, nil
}
