package server

import "github.com/cherish-chat/xxim-server/app/msg/internal/logic"

func (s *MsgServiceServer) Start() {
	l := logic.NewConsumerLogic(s.svcCtx)
	go l.Start()
}
