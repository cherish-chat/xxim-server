package gatewayservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayKeepAliveLogic {
	l := &GatewayKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
	if userKeepAliveMap == nil {
		userKeepAliveMap = new(sync.Map)
		go l.checkTimer()
	}
	return l
}

func (l *GatewayKeepAliveLogic) GatewayKeepAlive(in *pb.GatewayKeepAliveReq) (*pb.GatewayKeepAliveResp, error) {
	return &pb.GatewayKeepAliveResp{}, status.Error(codes.Unimplemented, "cannot call method GatewayKeepAlive")
}

var userKeepAliveMap *sync.Map // key: userId, value: time.Time

// KeepAlive 保持连接
// 客户端必须每隔 config.Websocket.KeepAliveSecond 秒发送一次心跳包
// 二次开发人员可以在这里修改逻辑，比如一致性算法安全校验等
func (l *GatewayKeepAliveLogic) KeepAlive(connection *UniversalConnection, c *pb.GatewayApiRequest) (*pb.GatewayApiResponse, error) {
	WsManager.KeepAlive(l.ctx, connection)
	_, err := l.svcCtx.CallbackService.UserAfterKeepAlive(l.ctx, &pb.UserAfterKeepAliveReq{
		Header: c.Header,
	})
	if err != nil {
		l.Errorf("UserAfterKeepAlive error: %v", err)
	}
	header := connection.GetHeader()
	if header.UserId != "" {
		if _, ok := userKeepAliveMap.Load(header.UserId); !ok {
			// 用户上线
			_, _ = l.svcCtx.CallbackService.UserAfterOnline(l.ctx, &pb.UserAfterOnlineReq{Header: header})
		}
	}
	userKeepAliveMap.Store(header.UserId, time.Now())
	return &pb.GatewayApiResponse{
		Header:    i18n.NewOkHeader(),
		RequestId: c.RequestId,
		Path:      c.Path,
		Body:      nil,
	}, nil
}

func (l *GatewayKeepAliveLogic) checkTimer() {
	ticker := time.NewTicker(time.Second * time.Duration(l.svcCtx.Config.Websocket.OfflineDeterminationSecond))
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			userKeepAliveMap.Range(func(key, value interface{}) bool {
				userId := key.(string)
				lastTime := value.(time.Time)
				if now.Sub(lastTime).Seconds() > float64(l.svcCtx.Config.Websocket.KeepAliveSecond) {
					// 用户下线
					_, _ = l.svcCtx.CallbackService.UserAfterOffline(l.ctx, &pb.UserAfterOfflineReq{UserId: userId})
					userKeepAliveMap.Delete(userId)
				}
				return true
			})
		}
	}
}
