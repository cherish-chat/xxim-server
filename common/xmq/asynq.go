package xmq

import (
	"context"
	"crypto/tls"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
	"os/signal"
	"time"
)

type Asynq struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
}

type logxLogger struct {
}

func (l *logxLogger) Debug(args ...interface{}) {
	logx.Debug(args...)
}

func (l *logxLogger) Info(args ...interface{}) {
	logx.Info(args...)
}

func (l *logxLogger) Warn(args ...interface{}) {
	logx.Info(args...)
}

func (l *logxLogger) Error(args ...interface{}) {
	logx.Error(args...)
}

func (l *logxLogger) Fatal(args ...interface{}) {
	logx.Error(args...)
	os.Exit(1)
}

func NewAsynq(
	redisConf redis.RedisConf,
	db int,
	logLevel string, // options=[debug,info,error,severe]
) MQ {
	mqLogLevel := asynq.DebugLevel
	switch logLevel {
	case "debug":
		mqLogLevel = asynq.DebugLevel
	case "info":
		mqLogLevel = asynq.InfoLevel
	case "error":
		mqLogLevel = asynq.ErrorLevel
	case "severe":
		mqLogLevel = asynq.FatalLevel
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if !redisConf.Tls {
		tlsConfig = nil
	}
	opt := asynq.RedisClientOpt{
		Addr:         redisConf.Host,
		Password:     redisConf.Pass,
		DB:           db,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  0,
		WriteTimeout: 0,
		PoolSize:     0,
		TLSConfig:    tlsConfig,
	}
	client := asynq.NewClient(opt)
	server := asynq.NewServer(opt, asynq.Config{
		// 任务的最大并发处理数。如果设置为零或负值，NewServer 会将该值覆盖为当前进程可用的 CPU 数。
		Concurrency: 0,
		// Logger 用于记录日志。如果未设置，将使用默认的 Logger。
		Logger: new(logxLogger),
		// LoggerLevel 用于设置日志级别。默认为 Info 级别。
		LogLevel: mqLogLevel,
	})
	mux := asynq.NewServeMux()
	return &Asynq{
		client: client,
		server: server,
		mux:    mux,
	}
}

func (a *Asynq) RegisterHandler(
	topic string,
	handler HandlerFunc,
) {
	a.mux.HandleFunc(topic, func(ctx context.Context, task *asynq.Task) error {
		return handler(ctx, task.Type(), task.Payload())
	})
}

func (a *Asynq) StartConsuming() {
	// 监听 Ctrl+C 信号 如果收到则优雅退出
	go func() {
		err := a.server.Run(a.mux)
		if err != nil {
			logx.Errorf("start consuming error: %v", err)
			os.Exit(1)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logx.Info("shutdown signal received, exiting...")
	a.server.Stop()
	a.server.Shutdown()
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func (a *Asynq) Produce(ctx context.Context, topic string, msg []byte, opts ...ProduceOptionFunction) error {
	task := asynq.NewTask(topic, msg)
	opt := a.defaultProduceOption()
	for _, o := range opts {
		o(opt)
	}
	var enqueueOptions []asynq.Option
	if opt.delay != nil {
		enqueueOptions = append(enqueueOptions, asynq.ProcessIn(*opt.delay))
	}
	_, err := a.client.EnqueueContext(ctx, task, enqueueOptions...)
	return err
}

func (a *Asynq) defaultProduceOption() *ProduceOption {
	return &ProduceOption{}
}
