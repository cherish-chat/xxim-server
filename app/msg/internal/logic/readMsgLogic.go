package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadMsgTask struct {
	readMsgReq *pb.ReadMsgReq
	errChan    chan error
}

func (l *ReadMsgLogic) readMsgTask() {
	ticker := time.NewTicker(l.svcCtx.ConfigMgr.ReadMsgTaskInterval(l.ctx))
	for {
		select {
		case <-ticker.C:
			notices, errChans := l.popReadMsgTask(l.svcCtx.ConfigMgr.ReadMsgTaskBatchSize(l.ctx))
			if len(notices) == 0 {
				continue
			}
			// 批量插入
			go l.handleReadMsgTask(notices, errChans)
		}
	}
}

type ReadMsgLogic struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	readMsgTasks chan *ReadMsgTask
	logx.Logger
}

var singleReadMsgLogic *ReadMsgLogic

func NewReadMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadMsgLogic {
	ctx = context.Background()
	if singleReadMsgLogic == nil {
		singleReadMsgLogic = &ReadMsgLogic{
			ctx:          ctx,
			svcCtx:       svcCtx,
			readMsgTasks: make(chan *ReadMsgTask, 0),
			Logger:       logx.WithContext(ctx),
		}
		go singleReadMsgLogic.readMsgTask()
	}
	return singleReadMsgLogic
}

// ReadMsg 设置会话已读
func (l *ReadMsgLogic) ReadMsg(in *pb.ReadMsgReq) (*pb.ReadMsgResp, error) {
	errChan := make(chan error, 0)
	l.readMsgTasks <- &ReadMsgTask{
		readMsgReq: in,
		errChan:    errChan,
	}
	err := <-errChan
	if err != nil {
		return &pb.ReadMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.ReadMsgResp{CommonResp: pb.NewSuccessResp()}, nil
}

func (l *ReadMsgLogic) readMsg(in *pb.ReadMsgReq) (*pb.ReadMsgResp, error) {
	notice := &noticemodel.Notice{
		ConvId: pb.HiddenConvId(in.ConvId),
		Options: noticemodel.NoticeOption{
			StorageForClient: false,
			UpdateConvNotice: false,
		},
		ContentType: int32(pb.NoticeType_READ),
		Content:     in.NoticeContent,
		UniqueId:    fmt.Sprintf("readSeq:%s-%s", in.SenderId, in.Seq),
		Title:       "",
		Ext:         nil,
	}
	err := notice.Insert(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis())
	if err != nil {
		l.Errorf("insert notice failed, err: %v", err)
		return &pb.ReadMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 通知
	xtrace.StartFuncSpan(l.ctx, "SendNotice", func(ctx context.Context) {
		utils.RetryProxy(ctx, 12, time.Second, func() error {
			_, err := l.svcCtx.NoticeService().GetUserNoticeData(ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				UserId:    "",
				ConvId:    pb.HiddenConvId(in.ConvId),
			})
			if err != nil {
				l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
				return err
			}
			return nil
		})
	})
	return &pb.ReadMsgResp{}, nil
}

func (l *ReadMsgLogic) popReadMsgTask(num int) ([]*noticemodel.Notice, []chan error) {
	if len(l.readMsgTasks) == 0 {
		return nil, nil
	}
	var (
		tasks    []*ReadMsgTask
		notices  []*noticemodel.Notice
		errChans []chan error
	)
	if len(l.readMsgTasks) < num {
		num = len(l.readMsgTasks)
	}
	for i := 0; i < num; i++ {
		task := <-l.readMsgTasks
		tasks = append(tasks, task)
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvId(task.readMsgReq.ConvId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: int32(pb.NoticeType_READ),
			Content:     task.readMsgReq.NoticeContent,
			UniqueId:    fmt.Sprintf("readSeq:%s-%s", task.readMsgReq.SenderId, task.readMsgReq.Seq),
			Title:       "",
			Ext:         nil,
		}
		notices = append(notices, notice)
		errChans = append(errChans, task.errChan)
	}
	return notices, errChans
}

func (l *ReadMsgLogic) handleReadMsgTask(notices []*noticemodel.Notice, errChans []chan error) {
	err := noticemodel.BatchInsert(l.svcCtx.Mysql(), notices, l.svcCtx.Redis())
	for _, errChan := range errChans {
		errChan <- err
	}
	// 通知
	xtrace.StartFuncSpan(l.ctx, "SendNotices", func(ctx context.Context) {
		var convIds []string
		for _, notice := range notices {
			convIds = append(convIds, notice.ConvId)
		}
		convIds = utils.Set(convIds)
		for _, convId := range convIds {
			utils.RetryProxy(ctx, 12, time.Second, func() error {
				_, err := l.svcCtx.NoticeService().GetUserNoticeData(ctx, &pb.GetUserNoticeDataReq{
					CommonReq: &pb.CommonReq{},
					UserId:    "",
					ConvId:    pb.HiddenConvId(convId),
				})
				if err != nil {
					l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
					return err
				}
				return nil
			})
		}
	})
}
