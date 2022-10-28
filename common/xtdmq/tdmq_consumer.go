package xtdmq

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	_ "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ConsumeFunc func(ctx context.Context, topic string, msgKey string, payload []byte) error
type ConsumeOpt struct {
}
type ConsumerOptFunc func(options *ConsumeOpt)

type TDMQConsumer struct {
	Config         TDMQConfig
	ConsumerConfig TDMQConsumerConfig
	consumer       pulsar.Consumer
	client         pulsar.Client
}

func NewTDMQConsumer(config TDMQConfig, consumerConfig TDMQConsumerConfig) *TDMQConsumer {
	p := &TDMQConsumer{Config: config, ConsumerConfig: consumerConfig}
	p.init()
	return p
}

func (p *TDMQConsumer) init() {
	// 创建pulsar客户端
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		// 服务接入地址
		URL: p.Config.VpcUrl,
		// 授权角色密钥
		Authentication:    pulsar.NewAuthenticationToken(p.Config.Token),
		OperationTimeout:  3 * time.Second,
		ConnectionTimeout: 3 * time.Second,
	})
	if err != nil {
		logx.Errorf("Could not instantiate Pulsar client: %v", err)
		panic(err)
	}
	p.client = client
	// 使用客户端创建生产者
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		// topic完整路径，格式为persistent://集群（租户）ID/命名空间/Topic名称
		Topic:                       p.ConsumerConfig.TopicName,
		SubscriptionName:            p.ConsumerConfig.GetSubName(),
		Type:                        pulsar.SubscriptionType(p.ConsumerConfig.SubType),
		SubscriptionInitialPosition: pulsar.SubscriptionInitialPosition(p.ConsumerConfig.SubInitialPosition),
		RetryEnable:                 p.ConsumerConfig.EnableRetry,
		ReceiverQueueSize:           p.ConsumerConfig.ReceiverQueueSize,
		Name:                        p.ConsumerConfig.GetConsumerName(),
	})
	if err != nil {
		logx.Errorf("Could not instantiate Pulsar consumer: %v", err)
		panic(err)
	}
	p.consumer = consumer
}

func (p *TDMQConsumer) Consume(
	ctx context.Context,
	f ConsumeFunc,
	opts ...ConsumerOptFunc,
) error {
	options := &ConsumeOpt{}
	for _, opt := range opts {
		opt(options)
	}
	for {
		receive, err := p.consumer.Receive(ctx)
		if err != nil {
			logx.Errorf("Could not receive message: %v", err)
			return err
		}
		// 获取 traceId
		traceId, _ := receive.Properties()["traceId"]
		xtrace.RunWithTrace(
			traceId,
			fmt.Sprintf("tdmqconsumer/topic:%s/subname:%s/consumername:%s", p.ConsumerConfig.TopicName, p.ConsumerConfig.SubName, p.ConsumerConfig.ConsumerName),
			func(ctx context.Context) {
				err := f(ctx, receive.Topic(), receive.Key(), receive.Payload())
				if err != nil {
					// 消费失败，重试
					// https://cloud.tencent.com/document/product/1179/49607
					p.consumer.Nack(receive)
				} else {
					// 消费成功，确认消费
					p.consumer.Ack(receive)
				}
			},
		)
	}
}
