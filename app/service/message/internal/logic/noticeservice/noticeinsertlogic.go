package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"
	"github.com/cherish-chat/xxim-server/app/service/message/messagemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NoticeInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeInsertLogic {
	return &NoticeInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type xSendMsgListTask struct {
	sendMsgListSyncReq *peerpb.NoticeSendReq
	errChan            chan error
}

var xSendMsgListTaskChan chan *xSendMsgListTask

func sendMsgListTaskChan(svcCtx *svc.ServiceContext) chan *xSendMsgListTask {
	if xSendMsgListTaskChan == nil {
		xSendMsgListTaskChan = make(chan *xSendMsgListTask, 100)
		go loopInsertMsgList(svcCtx)
	}
	return xSendMsgListTaskChan
}

func loopInsertMsgList(svcCtx *svc.ServiceContext) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(50))
	for {
		select {
		case <-ticker.C:
			var msgDataList []*peerpb.Message
			var errChan []chan error
			msgDataList, errChan = popSendMsgListTask(100)
			if len(msgDataList) == 0 {
				continue
			}
			if len(msgDataList) > 1 {
				logx.Infof("SendMsgListSyncLogic.SendMsgListSync msgDataList.len: %d", len(msgDataList))
			}
			_, err := insertMsgList(svcCtx, &peerpb.NoticeSendReq{
				Notices: msgDataList,
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

func popSendMsgListTask(num int) ([]*peerpb.Message, []chan error) {
	if len(xSendMsgListTaskChan) == 0 {
		return nil, nil
	}
	var msgDataList []*peerpb.Message
	var errChan []chan error
	length := len(xSendMsgListTaskChan)
	if length < num {
		// 如果任务数小于 num, 则全部执行
		num = length
	}
	for i := 0; i < num; i++ {
		task := <-xSendMsgListTaskChan
		msgDataList = append(msgDataList, task.sendMsgListSyncReq.Notices...)
		errChan = append(errChan, task.errChan)
	}
	return msgDataList, errChan
}

func insertMsgList(svcCtx *svc.ServiceContext, in *peerpb.NoticeSendReq) (*peerpb.NoticeSendResp, error) {
	now := time.Now()
	ctx := context.Background()
	var insertNotices = make([]*messagemodel.Message, 0)
	var newNotices = make([]*peerpb.Message, 0)
	for _, message := range in.Notices {
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
			seq, err := messagemodel.RedisSeq.IncrConvNoticeMaxSeq(ctx, messageFromPb.ConversationId, messageFromPb.ConversationType)
			if err != nil {
				logx.Errorf("NoticeInsert RedisSeq.IncrConvNoticeMaxSeq error: %v", err)
				return nil, err
			}
			messageFromPb.Seq = int64(seq)
			//messageId
			messageFromPb.GenerateMessageId()

			insertNotices = append(insertNotices, messageFromPb)
			newNotices = append(newNotices, messageFromPb.ToPb())
		} else {
			newNotices = append(newNotices, message)
		}
	}
	if len(insertNotices) > 0 {
		err := xmgo.BatchInsertMany(svcCtx.NoticeCollection, ctx, insertNotices, 1000)
		if err != nil {
			logx.Errorf("NoticeInsert BatchInsertMany error: %v", err)
			return nil, err
		}
	}

	go func() {
		_, _ = NewNoticePushLogic(ctx, svcCtx).NoticePush(&peerpb.NoticeSendReq{
			Header:  in.Header,
			Notices: newNotices,
		})
	}()
	return &peerpb.NoticeSendResp{}, nil
}

// NoticeInsert 插入消息
func (l *NoticeInsertLogic) NoticeInsert(in *peerpb.NoticeSendReq) (*peerpb.NoticeSendResp, error) {
	errChan := make(chan error)
	sendMsgListTaskChan(l.svcCtx) <- &xSendMsgListTask{
		sendMsgListSyncReq: in,
		errChan:            errChan,
	}
	err := <-errChan
	if err != nil {
		l.Errorf("NoticeInsert error: %v", err)
		return &peerpb.NoticeSendResp{}, err
	}
	return &peerpb.NoticeSendResp{}, nil
}

// ConsumeNotice 消费消息
func (l *NoticeInsertLogic) ConsumeNotice(topic string, msg []byte) error {
	in := &peerpb.NoticeSendReq{}
	err := utils.Json.Unmarshal(msg, in)
	if err != nil {
		l.Errorf("ConsumeNotice json.Unmarshal error: %v", err)
		return nil
	}
	_, err = l.NoticeInsert(in)
	return err
}
