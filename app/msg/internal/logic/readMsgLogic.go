package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadMsgLogic {
	return &ReadMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReadMsg 设置会话已读
func (l *ReadMsgLogic) ReadMsg(in *pb.ReadMsgReq) (*pb.ReadMsgResp, error) {
	notice := &noticemodel.Notice{
		ConvId: pb.HiddenConvId(in.ConvId),
		Options: noticemodel.NoticeOption{
			StorageForClient: false,
			UpdateConvNotice: false,
		},
		ContentType: int32(pb.NoticeType_READ),
		Content:     in.NoticeContent,
		UniqueId:    "info",
		Title:       "",
		Ext:         nil,
	}
	err := notice.Insert(l.ctx, l.svcCtx.Mysql())
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
