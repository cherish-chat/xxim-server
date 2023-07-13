package friendservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendAfterOfflineLogic {
	return &FriendAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendAfterOfflineLogic) FriendAfterOffline(in *peerpb.FriendAfterOfflineReq) (*peerpb.FriendAfterOfflineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.FriendAfterOfflineResp{}, nil
}
