package server

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/logic"
)

func (s *ImServiceServer) IMServer() {
	l := logic.NewConsumerStorage(context.Background(), s.svcCtx)
	l.Consume()
}
