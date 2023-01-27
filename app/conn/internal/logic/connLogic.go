package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xerr"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type ConnMap map[string]DeviceMap // key: platform, value: deviceMap

type DeviceMap map[string]*types.UserConn // key: deviceId, value: *types.UserConn

type ConnLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	userConnMap sync.Map // key: userId value: ConnMap
	// 未认证的连接
	unknownConnMap sync.Map // key: *types.UserConn
}

func (l *ConnLogic) LoadOk(userId string) (ConnMap, bool) {
	if v, ok := l.userConnMap.Load(userId); ok {
		return v.(ConnMap), true
	}
	return nil, false
}

func (l *ConnLogic) Load(userId string) ConnMap {
	if v, ok := l.userConnMap.Load(userId); ok {
		return v.(ConnMap)
	}
	cm := ConnMap{}
	l.Store(userId, cm)
	return cm
}

func (l *ConnLogic) Store(userId string, cm ConnMap) {
	l.userConnMap.Store(userId, cm)
}

func (l *ConnLogic) LoadPlatform(userId, platform string) DeviceMap {
	connMap := l.Load(userId)
	deviceMap, ok := connMap[platform]
	if !ok {
		deviceMap = DeviceMap{}
		l.UpdatePlatform(userId, platform, deviceMap)
	}
	return deviceMap
}

func (l *ConnLogic) LoadPlatformOk(userId, platform string) (DeviceMap, bool) {
	connMap := l.Load(userId)
	dm, ok := connMap[platform]
	return dm, ok
}

func (l *ConnLogic) UpdatePlatform(userId string, platform string, dm DeviceMap) {
	connMap := l.Load(userId)
	connMap[platform] = dm
	l.Store(userId, connMap)
}

func (l *ConnLogic) LoadDeviceOk(userId, platform, deviceId string) (*types.UserConn, bool) {
	dm := l.LoadPlatform(userId, platform)
	if conn, ok := dm[deviceId]; ok {
		return conn, true
	}
	return nil, false
}

func (l *ConnLogic) UpdateDevice(userId, platform, deviceId string, userConn *types.UserConn) {
	dm := l.LoadPlatform(userId, platform)
	dm[deviceId] = userConn
	l.UpdatePlatform(userId, platform, dm)
}

func (l *ConnLogic) DeleteDevice(userId string, platform string, deviceId string) {
	dm := l.LoadPlatform(userId, platform)
	delete(dm, deviceId)
	l.UpdatePlatform(userId, platform, dm)
}

func (l *ConnLogic) CheckNoUser(userId string) {
	if connMap, ok := l.LoadOk(userId); !ok {
		// 不在线
		return
	} else if len(connMap) == 0 {
		// 不在线
		l.userConnMap.Delete(userId)
	} else {
		// 遍历所有平台
		for platform, dm := range connMap {
			if len(dm) == 0 {
				// 不在线
				delete(connMap, platform)
			}
		}
		if len(connMap) == 0 {
			l.userConnMap.Delete(userId)
		} else {
			l.Store(userId, connMap)
		}
	}
}

var singletonConnLogic *ConnLogic

func InitConnLogic(svcCtx *svc.ServiceContext) *ConnLogic {
	if singletonConnLogic == nil {
		l := &ConnLogic{svcCtx: svcCtx}
		l.Logger = logx.WithContext(context.Background())
		l.userConnMap = sync.Map{}
		singletonConnLogic = l
	}
	return singletonConnLogic
}

func GetConnLogic() *ConnLogic {
	return singletonConnLogic
}

func (l *ConnLogic) BeforeConnect(ctx context.Context, param types.ConnParam) (int, error) {
	if param.UserId == "" || param.Token == "" {
		// 通过
		return 0, nil
	}
	resp, err := l.svcCtx.ImService().BeforeConnect(ctx, &pb.BeforeConnectReq{
		ConnParam: &pb.ConnParam{
			UserId:      param.UserId,
			Token:       param.Token,
			DeviceId:    param.DeviceId,
			Platform:    param.Platform,
			Ips:         param.Ips,
			NetworkUsed: param.NetworkUsed,
			Headers:     param.Headers,
		},
	})
	if err != nil {
		return 0, err
	}
	if resp.Msg != "" {
		return int(resp.Code), errors.New(resp.Msg)
	}
	return int(resp.Code), nil
}

func (l *ConnLogic) AddSubscriber(c *types.UserConn) {
	param := c.ConnParam
	l.Infof("user %s connected", utils.AnyToString(param))
	// 告知客户端连接成功
	_ = c.Conn.Write(context.Background(), int(websocket.MessageText), []byte("connected"))
	// 是否未认证的连接
	if param.UserId == "" || param.Token == "" {
		l.unknownConnMap.Store(c, struct{}{})
		return
	}
	// 删除未认证的连接
	l.unknownConnMap.Delete(c)
	// 加入用户连接
	{
		if userConn, ok := l.LoadDeviceOk(param.UserId, param.Platform, param.DeviceId); ok {
			userConn.Conn.Close(int(websocket.StatusNormalClosure), "duplicate connection")
		}
		l.UpdateDevice(param.UserId, param.Platform, param.DeviceId, c)
	}
	go func() {
		for {
			_, err := l.svcCtx.ImService().AfterConnect(c.Ctx, &pb.AfterConnectReq{
				ConnParam: &pb.ConnParam{
					UserId:      param.UserId,
					Token:       param.Token,
					DeviceId:    param.DeviceId,
					Platform:    param.Platform,
					Ips:         param.Ips,
					NetworkUsed: param.NetworkUsed,
					Headers:     param.Headers,
					PodIp:       l.svcCtx.PodIp,
				},
				ConnectedAt: utils.AnyToString(c.ConnectedAt.UnixMilli()),
			})
			if err != nil {
				// 是否是 context canceled
				if xerr.IsCanceled(err) {
					break
				}
				l.Errorf("AfterConnect error: %s", err.Error())
				utils.ProxySleep(c.Ctx)
			} else {
				break
			}
		}
	}()
	l.stats()
}

