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

type ConnLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	// 未认证的连接
	unknownConnMap sync.Map // key: *types.UserConn
}

var singletonConnLogic *ConnLogic

func InitConnLogic(svcCtx *svc.ServiceContext) *ConnLogic {
	InitUserConnStorage(svcCtx)
	if singletonConnLogic == nil {
		l := &ConnLogic{svcCtx: svcCtx}
		l.Logger = logx.WithContext(context.Background())
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
	l.Debugf("user %s connected", utils.AnyToString(param))
	// 是否未认证的连接
	if param.UserId == "" || param.Token == "" {
		l.unknownConnMap.Store(c, struct{}{})
		// 告知客户端连接成功
		_ = c.Conn.Write(context.Background(), int(websocket.MessageText), []byte("connected"))
		return
	}
	// 删除未认证的连接
	l.unknownConnMap.Delete(c)
	// 加入用户连接
	{
		if userConn, ok := singletonUserConnStorage.LoadDeviceOk(param.UserId, param.Platform, param.DeviceId); ok {
			// 是不是同一条连接
			if userConn.Conn == c.Conn {
				// 是同一条连接
				return
			}
			userConn.Conn.Close(int(websocket.StatusNormalClosure), "duplicate connection")
		}
		err := singletonUserConnStorage.UpdateDevice(param.UserId, param.Platform, param.DeviceId, c)
		if err != nil {
			return
		}
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
	l.Debugf("user %s disconnected", utils.AnyToString(c.ConnParam))
	// 是否未认证的连接
	if _, ok := l.unknownConnMap.Load(c); ok {
		l.unknownConnMap.Delete(c)
		return
	}
	// 删除用户连接
	{
		if _, ok := singletonUserConnStorage.LoadDeviceOk(c.ConnParam.UserId, c.ConnParam.Platform, c.ConnParam.DeviceId); !ok {
			return
		}
		singletonUserConnStorage.DeleteDevice(c.ConnParam.UserId, c.ConnParam.Platform, c.ConnParam.DeviceId)
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
	var onlineUserIds []string
	var onlineDeviceIds []string
	singletonUserConnStorage.Range(func(id uint32, conn *types.UserConn) bool {
		onlineUserIds = append(onlineUserIds, conn.ConnParam.UserId)
		onlineDeviceIds = append(onlineDeviceIds, conn.ConnParam.DeviceId)
		return true
	})
	onlineUserCount = len(utils.Set(onlineUserIds))
	onlineDeviceCount = len(utils.Set(onlineDeviceIds))
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
		singletonUserConnStorage.Range(func(id uint32, conn *types.UserConn) bool {
			if filter(conn) {
				conns = append(conns, conn)
			}
			return true
		})
	}
	return conns
}

func (l *ConnLogic) SendMsgToConn(c *types.UserConn, data []byte) error {
	return c.Conn.Write(c.Ctx, int(websocket.MessageBinary), data)
}
