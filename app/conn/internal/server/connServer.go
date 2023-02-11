package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic"
	"github.com/cherish-chat/xxim-server/app/conn/internal/server/tcp"
	"github.com/cherish-chat/xxim-server/app/conn/internal/server/ws"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
)

type ConnServer struct {
	svcCtx  *svc.ServiceContext
	Servers []types.IServer
}

func NewConnServer(svcCtx *svc.ServiceContext) *ConnServer {
	s := &ConnServer{
		svcCtx: svcCtx,
	}
	logic.InitConnLogic(svcCtx)
	l := logic.GetConnLogic()

	var servers = []types.IServer{ws.NewServer(svcCtx), tcp.NewServer(svcCtx)}
	for _, server := range servers {
		server.SetBeforeConnect(l.BeforeConnect)
		server.SetAddSubscriber(l.AddSubscriber)
		server.SetDeleteSubscriber(l.DeleteSubscriber)
		server.SetOnReceive(l.OnReceive)
	}
	s.Servers = servers
	s.registerGateway()
	go l.Stats()
	return s
}

func (s *ConnServer) Start() {
	for _, server := range s.Servers {
		go server.Start()
	}
}
