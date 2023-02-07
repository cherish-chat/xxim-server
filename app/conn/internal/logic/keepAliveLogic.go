package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/zeromicro/go-zero/core/mr"
	"runtime"
	"time"

	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeepAliveLogic {
	return &KeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KeepAliveLogic) KeepAlive(in *pb.KeepAliveReq) (*pb.KeepAliveResp, error) {
	return l.svcCtx.ImService().KeepAlive(l.ctx, in)
}

func (l *KeepAliveLogic) Start() {
	logic := GetConnLogic()
	ticker := time.NewTicker(time.Second * time.Duration(l.svcCtx.Config.Websocket.KaInterval))
	for {
		select {
		case <-ticker.C:
			// get all subscribers
			userConns := logic.GetConnsByFilter(func(c *types.UserConn) bool {
				return true
			})
			// number of cpu
			cpu := runtime.NumCPU()
			var fs []func()
			// 把所有的用户连接分成cpu个数的组，每个组一个goroutine去处理
			for i := 0; i < cpu; i++ {
				fs = append(fs, func() {
					// 根据 i 的值，计算出每个 goroutine 要处理的用户连接
					for j := i; j < len(userConns); j += cpu {
						conn := userConns[j]
						// send keep alive
						_, err := l.KeepAlive(&pb.KeepAliveReq{CommonReq: &pb.CommonReq{
							UserId:      conn.ConnParam.UserId,
							Token:       conn.ConnParam.Token,
							DeviceModel: conn.ConnParam.DeviceModel,
							DeviceId:    conn.ConnParam.DeviceId,
							OsVersion:   conn.ConnParam.OsVersion,
							Platform:    conn.ConnParam.Platform,
							AppVersion:  conn.ConnParam.AppVersion,
							Language:    conn.ConnParam.Language,
							Ip:          conn.ConnParam.Ips,
						}})
						if err != nil {
							l.Errorf("keep alive error: %s", err)
						}
					}
				})
			}
			// 并发执行
			if len(fs) > 0 {
				mr.FinishVoid(fs...)
			}
		case <-l.ctx.Done():
			return
		}
	}
}
