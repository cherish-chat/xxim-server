package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

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
func (l *NoticeSendLogic) NoticeSend(in *peerpb.NoticeSendReq) (*peerpb.NoticeSendResp, error) {
	err := l.svcCtx.MQ.Produce(context.Background(), xmq.TopicNoticeSend, utils.Proto.Marshal(in))
	if err != nil {
		l.Errorf("produce message error: %v", err)
		return &peerpb.NoticeSendResp{}, err
	}
	return &peerpb.NoticeSendResp{}, nil
}