func (l *ConnLogic) DeleteSubscriber(c *types.UserConn) {
	l.Infof("user %s disconnected", utils.AnyToString(c.ConnParam))
	// 是否未认证的连接
	if _, ok := l.unknownConnMap.Load(c); ok {
		l.unknownConnMap.Delete(c)
		return
	}
	// 删除用户连接
	{
		if _, ok := l.LoadOk(c.ConnParam.UserId); !ok {
			return
		}
		if _, ok := l.LoadPlatformOk(c.ConnParam.UserId, c.ConnParam.Platform); !ok {
		}
		if _, ok := l.LoadDeviceOk(c.ConnParam.UserId, c.ConnParam.Platform, c.ConnParam.DeviceId); !ok {
			return
		}
		l.DeleteDevice(c.ConnParam.UserId, c.ConnParam.Platform, c.ConnParam.DeviceId)
		l.CheckNoUser(c.ConnParam.UserId)
	}
	l.stats()
	go func() {
		for {
			_, err := l.svcCtx.ImService().AfterDisconnect(context.Background(), &pb.AfterDisconnectReq{
				ConnParam: &pb.ConnParam{
					UserId:      c.ConnParam.UserId,
					Token:       c.ConnParam.Token,
					DeviceId:    c.ConnParam.DeviceId,
					Platform:    c.ConnParam.Platform,
					Ips:         c.ConnParam.Ips,
					NetworkUsed: c.ConnParam.NetworkUsed,
					Headers:     c.ConnParam.Headers,
					PodIp:       l.svcCtx.PodIp,
				},
				ConnectedAt:    utils.AnyToString(c.ConnectedAt.UnixMilli()),
				DisconnectedAt: utils.AnyToString(time.Now().UnixMilli()),
			})
			if err != nil {
				// 是否是 context canceled
				if xerr.IsCanceled(err) {
					break
				}
				l.Errorf("AfterDisconnect error: %s", err.Error())
				utils.ProxySleep(c.Ctx)
			} else {
				break
			}
		}
	}()
}

func (l *ConnLogic) Stats() {
	l.stats()
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			l.stats()
		}
	}
}

func (l *ConnLogic) stats() {
	// 统计在线用户数
	onlineUserCount := 0
	onlineDeviceCount := 0
	l.userConnMap.Range(func(key, value any) bool {
		onlineUserCount++
		connMap := value.(ConnMap)
		for _, deviceMap := range connMap {
			onlineDeviceCount += len(deviceMap)
		}
		return true
	})
	l.Infof("online user count: %d, online device count: %d", onlineUserCount, onlineDeviceCount)
}

func (l *ConnLogic) BuildSearchUserConnFilter(in *pb.GetUserConnReq) func(conn *types.UserConn) bool {
	return func(conn *types.UserConn) bool {
		if len(in.UserIds) > 0 {
			in := utils.InSlice(in.UserIds, conn.ConnParam.UserId)
			if !in {
				return false
			}
		}
		if len(in.Platforms) > 0 {
			in := utils.InSlice(in.Platforms, conn.ConnParam.Platform)
			if !in {
				return false
			}
		}
		if len(in.Devices) > 0 {
			in := utils.InSlice(in.Devices, conn.ConnParam.DeviceId)
			if !in {
				return false
			}
		}
		return true
	}
}

func (l *ConnLogic) KickUserConn(ctx context.Context, in *pb.KickUserConnReq) error {
	var conns []*types.UserConn
	xtrace.StartFuncSpan(ctx, "GetConnsByFilter", func(ctx context.Context) {
		conns = l.GetConnsByFilter(l.BuildSearchUserConnFilter(in.GetUserConnReq))
	})
	for _, c := range conns {
		var err error
		xtrace.StartFuncSpan(ctx, "KickUserConn", func(ctx context.Context) {
			err = l.KickConn(c)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *ConnLogic) KickConn(c *types.UserConn) error {
	err := c.Conn.Close(1001, "kick")
	// 如果是 context deadline exceeded，说明连接已经断开了
	if xerr.IsCanceled(err) {
		return nil
	}
	return err
}

func (l *ConnLogic) SendMsg(in *pb.SendMsgReq) error {
	var err error
	{
		conns := l.GetConnsByFilter(l.BuildSearchUserConnFilter(in.GetUserConnReq))
		data, _ := proto.Marshal(&pb.PushBody{
			Event: in.Event,
			Data:  in.Data,
		})
		for _, c := range conns {
			err = l.SendMsgToConn(c, data)
			if err != nil {
				break
			}
		}
	}
	return err
}

func (l *ConnLogic) GetConnsByFilter(filter func(c *types.UserConn) bool) []*types.UserConn {
	conns := make([]*types.UserConn, 0)
	{
		l.userConnMap.Range(func(key, value any) bool {
			//userId := key.(string)
			cm := value.(ConnMap)
			for _, dm := range cm {
				for _, c := range dm {
					if filter(c) {
						conns = append(conns, c)
					}
				}
			}
			return true
		})
	}
	return conns
}

func (l *ConnLogic) SendMsgToConn(c *types.UserConn, data []byte) error {
	return c.Conn.Write(c.Ctx, int(websocket.MessageBinary), data)
}
