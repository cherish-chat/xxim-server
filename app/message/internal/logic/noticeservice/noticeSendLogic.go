package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeSendLogic {
	return &NoticeSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// NoticeSend 通知发送
func (l *NoticeSendLogic) NoticeSend(in *pb.NoticeSendReq) (*pb.NoticeSendResp, error) {
	req := &pb.NoticeBatchSendReq{
		Header:  in.Header,
		Notices: []*pb.NoticeSendReq{in},
	}
	err := l.svcCtx.MQ.Produce(l.ctx, xmq.TopicNoticeBatchSend, utils.Json.MarshalToBytes(req))
	if err != nil {
		l.Errorf("produce message error: %v", err)
		return &pb.NoticeSendResp{}, err
	}
	return &pb.NoticeSendResp{}, nil
}
