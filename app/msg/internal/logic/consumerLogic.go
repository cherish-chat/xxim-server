package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumerLogic(svcCtx *svc.ServiceContext) *ConsumerLogic {
	l := &ConsumerLogic{svcCtx: svcCtx}
	l.Logger = logx.WithContext(context.Background())
	return l
}

func (l *ConsumerLogic) Start() {
	pushConsumer := xtdmq.NewTDMQConsumer(l.svcCtx.Config.TDMQ.TDMQConfig, xtdmq.TDMQConsumerConfig{
		TopicName:          l.svcCtx.Config.TDMQ.TopicName,
		SubName:            "msg",
		ConsumerName:       "msg",
		SubInitialPosition: 0,
		SubType:            1,
		EnableRetry:        true,
		ReceiverQueueSize:  l.svcCtx.Config.TDMQ.ReceiverQueueSize,
		IsBroadcast:        false,
	})
	err := pushConsumer.Consume(context.Background(), l.Consumer)
	if err != nil {
		l.Errorf("pushConsumer.Consume error: %v", err)
		panic(err)
	}
}

func (l *ConsumerLogic) Consumer(ctx context.Context, topic string, key string, payload []byte) error {
	// TODO: 业务逻辑
	return nil
}
