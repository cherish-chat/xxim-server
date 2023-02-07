package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"time"
)

type SendMsgListTaskTask struct {
	sendMsgListSyncReq *pb.SendMsgListReq
	errChan            chan error
}
type ConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	SendMsgListTasks     chan *SendMsgListTaskTask
	sendMsgListSyncLogic *SendMsgListSyncLogic
}

func NewConsumerLogic(svcCtx *svc.ServiceContext) *ConsumerLogic {
	l := &ConsumerLogic{svcCtx: svcCtx}
	l.SendMsgListTasks = make(chan *SendMsgListTaskTask, svcCtx.Config.TDMQ.SendMsgListTaskNum)
	l.sendMsgListSyncLogic = NewSendMsgListSyncLogic(context.Background(), svcCtx)
	l.Logger = logx.WithContext(context.Background())
	go l.SendMsgListTask()
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
	err := pushConsumer.Consume(context.Background(), l.Consumer, xtdmq.ConsumerWithRc(l.svcCtx.Redis()))
	if err != nil {
		l.Errorf("pushConsumer.Consume error: %v", err)
		panic(err)
	}
}

func (l *ConsumerLogic) Consumer(ctx context.Context, topic string, key string, payload []byte) error {
	body := &pb.MsgMQBody{}
	err := proto.Unmarshal(payload, body)
	if err != nil {
		l.Errorf("proto.Unmarshal error: %v", err)
		return err
	}
	switch body.Event {
	case pb.MsgMQBody_SendMsgListSync:
		sendMsgListSyncReq := &pb.SendMsgListReq{}
		err = proto.Unmarshal(body.Data, sendMsgListSyncReq)
		if err != nil {
			l.Errorf("proto.Unmarshal error: %v", err)
			return err
		}
		// 加入task 并等待errChan返回
		errChan := make(chan error)
		l.SendMsgListTasks <- &SendMsgListTaskTask{
			sendMsgListSyncReq: sendMsgListSyncReq,
			errChan:            errChan,
		}
		err = <-errChan
		if err != nil {
			l.Errorf("SendMsgListSyncLogic.SendMsgListSync error: %v", err)
			return err
		}
	}
	return nil
}

// 定时器执行SendMsgListTask, 执行完成之后往 errChan 发送 nil 或者 error
func (l *ConsumerLogic) SendMsgListTask() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(l.svcCtx.Config.TDMQ.SendMsgListTaskInterval))
	for {
		select {
		case <-ticker.C:
			var msgDataList []*pb.MsgData
			var errChan []chan error
			msgDataList, errChan = l.PopSendMsgListTask(l.svcCtx.Config.TDMQ.SendMsgListTaskNum)
			if len(msgDataList) == 0 {
				continue
			}
			_, err := l.sendMsgListSyncLogic.SendMsgListSync(&pb.SendMsgListReq{
				MsgDataList: msgDataList,
			})
			if err != nil {
				l.Errorf("SendMsgListSyncLogic.SendMsgListSync error: %v", err)
			}
			for _, errChan := range errChan {
				errChan <- err
			}
		}
	}
}

func (l *ConsumerLogic) PopSendMsgListTask(num int) ([]*pb.MsgData, []chan error) {
	if len(l.SendMsgListTasks) == 0 {
		return nil, nil
	}
	var msgDataList []*pb.MsgData
	var errChan []chan error
	length := len(l.SendMsgListTasks)
	if length < num {
		// 如果任务数小于 num, 则全部执行
		num = length
	}
	for i := 0; i < num; i++ {
		task := <-l.SendMsgListTasks
		msgDataList = append(msgDataList, task.sendMsgListSyncReq.MsgDataList...)
		errChan = append(errChan, task.errChan)
	}
	return msgDataList, errChan
}
