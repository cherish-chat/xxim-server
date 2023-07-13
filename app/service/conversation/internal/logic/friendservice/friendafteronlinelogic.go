package friendservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendAfterOnlineLogic {
	return &FriendAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendAfterOnlineLogic) FriendAfterOnline(in *peerpb.FriendAfterOnlineReq) (*peerpb.FriendAfterOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.FriendAfterOnlineResp{}, nil
}
