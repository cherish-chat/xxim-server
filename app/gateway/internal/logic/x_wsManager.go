package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"nhooyr.io/websocket"
	"sync"
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
	w.wsConnectionMap.Set(id, wsConnection)
	return wsConnection, nil
}

func (w *wsManager) RemoveSubscriber(header *pb.RequestHeader, id int64) error {
	w.wsConnectionMap.Delete(header.UserId, id)
	return nil
}
