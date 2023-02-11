package logic

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

// 使用redis存储
// map[userId]map[platform]map[deviceId]*types.UserConn

// 避免因并发导致panic

// 读取效率优化

// UserConnStorage 用户连接的存储
type UserConnStorage struct {
	svcCtx     *svc.ServiceContext
	connMap    sync.Map // key: id(uint32) value: *types.UserConn
	connLength int
	redisKey   string
	maxId      uint32
}
type connValue struct {
	Id       uint32
	UserId   string
	Platform string
	DeviceId string
}

var singletonUserConnStorage *UserConnStorage

func InitUserConnStorage(svcCtx *svc.ServiceContext) *UserConnStorage {
	singletonUserConnStorage = &UserConnStorage{
		svcCtx:   svcCtx,
		redisKey: rediskey.ConnUserMap(svcCtx.PodIp),
	}
	return singletonUserConnStorage
}

// LoadOk 加载用户的连接
func (l *UserConnStorage) LoadOk(userId string) (map[string]map[string]*types.UserConn, bool) {
	// hgetall
	hgetall, err := l.svcCtx.Redis().Hgetall(l.redisKey)
	if err != nil {
		return nil, false
	}
	var connMap = make(map[string]map[string]*types.UserConn)
	// hgetall -> key: id value: connValue
	for _, v := range hgetall {
		//id := uint32(utils.AnyToInt64(k))
		cv := &connValue{}
		_ = json.Unmarshal([]byte(v), cv)
		if cv.UserId != userId {
			continue
		}
		if _, ok := connMap[cv.Platform]; !ok {
			connMap[cv.Platform] = make(map[string]*types.UserConn)
		}
		conn, ok := l.getConn(cv.Id)
		if !ok {
			continue
		}
		connMap[cv.Platform][cv.DeviceId] = conn
	}
	if len(connMap) == 0 {
		return nil, false
	}
	return connMap, true
}

// LoadDeviceOk 加载用户设备的连接
func (l *UserConnStorage) LoadDeviceOk(userId, platform, deviceId string) (*types.UserConn, bool) {
	userConnMap, ok := l.LoadOk(userId)
	if !ok {
		return nil, false
	}
	if platformConnMap, ok := userConnMap[platform]; ok {
		if conn, ok := platformConnMap[deviceId]; ok {
			return conn, true
		}
	}
	return nil, false
}

// LoadPlatformOk 加载用户平台的连接
func (l *UserConnStorage) LoadPlatformOk(userId, platform string) (map[string]*types.UserConn, bool) {
	userConnMap, ok := l.LoadOk(userId)
	if !ok {
		return nil, false
	}
	if platformConnMap, ok := userConnMap[platform]; ok {
		return platformConnMap, true
	}
	return nil, false
}

// UpdateDevice 更新用户设备的连接
func (l *UserConnStorage) UpdateDevice(userId, platform, deviceId string, conn *types.UserConn) error {
	// hset
	cv := &connValue{
		Id:       l.incId(),
		UserId:   userId,
		Platform: platform,
		DeviceId: deviceId,
	}
	conn.Id = cv.Id
	cvBytes, _ := json.Marshal(cv)
	l.connMap.Store(cv.Id, conn)
	l.connLength++
	err := l.svcCtx.Redis().Hset(l.redisKey, utils.AnyToString(cv.Id), string(cvBytes))
	if err != nil {
		logx.Errorf("UpdateDevice Hset error: %v", err)
	}
	go func() {
		// Expire 2小时
		err := l.svcCtx.Redis().Expire(l.redisKey, 2*60*60)
		if err != nil {
			logx.Errorf("UpdateDevice Expire error: %v", err)
		}
	}()
	return err
}

// DeleteDevice 删除用户设备的连接
func (l *UserConnStorage) DeleteDevice(userId, platform, deviceId string) error {
	conn, ok := l.LoadDeviceOk(userId, platform, deviceId)
	if !ok {
		return nil
	}
	l.connMap.Delete(conn.Id)
	l.connLength--
	// hdel
	_, err := l.svcCtx.Redis().Hdel(l.redisKey, utils.AnyToString(conn.Id))
	if err != nil {
		logx.Errorf("DeleteDevice Hdel error: %v", err)
		return err
	}
	return nil
}

func (l *UserConnStorage) Range(f func(id uint32, conn *types.UserConn) bool) {
	// hgetall
	hgetall, err := l.svcCtx.Redis().Hgetall(l.redisKey)
	if err != nil {
		return
	}
	// hgetall -> key: id value: connValue
	for k, v := range hgetall {
		id := uint32(utils.AnyToInt64(k))
		cv := &connValue{}
		_ = json.Unmarshal([]byte(v), cv)
		conn, ok := l.getConn(id)
		if !ok {
			continue
		}
		if !f(id, conn) {
			break
		}
	}
}

func (l *UserConnStorage) getConn(id uint32) (*types.UserConn, bool) {
	if v, ok := l.connMap.Load(id); ok {
		return v.(*types.UserConn), true
	}
	// redis del
	_, _ = l.svcCtx.Redis().Hdel(l.redisKey, utils.AnyToString(id))
	return nil, false
}

func (l *UserConnStorage) incId() uint32 {
	l.maxId++
	return l.maxId
}
