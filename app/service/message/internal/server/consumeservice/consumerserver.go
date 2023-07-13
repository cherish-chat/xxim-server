package consumeservice

import (
	"context"
	messageservicelogic "github.com/cherish-chat/xxim-server/app/service/message/internal/logic/messageservice"
	noticeservicelogic "github.com/cherish-chat/xxim-server/app/service/message/internal/logic/noticeservice"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xmq"
)

type ConsumerServer struct {
	svcCtx *svc.ServiceContext
}

func NewConsumerServer(svcCtx *svc.ServiceContext) *ConsumerServer {
	return &ConsumerServer{svcCtx: svcCtx}
}

func (s *ConsumerServer) Start() {
	s.svcCtx.MQ.RegisterHandler(xmq.TopicMessageInsert, func(ctx context.Context, topic string, msg []byte) error {
		return messageservicelogic.NewMessageInsertLogic(context.Background(), s.svcCtx).ConsumeMessage(topic, msg)
	})

	s.svcCtx.MQ.RegisterHandler(xmq.TopicNoticeSend, func(ctx context.Context, topic string, msg []byte) error {
		return noticeservicelogic.NewNoticeInsertLogic(context.Background(), s.svcCtx).ConsumeNotice(topic, msg)
	})
	go s.svcCtx.MQ.StartConsuming()
}
