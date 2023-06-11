package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

func InitWsManager(svcCtx *svc.ServiceContext) {
	WsManager = &wsManager{
		svcCtx: svcCtx,
		wsConnectionMap: &WsConnectionRWMutexMap{
			idConnectionMap:     make(map[int64]*WsConnection),
			userIdsMap:          make(map[string][]*WsConnection),
			idConnectionMapLock: sync.RWMutex{},
		},
	}
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

type WsConnectionRWMutexMap struct {
	idConnectionMap     map[int64]*WsConnection
	idConnectionMapLock sync.RWMutex
	userIdsMap          map[string][]*WsConnection
	userIdsMapLock      sync.RWMutex
	idAliveMap          sync.Map
}

func (w *WsConnectionRWMutexMap) GetByConnectionId(connectionId int64) (*WsConnection, bool) {
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	v, ok := w.idConnectionMap[connectionId]
	return v, ok
}

func (w *WsConnectionRWMutexMap) GetByUserId(userId string) ([]*WsConnection, bool) {
	w.userIdsMapLock.RLock()
	defer w.userIdsMapLock.RUnlock()
	v, ok := w.userIdsMap[userId]
	return v, ok
}

func (w *WsConnectionRWMutexMap) GetByUserIds(userIds []string) []*WsConnection {
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

func (w *WsConnectionRWMutexMap) GetAll() []*WsConnection {
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	var connections []*WsConnection
	for _, connection := range w.idConnectionMap {
		connections = append(connections, connection)
	}
	return connections
}

func (w *WsConnectionRWMutexMap) GetAliveTime(connectionId int64) (time.Time, bool) {
	v, ok := w.idAliveMap.Load(connectionId)
	if !ok {
		return time.Time{}, false
	}
	return v.(time.Time), true
}

func (w *WsConnectionRWMutexMap) DeleteAliveTime(connectionId int64) {
	w.idAliveMap.Delete(connectionId)
}

func (w *WsConnectionRWMutexMap) SetAliveTime(connectionId int64, aliveTime time.Time) {
	w.idAliveMap.Store(connectionId, aliveTime)
}

func (w *WsConnectionRWMutexMap) Set(connectionId int64, value *WsConnection) {
	w.idConnectionMapLock.Lock()
	w.idConnectionMap[connectionId] = value
	w.idConnectionMapLock.Unlock()
	w.userIdsMapLock.Lock()
	w.userIdsMap[value.Header.UserId] = append(w.userIdsMap[value.Header.UserId], value)
	w.userIdsMapLock.Unlock()
}

func (w *WsConnectionRWMutexMap) Delete(userId string, connectionId int64) {
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
	w.DeleteAliveTime(connectionId)
}

type wsManager struct {
	svcCtx          *svc.ServiceContext
	wsConnectionMap *WsConnectionRWMutexMap
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
	w.KeepAlive(wsConnection)
	w.wsConnectionMap.Set(id, wsConnection)
	return wsConnection, nil
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
			if !ok || time.Now().Sub(aliveTime) > time.Second*time.Duration(w.svcCtx.Config.Websocket.KeepAliveSecond) {
				// 删除连接
				w.RemoveSubscriber(connection.Header, connection.Id, websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout), "heartbeat timeout")
				return
			}
		}
	}
}

func (w *wsManager) KeepAlive(connection *WsConnection) {
	w.wsConnectionMap.SetAliveTime(connection.Id, time.Now())
}

func (w *wsManager) RemoveSubscriber(header *pb.RequestHeader, id int64, closeCode websocket.StatusCode, closeReason string) error {
	connection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if ok {
		_ = connection.Connection.Close(closeCode, closeReason)
	}
	w.wsConnectionMap.Delete(header.UserId, id)
	return nil
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
