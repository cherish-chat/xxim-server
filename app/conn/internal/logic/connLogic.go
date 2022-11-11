package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type connMap map[string]deviceMap // key: platform, value: deviceMap

type deviceMap map[string]*types.UserConn // key: deviceId, value: *types.UserConn

type ConnLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	userConnMapLock sync.RWMutex
	userConnMap     map[string]connMap // key: userId, value: connMap
}

var singletonConnLogic *ConnLogic

func InitConnLogic(svcCtx *svc.ServiceContext) *ConnLogic {
	if singletonConnLogic == nil {
		l := &ConnLogic{svcCtx: svcCtx}
		l.Logger = logx.WithContext(context.Background())
		l.userConnMap = map[string]connMap{}
		singletonConnLogic = l
	}
	return singletonConnLogic
}

func GetConnLogic() *ConnLogic {
	return singletonConnLogic
}

func (l *ConnLogic) BeforeConnect(ctx context.Context, param types.ConnParam) (int, error) {
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
	// 加入用户连接
	l.userConnMapLock.Lock()
	if _, ok := l.userConnMap[param.UserId]; !ok {
		l.userConnMap[param.UserId] = connMap{}
	}
	if _, ok := l.userConnMap[param.UserId][param.Platform]; !ok {
		l.userConnMap[param.UserId][param.Platform] = deviceMap{}
	}
	if _, ok := l.userConnMap[param.UserId][param.Platform][param.DeviceId]; ok {
		l.userConnMap[param.UserId][param.Platform][param.DeviceId].Conn.Close(int(websocket.StatusNormalClosure), "duplicate connection")
	}
	l.userConnMap[param.UserId][param.Platform][param.DeviceId] = c
	l.userConnMapLock.Unlock()
	// 告知客户端连接成功
	_ = c.Conn.Write(context.Background(), int(websocket.MessageText), []byte("connected"))
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

func (l *ConnLogic) KickUserConn(in *pb.KickUserConnReq) error {
	conns := l.GetConnsByFilter(l.BuildSearchUserConnFilter(in.GetUserConnReq))
	for _, c := range conns {
		err := l.KickConn(c)
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
	conns := l.GetConnsByFilter(l.BuildSearchUserConnFilter(in.GetUserConnReq))
	data, _ := proto.Marshal(&pb.PushBody{
		Event: in.Event,
		Data:  in.Data,
	})
	for _, c := range conns {
		err := l.SendMsgToConn(c, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *ConnLogic) GetConnsByFilter(filter func(c *types.UserConn) bool) []*types.UserConn {
	l.userConnMapLock.RLock()
	defer l.userConnMapLock.RUnlock()
	conns := make([]*types.UserConn, 0)
	for _, cm := range l.userConnMap {
		for _, dm := range cm {
			for _, c := range dm {
				if filter(c) {
					conns = append(conns, c)
				}
			}
		}
	}
	return conns
}

func (l *ConnLogic) SendMsgToConn(c *types.UserConn, data []byte) error {
	return c.Conn.Write(c.Ctx, int(websocket.MessageBinary), data)
}
