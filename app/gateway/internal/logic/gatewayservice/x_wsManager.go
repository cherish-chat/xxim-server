package gatewayservicelogic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

func InitWsManager(svcCtx *svc.ServiceContext) {
	WsManager = &wsManager{
		svcCtx: svcCtx,
		wsConnectionMap: &WsConnectionMap{
			userIdsMap:          make(map[string][]*UniversalConnection),
			userIdsMapLock:      sync.RWMutex{},
			idConnectionMap:     make(map[int64]*UniversalConnection),
			idConnectionMapLock: sync.RWMutex{},
			idAliveTimeMap:      make(map[int64]time.Time),
			idAliveTimeMapLock:  sync.RWMutex{},
		},
	}
	go WsManager.loopCheck()
}

type Connection interface {
	Close(code pb.WebsocketCustomCloseCode, reason string) error
	Write(ctx context.Context, data []byte) error
}

type WsForConnection struct {
	Connection *websocket.Conn
}

func (w *WsForConnection) Close(code pb.WebsocketCustomCloseCode, reason string) error {
	return w.Connection.Close(websocket.StatusCode(code.Number()), reason)
}

func (w *WsForConnection) Write(ctx context.Context, data []byte) error {
	return w.Connection.Write(ctx, websocket.MessageBinary, data)
}

func NewWsForConnection(c *websocket.Conn) *WsForConnection {
	return &WsForConnection{Connection: c}
}

type RtcForConnection struct {
	DataChannel *webrtc.DataChannel
}

func (r *RtcForConnection) Close(code pb.WebsocketCustomCloseCode, reason string) error {
	return r.DataChannel.Close()
}

func (r *RtcForConnection) Write(ctx context.Context, data []byte) error {
	return r.DataChannel.Send(data)
}

func NewRtcForConnection(c *webrtc.DataChannel) *RtcForConnection {
	return &RtcForConnection{DataChannel: c}
}

type UniversalConnection struct {
	Id              int64
	Connection      Connection
	Header          *pb.RequestHeader
	Ctx             context.Context
	ReSetHeaderLock sync.RWMutex
}

func (c *UniversalConnection) ToPb() *pb.WsConnection {
	return &pb.WsConnection{
		Id:     c.Id,
		Header: c.GetHeader(),
	}
}

func (c *UniversalConnection) ReSetHeader(header *pb.RequestHeader) *pb.RequestHeader {
	c.ReSetHeaderLock.Lock()
	defer c.ReSetHeaderLock.Unlock()
	old := c.Header
	copyHeader := &pb.RequestHeader{
		AppId:        header.AppId,
		UserId:       header.UserId,
		UserToken:    header.UserToken,
		ClientIp:     header.ClientIp,
		InstallId:    header.InstallId,
		Platform:     header.Platform,
		GatewayPodIp: header.GatewayPodIp,
		DeviceModel:  header.DeviceModel,
		OsVersion:    header.OsVersion,
		AppVersion:   header.AppVersion,
		Language:     header.Language,
		ConnectTime:  old.ConnectTime,
		Encoding:     old.Encoding,
		Extra:        header.Extra,
	}
	c.Header = copyHeader

	go func() {
		if header.UserId != "" {
			_, e := WsManager.svcCtx.CallbackService.UserAfterOnline(c.Ctx, &pb.UserAfterOnlineReq{Header: copyHeader})
			if e != nil {
				logx.Errorf("UserAfterOnline error: %s", e.Error())
			}
		} else {
			_, e := WsManager.svcCtx.CallbackService.UserAfterOffline(c.Ctx, &pb.UserAfterOfflineReq{Header: old})
			if e != nil {
				logx.Errorf("UserAfterOnline error: %s", e.Error())
			}
		}
	}()
	return copyHeader
}

func (c *UniversalConnection) GetHeader() *pb.RequestHeader {
	c.ReSetHeaderLock.RLock()
	defer c.ReSetHeaderLock.RUnlock()
	return c.Header
}

type WsConnectionMap struct {
	idConnectionMap     map[int64]*UniversalConnection
	idConnectionMapLock sync.RWMutex
	userIdsMap          map[string][]*UniversalConnection
	userIdsMapLock      sync.RWMutex
	idAliveTimeMap      map[int64]time.Time
	idAliveTimeMapLock  sync.RWMutex
}

