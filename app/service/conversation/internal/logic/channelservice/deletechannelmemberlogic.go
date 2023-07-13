package channelservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteChannelMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteChannelMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteChannelMemberLogic {
	return &DeleteChannelMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteChannelMemberLogic) DeleteChannelMember(in *peerpb.DeleteChannelMemberReq) (*peerpb.DeleteChannelMemberResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.DeleteChannelMemberResp{}, nil
}
