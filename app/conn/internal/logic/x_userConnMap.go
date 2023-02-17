package logic

import (
	"encoding/json"
	"fmt"
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
	UserId   string
	Platform string
	DeviceId string
	Pointer  string // 指针
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
		conn, ok := l.getConn(cv.Pointer)
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
		UserId:   userId,
		Platform: platform,
		DeviceId: deviceId,
		Pointer:  fmt.Sprintf("%p", conn),
	}
	conn.Pointer = cv.Pointer
	cvBytes, _ := json.Marshal(cv)
	l.connMap.Store(cv.Pointer, conn)
	l.connLength++
	err := l.svcCtx.Redis().Hset(l.redisKey, utils.AnyToString(cv.Pointer), string(cvBytes))
	if err != nil {
		logx.Errorf("UpdateDevice Hset error: %v", err)
	}
	go func() {
		// Expire 7天
		err := l.svcCtx.Redis().Expire(l.redisKey, 24*60*60*7)
		if err != nil {
			logx.Errorf("UpdateDevice Expire error: %v", err)
		}
	}()
	return err
}

// DeleteDevice 删除用户设备的连接
func (l *UserConnStorage) DeleteDevice(userId, platform, deviceId string) error {
	conn, ok := l.LoadDeviceOk(userId, platform, deviceId)
	logx.Debugf("DeleteDevice LoadDeviceOk: %v %v %v %v", userId, platform, deviceId, ok)
	if !ok {
		return nil
	}
	l.connMap.Delete(conn.Pointer)
	l.connLength--
	// hdel
	_, err := l.svcCtx.Redis().Hdel(l.redisKey, utils.AnyToString(conn.Pointer))
	logx.Debugf("DeleteDevice Hdel: %v %v", conn.Pointer, err)
	if err != nil {
		logx.Errorf("DeleteDevice Hdel error: %v", err)
		return err
	}
	return nil
}

func (l *UserConnStorage) Range(f func(id string, conn *types.UserConn) bool) {
	// hgetall
	hgetall, err := l.svcCtx.Redis().Hgetall(l.redisKey)
	if err != nil {
		return
	}
	var userDeviceMap = make(map[string]map[string]*types.UserConn)
	// hgetall -> key: id value: connValue
	for pointer := range hgetall {
		conn, ok := l.getConn(pointer)
		if !ok {
			continue
		}
		if _, ok := userDeviceMap[conn.ConnParam.UserId]; !ok {
			userDeviceMap[conn.ConnParam.UserId] = make(map[string]*types.UserConn)
		}
		if found, ok := userDeviceMap[conn.ConnParam.UserId][conn.ConnParam.DeviceId]; !ok {
			userDeviceMap[conn.ConnParam.UserId][conn.ConnParam.DeviceId] = conn
		} else {
			// 删除旧的连接
			if found.Pointer != conn.Pointer {
				logx.Infof("duplicate connection: %v %v %v %v %v", conn.ConnParam.UserId, conn.ConnParam.Platform, conn.ConnParam.DeviceId, found.ConnectedAt, conn.ConnectedAt)
			}
		}
		if !f(pointer, conn) {
			break
		}
	}
}

func (l *UserConnStorage) getConn(id string) (*types.UserConn, bool) {
	if v, ok := l.connMap.Load(id); ok {
		return v.(*types.UserConn), true
	}
	// redis del
	_, _ = l.svcCtx.Redis().Hdel(l.redisKey, utils.AnyToString(id))
	return nil, false
}
