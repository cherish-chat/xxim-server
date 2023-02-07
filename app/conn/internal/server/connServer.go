package server

import (
	"context"
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
	server.SetOnReceive(l.OnReceive)
	s.Server = server
	s.registerGateway()
	go logic.NewKeepAliveLogic(context.Background(), svcCtx).Start()
	go l.Stats()
	return s
}

func (s *ConnServer) Start() {
	s.Server.Start()
}
