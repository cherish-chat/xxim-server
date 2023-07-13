package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOnlineLogic {
	return &UserAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAfterOnlineLogic) UserAfterOnline(in *peerpb.UserAfterOnlineReq) (*peerpb.UserAfterOnlineResp, error) {
	//1. 订阅号的逻辑
	{
		_, err := l.svcCtx.ChannelService.ChannelAfterOnline(context.Background(), &peerpb.ChannelAfterOnlineReq{
			Header: in.Header,
		})
		if err != nil {
			l.Errorf("channel after online error: %v", err)
		}
	}
	return &peerpb.UserAfterOnlineResp{}, nil
}
