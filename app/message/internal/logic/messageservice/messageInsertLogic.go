package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/message/messagemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInsertLogic {
	return &MessageInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type xSendMsgListTask struct {
	sendMsgListSyncReq *pb.MessageInsertReq
	errChan            chan error
}

var xSendMsgListTaskChan chan *xSendMsgListTask

func sendMsgListTaskChan(svcCtx *svc.ServiceContext) chan *xSendMsgListTask {
	if xSendMsgListTaskChan == nil {
		xSendMsgListTaskChan = make(chan *xSendMsgListTask, svcCtx.Config.InsertMsgBuffer.Size)
		go loopInsertMsgList(svcCtx)
	}
	return xSendMsgListTaskChan
}

func loopInsertMsgList(svcCtx *svc.ServiceContext) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(svcCtx.Config.InsertMsgBuffer.LoopInterval))
	for {
		select {
		case <-ticker.C:
			var msgDataList []*pb.Message
			var errChan []chan error
			msgDataList, errChan = popSendMsgListTask(svcCtx.Config.InsertMsgBuffer.Size)
			if len(msgDataList) == 0 {
				continue
			}
			_, err := insertMsgList(svcCtx, &pb.MessageInsertReq{
				Messages: msgDataList,
			})
			if err != nil {
				logx.Errorf("SendMsgListSyncLogic.SendMsgListSync error: %v", err)
			}
			for _, errChan := range errChan {
				errChan <- err
			}
		}
	}
}

func popSendMsgListTask(num int) ([]*pb.Message, []chan error) {
	if len(xSendMsgListTaskChan) == 0 {
		return nil, nil
	}
	var msgDataList []*pb.Message
	var errChan []chan error
	length := len(xSendMsgListTaskChan)
	if length < num {
		// 如果任务数小于 num, 则全部执行
		num = length
	}
	for i := 0; i < num; i++ {
		task := <-xSendMsgListTaskChan
		msgDataList = append(msgDataList, task.sendMsgListSyncReq.Messages...)
		errChan = append(errChan, task.errChan)
	}
	return msgDataList, errChan
}

func insertMsgList(svcCtx *svc.ServiceContext, in *pb.MessageInsertReq) (*pb.MessageInsertResp, error) {
	now := time.Now()
	ctx := context.Background()
	var insertMessages = make([]*messagemodel.Message, 0)
	var newMessages = make([]*pb.Message, 0)
	for _, message := range in.Messages {
		if message.GetOption().GetStorageForServer() {
			messageFromPb := messagemodel.MessageFromPb(message)
			//赋值
			if messageFromPb.SendTime == 0 {
				messageFromPb.SendTime = primitive.NewDateTimeFromTime(now)
			}
			if messageFromPb.InsertTime == 0 {
				messageFromPb.InsertTime = primitive.NewDateTimeFromTime(now)
			}
			if messageFromPb.Uuid == "" {
				messageFromPb.Uuid = utils.Snowflake.String()
			}
			//seq
			seq, err := messagemodel.RedisSeq.IncrConvMaxSeq(ctx, messageFromPb.ConversationId, messageFromPb.ConversationType)
			if err != nil {
				logx.Errorf("MessageInsert RedisSeq.IncrConvMaxSeq error: %v", err)
				return nil, err
			}
			messageFromPb.Seq = int64(seq)
			//messageId
			messageFromPb.GenerateMessageId()

			insertMessages = append(insertMessages, messageFromPb)
			newMessages = append(newMessages, messageFromPb.ToPb())
		} else {
			newMessages = append(newMessages, message)
		}
	}
	if len(insertMessages) > 0 {
		err := xmgo.BatchInsertMany(svcCtx.MessageCollection, ctx, insertMessages, 1000)
		if err != nil {
			logx.Errorf("MessageInsert BatchInsertMany error: %v", err)
			return nil, err
		}
	}

	_, err := NewMessagePushLogic(ctx, svcCtx).MessagePush(&pb.MessagePushReq{
		Header:  in.Header,
		Message: newMessages,
	})
	return &pb.MessageInsertResp{}, err
}

// MessageInsert 插入消息
func (l *MessageInsertLogic) MessageInsert(in *pb.MessageInsertReq) (*pb.MessageInsertResp, error) {
	errChan := make(chan error)
	sendMsgListTaskChan(l.svcCtx) <- &xSendMsgListTask{
		sendMsgListSyncReq: in,
		errChan:            errChan,
	}
	err := <-errChan
	if err != nil {
		l.Errorf("MessageInsert error: %v", err)
		return &pb.MessageInsertResp{}, err
	}
	return &pb.MessageInsertResp{}, nil
}

// ConsumeMessage 消费消息
func (l *MessageInsertLogic) ConsumeMessage(topic string, msg []byte) error {
	in := &pb.MessageInsertReq{}
	err := utils.Json.Unmarshal(msg, in)
	if err != nil {
		l.Errorf("ConsumeMessage json.Unmarshal error: %v", err)
		return nil
	}
	_, err = l.MessageInsert(in)
	return err
}
