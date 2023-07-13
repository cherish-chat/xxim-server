package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterKeepAliveLogic {
	return &UserAfterKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAfterKeepAliveLogic) UserAfterKeepAlive(in *peerpb.UserAfterKeepAliveReq) (*peerpb.UserAfterKeepAliveResp, error) {
	// 1. 订阅群聊消息
	{
		_, err := l.svcCtx.GroupService.GroupAfterKeepAlive(context.Background(), &peerpb.GroupAfterKeepAliveReq{
			Header: in.Header,
		})
		if err != nil {
			l.Errorf("subscribe group message error: %v", err)
		}
	}
	// 2. 订阅 订阅号 消息
	{
		_, err := l.svcCtx.ChannelService.ChannelAfterKeepAlive(context.Background(), &peerpb.ChannelAfterKeepAliveReq{
			Header: in.Header,
		})
		if err != nil {
			l.Errorf("subscribe subscription message error: %v", err)
		}
	}
	return &peerpb.UserAfterKeepAliveResp{}, nil
}
