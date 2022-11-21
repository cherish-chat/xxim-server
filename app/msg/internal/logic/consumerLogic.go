package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
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
		_, err = NewSendMsgListSyncLogic(ctx, l.svcCtx).SendMsgListSync(sendMsgListSyncReq)
		if err != nil {
			l.Errorf("SendMsgListSyncLogic.SendMsgListSync error: %v", err)
			return err
		}
		//case pb.MsgMQBody_BatchSendMsgSync:
		//	batchSendMsgSyncReq := &pb.BatchSendMsgReq{}
		//	err = proto.Unmarshal(body.Data, batchSendMsgSyncReq)
		//	if err != nil {
		//		l.Errorf("proto.Unmarshal error: %v", err)
		//		return err
		//	}
		//	_, err = NewBatchSendMsgSyncLogic(ctx, l.svcCtx).BatchSendMsgSync(batchSendMsgSyncReq)
		//	if err != nil {
		//		l.Errorf("BatchSendMsgSyncLogic.BatchSendMsgSync error: %v", err)
		//		return err
		//	}
	}
	return nil
}
