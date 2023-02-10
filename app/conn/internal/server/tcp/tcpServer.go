package tcp

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/pkg/zinx/ziface"
	"github.com/cherish-chat/xxim-server/app/conn/internal/pkg/zinx/znet"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"net"
	"strconv"
	"sync"
	"time"
)

var connIdMap sync.Map // key: uint32 value: *types.UserConn

type userConn struct {
	ctx        context.Context
	tcp        *net.TCPConn
	dataPacker ziface.Packet
}

func (c *userConn) Close(code int, desc string) error {
	return c.tcp.Close()
}

func (c *userConn) Write(ctx context.Context, typ int, msg []byte) error {
	msgpkg := znet.NewMsgPackage(0, msg)
	buf, err := c.dataPacker.Pack(msgpkg)
	if err != nil {
		return err
	}
	_, err = c.tcp.Write(buf)
	return err
}

func (c *userConn) Read(ctx context.Context) (typ int, msg []byte, err error) {
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}

// Server enables broadcasting to a set of subscribers.
type Server struct {
	svcCtx           *svc.ServiceContext
	addSubscriber    func(c *types.UserConn)
	deleteSubscriber func(c *types.UserConn)
	beforeConnect    func(ctx context.Context, param types.ConnParam) (int, error)
	onReceive        func(ctx context.Context, c *types.UserConn, typ int, msg []byte)
	zinx             ziface.IServer
}

func (s *Server) SetOnReceive(f func(ctx context.Context, c *types.UserConn, typ int, msg []byte)) {
	s.onReceive = f
}

func (s *Server) SetBeforeConnect(f func(ctx context.Context, param types.ConnParam) (int, error)) {
	s.beforeConnect = f
}

func (s *Server) SetAddSubscriber(f func(c *types.UserConn)) {
	s.addSubscriber = f
}

func (s *Server) SetDeleteSubscriber(f func(c *types.UserConn)) {
	s.deleteSubscriber = f
}

type zinxHandler struct {
	znet.BaseRouter
	server *Server
}

func (l *zinxHandler) Handle(request ziface.IRequest) {
	uc := l.server.iConnection2UserConn(request.GetConnection())
	if uc == nil {
		return
	}
	msg := request.GetData()
	go xtrace.RunWithTrace("", "ReadFromConn", func(ctx context.Context) {
		l.server.onReceive(uc.Ctx, uc, 2, msg)
	}, propagation.MapCarrier{
		"length":      strconv.Itoa(len(msg)),
		"userId":      uc.ConnParam.UserId,
		"platform":    uc.ConnParam.Platform,
		"deviceId":    uc.ConnParam.DeviceId,
		"ips":         uc.ConnParam.Ips,
		"networkUsed": uc.ConnParam.NetworkUsed,
	})
}

// NewServer constructs a Server with the defaults.
func NewServer(
	svcCtx *svc.ServiceContext,
) types.IServer {
	s := &Server{
		svcCtx:           svcCtx,
		addSubscriber:    func(c *types.UserConn) {},
		deleteSubscriber: func(c *types.UserConn) {},
		beforeConnect:    func(ctx context.Context, param types.ConnParam) (int, error) { return 0, nil },
		zinx:             znet.NewServer(),
	}
	s.zinx.SetOnConnStop(s.onConnStop)
	s.zinx.SetOnConnStart(s.onConnStart)
	handler := &zinxHandler{server: s}
	s.zinx.AddRouter(0, handler)
	return s
}

func (s *Server) Start() error {
	s.zinx.Start()
	return nil
}

func (s *Server) onConnStop(iConnection ziface.IConnection) {
	if iConnection == nil {
		return
	}
	uc := s.iConnection2UserConn(iConnection)
	if uc == nil {
		return
	}
	s.deleteSubscriber(uc)
	s.deleteConn(iConnection)
}

func (s *Server) iConnection2UserConn(iConnection ziface.IConnection) *types.UserConn {
	conn, ok := connIdMap.Load(iConnection.GetConnID())
	if !ok {
		iConnection.GetTCPConnection().Close()
		return nil
	}
	return conn.(*types.UserConn)
}

func (s *Server) onConnStart(iConnection ziface.IConnection) {
	if iConnection == nil {
		return
	}
	uc := &userConn{
		ctx:        iConnection.Context(),
		tcp:        iConnection.GetTCPConnection(),
		dataPacker: znet.NewDataPack(),
	}
	now := time.Now()
	typeConn := &types.UserConn{
		Conn: uc,
		ConnParam: types.ConnParam{
			UserId:      "",
			Token:       "",
			DeviceId:    "",
			DeviceModel: "",
			OsVersion:   "",
			AppVersion:  "",
			Language:    "",
			Platform:    "",
			Ips:         iConnection.RemoteAddr().String(),
			NetworkUsed: "",
			Headers:     nil,
			Timestamp:   now.UnixMilli(),
		},
		Ctx:         iConnection.Context(),
		ConnectedAt: now,
	}
	connIdMap.Store(iConnection.GetConnID(), typeConn)
	s.addSubscriber(typeConn)
}

func (s *Server) deleteConn(iConnection ziface.IConnection) {
	connIdMap.Delete(iConnection.GetConnID())
}
