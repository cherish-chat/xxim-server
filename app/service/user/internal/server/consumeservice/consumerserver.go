package consumeservice

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/logic/consumelogic"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xmq"
)

type ConsumerServer struct {
	svcCtx *svc.ServiceContext
}

func NewConsumerServer(svcCtx *svc.ServiceContext) *ConsumerServer {
	return &ConsumerServer{svcCtx: svcCtx}
}

func (s *ConsumerServer) Start() {
	s.svcCtx.MQ.RegisterHandler(xmq.TopicAfterRegister, func(ctx context.Context, topic string, msg []byte) error {
		return consumelogic.NewUserAfterRegisterLogic(ctx, s.svcCtx).AfterRegister(topic, msg)
	})
	go s.svcCtx.MQ.StartConsuming()
}
