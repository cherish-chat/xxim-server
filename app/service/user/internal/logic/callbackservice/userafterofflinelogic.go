package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOfflineLogic {
	return &UserAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAfterOfflineLogic) UserAfterOffline(in *peerpb.UserAfterOfflineReq) (*peerpb.UserAfterOfflineResp, error) {
	//1. 订阅号的逻辑
	{
		_, err := l.svcCtx.ChannelService.ChannelAfterOffline(context.Background(), &peerpb.ChannelAfterOfflineReq{
			UserId: in.UserId,
		})
		if err != nil {
			l.Errorf("channel after online error: %v", err)
		}
	}
	return &peerpb.UserAfterOfflineResp{}, nil
}
