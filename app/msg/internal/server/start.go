package server

import "github.com/cherish-chat/xxim-server/app/msg/internal/logic"

func (s *MsgServiceServer) Start() {
	if s.svcCtx.Config.TDMQ.Enabled {
		{
			l := logic.NewConsumerLogic(s.svcCtx)
			go l.Start()
		}
	}
	{
		l := logic.NewTimerCleanSubscriptionLogic(s.svcCtx)
		go l.Start()
	}
}