func (w *WsConnectionMap) GetByConnectionId(connectionId int64) (*UniversalConnection, bool) {
	// RLock() 读锁
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	v, ok := w.idConnectionMap[connectionId]
	return v, ok
}

func (w *WsConnectionMap) GetByUserId(userId string) ([]*UniversalConnection, bool) {
	w.userIdsMapLock.RLock()
	defer w.userIdsMapLock.RUnlock()
	v, ok := w.userIdsMap[userId]
	return v, ok
}

func (w *WsConnectionMap) GetByUserIds(userIds []string) []*UniversalConnection {
	w.userIdsMapLock.RLock()
	defer w.userIdsMapLock.RUnlock()
	var connections []*UniversalConnection
	for _, userId := range userIds {
		v, ok := w.userIdsMap[userId]
		if ok {
			connections = append(connections, v...)
		}
	}
	return connections
}

func (w *WsConnectionMap) GetAll() []*UniversalConnection {
	var connections []*UniversalConnection
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

func (w *WsConnectionMap) Set(connectionId int64, value *UniversalConnection) {
	w.idConnectionMapLock.Lock()
	w.idConnectionMap[connectionId] = value
	w.idConnectionMapLock.Unlock()
	userId := value.GetHeader().UserId
	if userId != "" {
		w.userIdsMapLock.Lock()
		w.userIdsMap[userId] = append(w.userIdsMap[userId], value)
		w.userIdsMapLock.Unlock()
	}
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
	var newConnections []*UniversalConnection
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

func (w *wsManager) AddSubscriber(ctx context.Context, header *pb.RequestHeader, connection Connection, id int64) (*UniversalConnection, error) {
	wsConnection := &UniversalConnection{
		Id:         id,
		Connection: connection,
		Header:     header,
		Ctx:        ctx,
	}
	//启动定时器 定时删掉连接
	go w.clearConnectionTimer(wsConnection)
	w.wsConnectionMap.Set(id, wsConnection)
	go func() {
		if header.UserId != "" {
			_, e := w.svcCtx.CallbackService.UserAfterOnline(ctx, &pb.UserAfterOnlineReq{Header: header})
			if e != nil {
				logx.Errorf("UserAfterOnline error: %s", e.Error())
			}
		}
	}()
	return wsConnection, nil
}

func (w *wsManager) RemoveSubscriber(header *pb.RequestHeader, id int64, closeCode pb.WebsocketCustomCloseCode, closeReason string) error {
	connection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if ok {
		_ = connection.Connection.Close(closeCode, closeReason)
	}
	if header.UserId != "" {
		w.wsConnectionMap.Delete(header.UserId, id)
	}
	go func() {
		if ok {
			_, e := w.svcCtx.CallbackService.UserAfterOffline(context.Background(), &pb.UserAfterOfflineReq{Header: header})
			if e != nil {
				logx.Errorf("UserAfterOffline error: %s", e.Error())
			}
		}
	}()
	return nil
}

// clearConnectionTimer 定时器清除连接
func (w *wsManager) clearConnectionTimer(connection *UniversalConnection) {
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
				w.RemoveSubscriber(connection.GetHeader(), connection.Id, pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout, "heartbeat timeout")
				return
			}
		}
	}
}

func (w *wsManager) KeepAlive(ctx context.Context, connection *UniversalConnection) {
	w.wsConnectionMap.SetAliveTime(ctx, connection.Id, time.Now())
}

func (w *wsManager) WriteData(id int64, writeData *pb.GatewayWriteDataContent) bool {
	wsConnection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if !ok {
		return false
	}
	var data []byte
	if wsConnection.GetHeader().GetEncoding() == pb.EncodingProto_JSON {
		data, _ = json.Marshal(writeData)
	} else {
		data, _ = proto.Marshal(writeData)
	}
	err := wsConnection.Connection.Write(wsConnection.Ctx, data)
	if err != nil {
		logx.Debugf("WriteData error:%v, %s", wsConnection.GetHeader(), err.Error())
		return false
	}
	return true
}

func (w *wsManager) CloseConnection(id int64, code pb.WebsocketCustomCloseCode, reason string) {
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
						w.RemoveSubscriber(connection.GetHeader(), connection.Id, pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout, "heartbeat timeout")
					}
				}
			}
		}
	}()
}
