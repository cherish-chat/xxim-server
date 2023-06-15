package server

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/logic"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xmq"
)

type ConsumerServer struct {
	svcCtx *svc.ServiceContext
}

func NewConsumerServer(svcCtx *svc.ServiceContext) *ConsumerServer {
	return &ConsumerServer{svcCtx: svcCtx}
}

func (s *ConsumerServer) Start() {
	consumerLogic := logic.NewConsumerLogic(context.Background(), s.svcCtx)
	s.svcCtx.MQ.RegisterHandler(xmq.TopicAfterRegister, consumerLogic.AfterRegister)
	go s.svcCtx.MQ.StartConsuming()
}
