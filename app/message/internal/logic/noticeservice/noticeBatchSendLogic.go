package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeBatchSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeBatchSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeBatchSendLogic {
	return &NoticeBatchSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// NoticeBatchSend 通知批量发送
func (l *NoticeBatchSendLogic) NoticeBatchSend(in *pb.NoticeBatchSendReq) (*pb.NoticeBatchSendResp, error) {
	err := l.svcCtx.MQ.Produce(l.ctx, xmq.TopicNoticeBatchSend, utils.Json.MarshalToBytes(in))
	if err != nil {
		l.Errorf("produce message error: %v", err)
		return &pb.NoticeBatchSendResp{}, err
	}
	return &pb.NoticeBatchSendResp{}, nil
}
