package gatewayservicelogic

import (
	"context"
	"crypto"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"sync"
	"time"
)

type xConnectionLogic struct {
	svcCtx *svc.ServiceContext

	// 连接暂存区
	connections []*Connection

	// 连接暂存区锁
	connectionsLock sync.RWMutex

	// 通过`传输协议验证`后的连接池
	verifiedConnections []*Connection

	// 通过`传输协议验证`后的连接池锁
	verifiedConnectionsLock sync.RWMutex

	// 设置好userId后的连接池 Map
	userConnectionsMap map[string]*installConnectionMap // key: userId value: *installConnectionMap

	// 设置好userId后的连接池 Map锁
	userConnectionsMapLock sync.RWMutex
}

type installConnectionMap struct {
	m map[string]*Connection // key: installId value: *Connection
	l sync.RWMutex
}

type Connection struct {
	ctx context.Context
	//请求头
	header *pb.RequestHeader
	//设置请求头的锁
	headerLock sync.RWMutex
	//服务端生成的公私钥
	ServerPrivateKey crypto.PrivateKey
	ServerPublicKey  crypto.PublicKey

	//客户端生成的公钥
	ClientPublicKey crypto.PublicKey
	//经过一致性算法后的共享密钥
	SharedSecret []byte
	//设置公钥的锁
	PublicKeyLock sync.RWMutex

	//连接
	Connection types.LongConnection
	//连接时间
	ConnectedTime time.Time
}

func (c *Connection) SendMessage(ctx context.Context, body []byte) error {
	// 是否加密 取决于 ClientPublicKey 是否为空
	var aesKey []byte
	var aesIv []byte
	var isEncrypt bool

	c.PublicKeyLock.RLock()
	{
		if len(c.SharedSecret) == 0 {
			// 不加密
			isEncrypt = false
		} else {
			// 加密
			isEncrypt = true
			aesKey = c.SharedSecret[:]
			aesIv = c.SharedSecret[8:24]
		}
	}
	c.PublicKeyLock.RUnlock()

	if isEncrypt {
		// 加密
		body = utils.Aes.Encrypt(aesKey, aesIv, body)
	}
	return c.Connection.SendMessage(ctx, body)
}

func (c *Connection) GetHeader() *pb.RequestHeader {
	c.headerLock.RLock()
	defer c.headerLock.RUnlock()
	return c.header
}

func (c *Connection) ToPb() *pb.LongConnection {
	return &pb.LongConnection{
		Header: c.GetHeader(),
	}
}

var ConnectionLogic *xConnectionLogic

func InitConnectionLogic(svcCtx *svc.ServiceContext) {
	ConnectionLogic = &xConnectionLogic{
		svcCtx:                  svcCtx,
		connections:             make([]*Connection, 0),
		connectionsLock:         sync.RWMutex{},
		verifiedConnections:     make([]*Connection, 0),
		verifiedConnectionsLock: sync.RWMutex{},
		userConnectionsMap:      make(map[string]*installConnectionMap),
		userConnectionsMapLock:  sync.RWMutex{},
	}
	go ConnectionLogic.loopCheck()
}

func (l *xConnectionLogic) GetConnectionsByUserIds(ids []string) []*Connection {
	var connections []*Connection
	l.userConnectionsMapLock.RLock()
	for _, id := range ids {
		userConnectionMap, ok := l.userConnectionsMap[id]
		if !ok {
			continue
		}
		userConnectionMap.l.RLock()
		for _, connection := range userConnectionMap.m {
			connections = append(connections, connection)
		}
		userConnectionMap.l.RUnlock()
	}
	l.userConnectionsMapLock.RUnlock()
	return connections
}

func (l *xConnectionLogic) GetAllConnections() []*Connection {
	var connections []*Connection
	l.userConnectionsMapLock.RLock()
	for _, userConnectionMap := range l.userConnectionsMap {
		userConnectionMap.l.RLock()
		for _, connection := range userConnectionMap.m {
			connections = append(connections, connection)
		}
		userConnectionMap.l.RUnlock()
	}
	l.userConnectionsMapLock.RUnlock()
	return connections
}

