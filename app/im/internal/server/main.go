package server

import (
	"github.com/cherish-chat/xxim-server/app/im/internal/ws"
)

func (s *ImServiceServer) Start() {
	s.http()
	ws.NewConsumerLogic(s.svcCtx).Start()
}
