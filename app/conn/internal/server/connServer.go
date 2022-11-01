package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic"
	"github.com/cherish-chat/xxim-server/app/conn/internal/server/ws"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
)

type ConnServer struct {
	svcCtx *svc.ServiceContext
	Server types.IServer
}

func NewConnServer(svcCtx *svc.ServiceContext) *ConnServer {
	s := &ConnServer{
		svcCtx: svcCtx,
	}
	server := ws.NewServer(svcCtx)
	logic.InitConnLogic(svcCtx)
	l := logic.GetConnLogic()
	server.SetBeforeConnect(l.BeforeConnect)
	server.SetAddSubscriber(l.AddSubscriber)
	server.SetDeleteSubscriber(l.DeleteSubscriber)
	s.Server = server
	go l.Stats()
	return s
}

func (s *ConnServer) Start() {
	s.Server.Start()
}