func (l *xConnectionLogic) OnConnect(conn *Connection) {
	l.connectionsLock.Lock()
	l.connections = append(l.connections, conn)
	l.connectionsLock.Unlock()
}

func (l *xConnectionLogic) OnVerified(conn *Connection) {
	l.verifiedConnectionsLock.Lock()
	l.verifiedConnections = append(l.verifiedConnections, conn)
	l.verifiedConnectionsLock.Unlock()
	// 删除未验证的连接
	l.connectionsLock.Lock()
	for i, connection := range l.connections {
		if connection == conn {
			l.connections = append(l.connections[:i], l.connections[i+1:]...)
			break
		}
	}
	l.connectionsLock.Unlock()
}

func (l *xConnectionLogic) OnLogin(conn *Connection) {
	l.userConnectionsMapLock.Lock()
	header := conn.GetHeader()
	userConnectionMap, ok := l.userConnectionsMap[header.UserId]
	if !ok {
		userConnectionMap = &installConnectionMap{
			m: make(map[string]*Connection),
			l: sync.RWMutex{},
		}
		l.userConnectionsMap[header.UserId] = userConnectionMap
	} else {
		userConnectionMap.l.Lock()
		old, exist := userConnectionMap.m[header.InstallId]
		if exist {
			old.Connection.CloseConnection(pb.WebsocketCustomCloseCode_CloseCodeDuplicateConnection, "duplicate Connection")
		}
		userConnectionMap.m[header.InstallId] = conn
		userConnectionMap.l.Unlock()
	}
	l.userConnectionsMapLock.Unlock()
	l.verifiedConnectionsLock.Lock()
	for i, connection := range l.verifiedConnections {
		if connection == conn {
			l.verifiedConnections = append(l.verifiedConnections[:i], l.verifiedConnections[i+1:]...)
			break
		}
	}
	l.verifiedConnectionsLock.Unlock()
}

func (l *xConnectionLogic) OnDisconnect(conn *Connection) {
	l.connectionsLock.Lock()
	for i, connection := range l.connections {
		if connection == conn {
			l.connections = append(l.connections[:i], l.connections[i+1:]...)
			break
		}
	}
	l.connectionsLock.Unlock()
	l.verifiedConnectionsLock.Lock()
	for i, connection := range l.verifiedConnections {
		if connection == conn {
			l.verifiedConnections = append(l.verifiedConnections[:i], l.verifiedConnections[i+1:]...)
			break
		}
	}
	l.verifiedConnectionsLock.Unlock()
	l.userConnectionsMapLock.Lock()
	for userId, userConnectionMap := range l.userConnectionsMap {
		userConnectionMap.l.Lock()
		for installId, connection := range userConnectionMap.m {
			if connection == conn {
				delete(userConnectionMap.m, installId)
				break
			}
		}
		userConnectionMap.l.Unlock()
		if len(userConnectionMap.m) == 0 {
			delete(l.userConnectionsMap, userId)
		}
	}
	l.userConnectionsMapLock.Unlock()
}

func (l *xConnectionLogic) loopCheck() {
	// 定期检测
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			//查询所有连接
			var connections = make(map[*Connection]struct{})
			l.connectionsLock.RLock()
			for _, conn := range l.connections {
				connections[conn] = struct{}{}
			}
			l.connectionsLock.RUnlock()
			l.verifiedConnectionsLock.RLock()
			for _, conn := range l.verifiedConnections {
				connections[conn] = struct{}{}
			}
			l.verifiedConnectionsLock.RUnlock()
			for connection := range connections {
				// 判断是否应该清除
				// 连接后超过n分钟未验证 就应该清除
				if time.Now().Sub(connection.ConnectedTime) > time.Minute*5 {
					connection.Connection.CloseConnection(pb.WebsocketCustomCloseCode_CloseCodeNoStatusReceived, "Connection timeout")
				}
			}
		}
	}
}
