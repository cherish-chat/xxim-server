package server

import (
	"context"
	noticeservicelogic "github.com/cherish-chat/xxim-server/app/message/internal/logic/noticeservice"
	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xmq"
)

type ConsumerServer struct {
	svcCtx *svc.ServiceContext
}

func NewConsumerServer(svcCtx *svc.ServiceContext) *ConsumerServer {
	return &ConsumerServer{svcCtx: svcCtx}
}

func (s *ConsumerServer) Start() {
	s.svcCtx.MQ.RegisterHandler(xmq.TopicNoticeBatchSend, func(ctx context.Context, topic string, msg []byte) error {
		return noticeservicelogic.NewConsumerLogic(ctx, s.svcCtx).NoticeBatchSend(topic, msg)
	})
	go s.svcCtx.MQ.StartConsuming()
}
