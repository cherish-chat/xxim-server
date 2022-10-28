package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

var singletonConsumerLogic *ConsumerLogic

func NewConsumerLogic(svcCtx *svc.ServiceContext) *ConsumerLogic {
	if singletonConsumerLogic == nil {
		c := &ConsumerLogic{
			svcCtx: svcCtx,
			Logger: logx.WithContext(context.Background()),
		}
		singletonConsumerLogic = c
	}
	return singletonConsumerLogic
}

func (l *ConsumerLogic) Start() {
	pushConsumer := xtdmq.NewTDMQConsumer(l.svcCtx.Config.TDMQ.TDMQConfig, xtdmq.TDMQConsumerConfig{
		TopicName:          l.svcCtx.Config.TDMQ.TopicName,
		SubName:            "conn",
		ConsumerName:       "conn",
		SubInitialPosition: 0,
		SubType:            0,
		EnableRetry:        true,
		ReceiverQueueSize:  l.svcCtx.Config.TDMQ.ReceiverQueueSize,
		IsBroadcast:        true,
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
