package xmq

import (
	"context"
	"time"
)

type Topic = string

type ProduceOption struct {
	// Delay 延迟时间
	delay *time.Duration
}

type ProduceOptionFunction func(o *ProduceOption)

// ProduceWithDelay 设置延迟时间
func ProduceWithDelay(delay time.Duration) ProduceOptionFunction {
	return func(o *ProduceOption) {
		o.delay = &delay
	}
}

type HandlerFunc func(ctx context.Context, topic string, msg []byte) error

type MQ interface {
	// Produce 生产消息
	Produce(ctx context.Context, topic Topic, msg []byte, opts ...ProduceOptionFunction) error
	RegisterHandler(topic Topic, handler HandlerFunc)
	StartConsuming()
}
