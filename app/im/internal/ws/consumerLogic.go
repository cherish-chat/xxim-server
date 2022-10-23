package ws

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
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
	go func() {
		pushConsumer := xtdmq.NewTDMQConsumer(l.svcCtx.Config.TDMQ.TDMQConfig, l.svcCtx.Config.TDMQ.PushConsumer)
		err := pushConsumer.Consume(context.Background(), l.pushConsumer)
		if err != nil {
			l.Errorf("pushConsumer.Consume error: %v", err)
			panic(err)
		}
	}()
	{
		internalConsumer := xtdmq.NewTDMQConsumer(l.svcCtx.Config.TDMQ.TDMQConfig, l.svcCtx.Config.TDMQ.InternalConsumer)
		err := internalConsumer.Consume(context.Background(), l.internalConsumer)
		if err != nil {
			l.Errorf("internalConsumer.Consume error: %v", err)
			panic(err)
		}
	}
}

func (l *ConsumerLogic) pushConsumer(ctx context.Context, topic string, key string, payload []byte) error {
	return NewConsumePushLogic(l.svcCtx, ctx).ConsumePush(key, payload)
}

func (l *ConsumerLogic) internalConsumer(ctx context.Context, topic string, key string, payload []byte) error {
	return NewConsumeInternalLogic(l.svcCtx, ctx).ConsumeInternal(key, payload)
}
