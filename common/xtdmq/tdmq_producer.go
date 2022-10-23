package xtdmq

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ProduceOpt struct {
	properties   *map[string]string
	deliverAfter *time.Duration
	deliverAt    *time.Time
}
type ProducerOptFunc func(options *ProduceOpt)

func ProduceWithProperties(properties map[string]string) ProducerOptFunc {
	return func(options *ProduceOpt) {
		options.properties = &properties
	}
}

// ProduceWithDeliverAfter 发送延迟消息
// 注意：此方法只能发送延迟消息 消费者订阅类型必须为Shared
func ProduceWithDeliverAfter(delay time.Duration) ProducerOptFunc {
	return func(options *ProduceOpt) {
		options.deliverAfter = &delay
	}
}

// ProduceWithDeliverAt 发送定时消息
// 注意：此方法只能发送定时消息 消费者订阅类型必须为Shared
func ProduceWithDeliverAt(t time.Time) ProducerOptFunc {
	return func(options *ProduceOpt) {
		options.deliverAt = &t
	}
}

type TDMQProducer struct {
	Config         TDMQConfig
	ProducerConfig TDMQProducerConfig
	producer       pulsar.Producer
	client         pulsar.Client
}

func NewTDMQProducer(config TDMQConfig, producerConfig TDMQProducerConfig) *TDMQProducer {
	p := &TDMQProducer{Config: config, ProducerConfig: producerConfig}
	p.init()
	return p
}

func (p *TDMQProducer) init() {
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
	GetApiClient(p.ProducerConfig.TopicName, p.Config)
	// 使用客户端创建生产者
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		// topic完整路径，格式为persistent://集群（租户）ID/命名空间/Topic名称
		Topic:       p.ProducerConfig.TopicName,
		Name:        p.ProducerConfig.GetProducerName(),
		SendTimeout: time.Duration(p.ProducerConfig.SendTimeout) * time.Millisecond,
	})
	if err != nil {
		logx.Errorf("Could not instantiate Pulsar producer: %v", err)
		panic(err)
	}
	p.producer = producer
}

func (p *TDMQProducer) Produce(
	ctx context.Context,
	key string,
	payload []byte,
	opts ...ProducerOptFunc,
) (string, error) {
	options := &ProduceOpt{}
	for _, opt := range opts {
		opt(options)
	}
	msg := &pulsar.ProducerMessage{
		Payload: payload,
		Key:     key,
	}
	if options.properties != nil {
		tmp := *options.properties
		tmp["traceId"] = xtrace.TraceIdFromContext(ctx)
		msg.Properties = tmp
	} else {
		msg.Properties = map[string]string{"traceId": xtrace.TraceIdFromContext(ctx)}
	}
	if options.deliverAfter != nil {
		msg.DeliverAfter = *options.deliverAfter
	} else if options.deliverAt != nil {
		msg.DeliverAt = *options.deliverAt
	}
	// 发送消息
	var err error
	var msgId string
	xtrace.StartFuncSpan(ctx, fmt.Sprintf("tdmqproducer/topic:%s/producername:%s", p.ProducerConfig.TopicName, p.ProducerConfig.GetProducerName()), func(ctx context.Context) {
		msgID, e := p.producer.Send(ctx, msg)
		if e != nil {
			logx.Errorf("Could not publish message: %v", err)
			msgId = ""
			err = e
			return
		}
		msgId = msgID.(fmt.Stringer).String()
		return
	})
	return msgId, err
}

func (p *TDMQProducer) SendMessage(ctx context.Context, v []byte, id string) (interface{}, string, error) {
	id, err := p.Produce(ctx, id, v)
	return nil, id, err
}
