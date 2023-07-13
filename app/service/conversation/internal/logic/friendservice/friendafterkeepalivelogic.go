package friendservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendAfterKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendAfterKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendAfterKeepAliveLogic {
	return &FriendAfterKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendAfterKeepAliveLogic) FriendAfterKeepAlive(in *peerpb.FriendAfterKeepAliveReq) (*peerpb.FriendAfterKeepAliveResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.FriendAfterKeepAliveResp{}, nil
}
