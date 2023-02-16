package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TimerCleanSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTimerCleanSubscriptionLogic(svcCtx *svc.ServiceContext) *TimerCleanSubscriptionLogic {
	l := &TimerCleanSubscriptionLogic{svcCtx: svcCtx, ctx: context.Background()}
	l.Logger = logx.WithContext(l.ctx)
	return l
}

func (l *TimerCleanSubscriptionLogic) Start() {
	// 每天执行一次 删除前一天之前的垃圾数据
	// 1. 删除过期的订阅关系
	l.cleanSubscription()
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-ticker.C:
			l.cleanSubscription()
		}
	}
}

func (l *TimerCleanSubscriptionLogic) cleanSubscription() {
	regKey := rediskey.ConvMembersSubscribed("*")
	keysChan := make(chan []string)
	stopChan := make(chan struct{})
	go xredis.Scan(l.svcCtx.RedisSub(), l.ctx, regKey, keysChan, stopChan)
	var keys []string
GETKEY:
	for {
		select {
		case ks := <-keysChan:
			keys = append(keys, ks...)
		case <-stopChan:
			break GETKEY
		}
	}
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		// 	ZREMRANGEBYSCORE key 0 time.Now().Add(time.Second * 60 * 3).UnixMilli()
		_, err := l.svcCtx.RedisSub().ZremrangebyscoreCtx(l.ctx, key, 0, time.Now().Add(time.Second*60*3).UnixMilli())
		if err != nil {
			l.Errorf("ZremrangebyscoreCtx error: %s", err.Error())
		}
	}
}
