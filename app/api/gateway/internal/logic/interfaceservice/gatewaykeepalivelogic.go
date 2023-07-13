package interfaceservicelogic

import (
	"context"
	"sync"
	"time"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayKeepAliveLogic {
	l := &GatewayKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
	if userKeepAliveMap == nil {
		userKeepAliveMap = new(sync.Map)
		go l.checkTimer()
	}
	return l
}

func (l *GatewayKeepAliveLogic) GatewayKeepAlive(in *peerpb.GatewayKeepAliveReq) (*peerpb.GatewayKeepAliveResp, error) {
	_, err := l.svcCtx.CallbackService.UserAfterKeepAlive(context.Background(), &peerpb.UserAfterKeepAliveReq{
		Header: in.Header,
	})
	if err != nil {
		l.Errorf("UserAfterKeepAlive error: %v", err)
	}
	if _, ok := userKeepAliveMap.Load(in.Header.UserId); !ok {
		_, _ = l.svcCtx.CallbackService.UserAfterOnline(context.Background(), &peerpb.UserAfterOnlineReq{Header: in.Header})
	}
	userKeepAliveMap.Store(in.Header.UserId, time.Now())

	return &peerpb.GatewayKeepAliveResp{}, nil
}

var userKeepAliveMap *sync.Map // key: userId, value: time.Time

func (l *GatewayKeepAliveLogic) checkTimer() {
	ticker := time.NewTicker(time.Second * time.Duration(60))
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			userKeepAliveMap.Range(func(key, value interface{}) bool {
				userId := key.(string)
				lastTime := value.(time.Time)
				if now.Sub(lastTime).Seconds() > float64(60) {
					// 用户下线
					_, _ = l.svcCtx.CallbackService.UserAfterOffline(context.Background(), &peerpb.UserAfterOfflineReq{UserId: userId})
					userKeepAliveMap.Delete(userId)
				}
				return true
			})
		}
	}
}
