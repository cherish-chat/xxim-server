package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

func InitWsManager(svcCtx *svc.ServiceContext) {
	WsManager = &wsManager{
		svcCtx: svcCtx,
		wsConnectionMap: &WsConnectionMap{
			userIdsMap:          make(map[string][]*WsConnection),
			userIdsMapLock:      sync.RWMutex{},
			idConnectionMap:     make(map[int64]*WsConnection),
			idConnectionMapLock: sync.RWMutex{},
			idAliveTimeMap:      make(map[int64]time.Time),
			idAliveTimeMapLock:  sync.RWMutex{},
		},
	}
	go WsManager.loopCheck()
}

type WsConnection struct {
	Id         int64
	Connection *websocket.Conn
	Header     *pb.RequestHeader
	Ctx        context.Context
}

func (c *WsConnection) ToPb() *pb.WsConnection {
	return &pb.WsConnection{
		Id:     c.Id,
		Header: c.Header,
	}
}

type WsConnectionMap struct {
	idConnectionMap     map[int64]*WsConnection
	idConnectionMapLock sync.RWMutex
	userIdsMap          map[string][]*WsConnection
	userIdsMapLock      sync.RWMutex
	idAliveTimeMap      map[int64]time.Time
	idAliveTimeMapLock  sync.RWMutex
}

func (w *WsConnectionMap) GetByConnectionId(connectionId int64) (*WsConnection, bool) {
	// RLock() 读锁
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	v, ok := w.idConnectionMap[connectionId]
	return v, ok
}

func (w *WsConnectionMap) GetByUserId(userId string) ([]*WsConnection, bool) {
	w.userIdsMapLock.RLock()
	defer w.userIdsMapLock.RUnlock()
	v, ok := w.userIdsMap[userId]
	return v, ok
}

func (w *WsConnectionMap) GetByUserIds(userIds []string) []*WsConnection {
	w.userIdsMapLock.RLock()
	defer w.userIdsMapLock.RUnlock()
	var connections []*WsConnection
	for _, userId := range userIds {
		v, ok := w.userIdsMap[userId]
		if ok {
			connections = append(connections, v...)
		}
	}
	return connections
}

func (w *WsConnectionMap) GetAll() []*WsConnection {
	var connections []*WsConnection
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	for _, v := range w.idConnectionMap {
		connections = append(connections, v)
	}
	return connections
}

func (w *WsConnectionMap) GetAliveTime(connectionId int64) (time.Time, bool) {
	// RLock() 读锁
	w.idAliveTimeMapLock.RLock()
	defer w.idAliveTimeMapLock.RUnlock()
	v, ok := w.idAliveTimeMap[connectionId]
	return v, ok
}

func (w *WsConnectionMap) SetAliveTime(ctx context.Context, connectionId int64, aliveTime time.Time) {
	w.idAliveTimeMapLock.Lock()
	w.idAliveTimeMap[connectionId] = aliveTime
	w.idAliveTimeMapLock.Unlock()
}

func (w *WsConnectionMap) Set(connectionId int64, value *WsConnection) {
	w.idConnectionMapLock.Lock()
	w.idConnectionMapLock.Unlock()
	w.userIdsMapLock.Lock()
	w.userIdsMap[value.Header.UserId] = append(w.userIdsMap[value.Header.UserId], value)
	w.userIdsMapLock.Unlock()
	w.idAliveTimeMapLock.Lock()
	w.idAliveTimeMap[connectionId] = time.Now()
	w.idAliveTimeMapLock.Unlock()
}

func (w *WsConnectionMap) Delete(userId string, connectionId int64) {
	w.idConnectionMapLock.Lock()
	delete(w.idConnectionMap, connectionId)
	w.idConnectionMapLock.Unlock()
	w.userIdsMapLock.Lock()
	defer w.userIdsMapLock.Unlock()
	//获取用户的所有连接
	connections, ok := w.userIdsMap[userId]
	if !ok {
		return
	}
	//删除用户的某个连接
	var newConnections []*WsConnection
	for _, connection := range connections {
		if connection.Id != connectionId {
			newConnections = append(newConnections, connection)
		}
	}
	w.userIdsMap[userId] = newConnections
	w.idAliveTimeMapLock.Lock()
	delete(w.idAliveTimeMap, connectionId)
	w.idAliveTimeMapLock.Unlock()
}

type wsManager struct {
	svcCtx          *svc.ServiceContext
	wsConnectionMap *WsConnectionMap
}

var WsManager *wsManager

func (w *wsManager) AddSubscriber(ctx context.Context, header *pb.RequestHeader, connection *websocket.Conn, id int64) (*WsConnection, error) {
	wsConnection := &WsConnection{
		Id:         id,
		Connection: connection,
		Header:     header,
		Ctx:        ctx,
	}
	//启动定时器 定时删掉连接
	go w.clearConnectionTimer(wsConnection)
	w.wsConnectionMap.Set(id, wsConnection)
	go func() {
		_, e := w.svcCtx.UserService.UserAfterOnline(ctx, &pb.UserAfterOnlineReq{Header: header})
		if e != nil {
			logx.Errorf("UserAfterOnline error: %s", e.Error())
		}
	}()
	return wsConnection, nil
}

func (w *wsManager) RemoveSubscriber(header *pb.RequestHeader, id int64, closeCode websocket.StatusCode, closeReason string) error {
	connection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if ok {
		_ = connection.Connection.Close(closeCode, closeReason)
	}
	w.wsConnectionMap.Delete(header.UserId, id)
	go func() {
		if ok {
			_, e := w.svcCtx.UserService.UserAfterOffline(connection.Ctx, &pb.UserAfterOfflineReq{Header: header})
			if e != nil {
				logx.Errorf("UserAfterOffline error: %s", e.Error())
			}
		}
	}()
	return nil
}

// clearConnectionTimer 定时器清除连接
func (w *wsManager) clearConnectionTimer(connection *WsConnection) {
	ticker := time.NewTicker(time.Second * time.Duration(w.svcCtx.Config.Websocket.KeepAliveTickerSecond))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//使用 id 查询连接最后活跃时间
			aliveTime, ok := w.wsConnectionMap.GetAliveTime(connection.Id)
			sub := time.Now().Sub(aliveTime)
			if !ok || sub > time.Second*time.Duration(w.svcCtx.Config.Websocket.KeepAliveSecond) {
				// 删除连接
				w.RemoveSubscriber(connection.Header, connection.Id, websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout), "heartbeat timeout")
				return
			}
		}
	}
}

func (w *wsManager) KeepAlive(ctx context.Context, connection *WsConnection) {
	w.wsConnectionMap.SetAliveTime(ctx, connection.Id, time.Now())
}

func (w *wsManager) WriteData(id int64, data []byte) bool {
	wsConnection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if !ok {
		return false
	}
	err := wsConnection.Connection.Write(wsConnection.Ctx, websocket.MessageText, data)
	if err != nil {
		return false
	}
	return true
}

func (w *wsManager) CloseConnection(id int64, code websocket.StatusCode, reason string) {
	wsConnection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if !ok {
		return
	}
	_ = wsConnection.Connection.Close(code, reason)
}

func (w *wsManager) loopCheck() {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(w.svcCtx.Config.Websocket.KeepAliveTickerSecond))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				connections := w.wsConnectionMap.GetAll()
				for _, connection := range connections {
					_, ok := w.wsConnectionMap.GetAliveTime(connection.Id)
					if !ok {
						// 删除连接
						w.RemoveSubscriber(connection.Header, connection.Id, websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout), "heartbeat timeout")
					}
				}
			}
		}
	}()
}
