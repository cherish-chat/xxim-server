package xtdmq

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	_ "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/zeromicro/go-zero/core/logx"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"go.opentelemetry.io/otel/propagation"
	"time"
)

type ConsumeFunc func(ctx context.Context, topic string, msgKey string, payload []byte) error
type ConsumeOpt struct {
	rc *zedis.Redis
}
type ConsumerOptFunc func(options *ConsumeOpt)

func ConsumerWithRc(rc *zedis.Redis) ConsumerOptFunc {
	return func(options *ConsumeOpt) {
		options.rc = rc
	}
}

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
		properties := receive.Properties()
		{
			properties["publishTime"] = receive.PublishTime().String()
			properties["eventTime"] = receive.EventTime().String()
			properties["now"] = time.Now().String()
			properties["key"] = receive.Key()
			properties["orderingKey"] = receive.OrderingKey()
			properties["producerName"] = receive.ProducerName()
			properties["topic"] = receive.Topic()
		}
		traceId, _ := properties["traceId"]
		xtrace.RunWithTrace(
			traceId,
			fmt.Sprintf("tdmqconsumer/topic:%s/subname:%s/consumername:%s", p.ConsumerConfig.TopicName, p.ConsumerConfig.SubName, p.ConsumerConfig.ConsumerName),
			func(ctx context.Context) {
				err := f(ctx, receive.Topic(), receive.Key(), receive.Payload())
				if err != nil {
					// 消费失败，重试
					// https://cloud.tencent.com/document/product/1179/49607
					if options.rc != nil {
						key := rediskey.MQRetryCount(traceId)
						count, err := xredis.IncrEx(options.rc, ctx, key, 24*60*60, 1)
						if err != nil {
							logx.Errorf("redis incr error:%v", err)
						}
						if count == 12 {
							// 重试12次，放弃，TODO 告警
						}
						p.consumer.ReconsumeLater(receive, GetRetryDelay(count))
					} else {
						p.consumer.Nack(receive)
					}
				} else {
					// 消费成功，确认消费
					p.consumer.Ack(receive)
				}
			},
			propagation.MapCarrier(properties),
		)
	}
}

var retryDelayMap = map[int64]time.Duration{
	0:  1 * time.Second,
	1:  10 * time.Second,
	2:  10 * time.Second,
	3:  10 * time.Second,
	4:  10 * time.Second,
	5:  10 * time.Second,
	6:  20 * time.Second,
	7:  20 * time.Second,
	8:  20 * time.Second,
	9:  20 * time.Second,
	10: 30 * time.Second,
	11: 30 * time.Second,
	12: 30 * time.Second,
}

func GetRetryDelay(times int64) time.Duration {
	if times >= 12 {
		return retryDelayMap[12]
	}
	return retryDelayMap[times]
}
