package server

import "github.com/cherish-chat/xxim-server/app/notice/internal/logic"

func (s *NoticeServiceServer) Start() {
	{
		l := logic.NewTimerCleanSubscriptionLogic(s.svcCtx)
		go l.Start()
	}
}
