package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type connMap map[string]deviceMap // key: platform, value: deviceMap

type deviceMap map[string]*types.UserConn // key: deviceId, value: *types.UserConn

type ServerLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	userConnMapLock sync.RWMutex
	userConnMap     map[string]connMap // key: userId, value: connMap
}

var singletonServerLogic *ServerLogic

func GetServerLogic(svcCtx *svc.ServiceContext) *ServerLogic {
	if singletonServerLogic == nil {
		l := &ServerLogic{svcCtx: svcCtx}
		l.Logger = logx.WithContext(context.Background())
		l.userConnMap = map[string]connMap{}
		singletonServerLogic = l
	}
	return singletonServerLogic
}

func (l *ServerLogic) BeforeConnect(ctx context.Context, param types.ConnParam) (int, error) {
	// TODO handle before connect
	return 0, nil
}

func (l *ServerLogic) AddSubscriber(c *types.UserConn) {
	l.Infof("user %s connected", utils.AnyToString(c.ConnParam))
	// 加入用户连接
	l.userConnMapLock.Lock()
	if _, ok := l.userConnMap[c.ConnParam.UserId]; !ok {
		l.userConnMap[c.ConnParam.UserId] = connMap{}
	}
	if _, ok := l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform]; !ok {
		l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform] = deviceMap{}
	}
	if _, ok := l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform][c.ConnParam.DeviceId]; ok {
		l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform][c.ConnParam.DeviceId].Conn.Close(int(websocket.StatusNormalClosure), "duplicate connection")
	}
	l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform][c.ConnParam.DeviceId] = c
	l.userConnMapLock.Unlock()
	// 告知客户端连接成功
	c.Conn.Write(context.Background(), int(websocket.MessageText), []byte("connected"))
	l.stats()
}

func (l *ServerLogic) DeleteSubscriber(c *types.UserConn) {
	l.Infof("user %s disconnected", utils.AnyToString(c.ConnParam))
	// 删除用户连接
	l.userConnMapLock.Lock()
	if _, ok := l.userConnMap[c.ConnParam.UserId]; !ok {
		return
	}
	if _, ok := l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform]; !ok {
		return
	}
	if _, ok := l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform][c.ConnParam.DeviceId]; !ok {
		return
	}
	delete(l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform], c.ConnParam.DeviceId)
	if len(l.userConnMap[c.ConnParam.UserId][c.ConnParam.Platform]) == 0 {
		delete(l.userConnMap[c.ConnParam.UserId], c.ConnParam.Platform)
	}
	if len(l.userConnMap[c.ConnParam.UserId]) == 0 {
		delete(l.userConnMap, c.ConnParam.UserId)
	}
	l.userConnMapLock.Unlock()
	l.stats()
}

func (l *ServerLogic) Stats() {
	l.stats()
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			l.stats()
		}
	}
}

func (l *ServerLogic) stats() {
	l.userConnMapLock.RLock()
	defer l.userConnMapLock.RUnlock()
	// 统计在线用户数
	onlineUserCount := len(l.userConnMap)
	// 统计在线设备数
	onlineDeviceCount := 0
	for _, cm := range l.userConnMap {
		for _, dm := range cm {
			onlineDeviceCount += len(dm)
		}
	}
	l.Infof("online user count: %d, online device count: %d", onlineUserCount, onlineDeviceCount)
}
