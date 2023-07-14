package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/message/messagemodel"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConvMessageSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConvMessageSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConvMessageSeqLogic {
	return &GetConvMessageSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetConvMessageSeq 获取会话消息序列号
func (l *GetConvMessageSeqLogic) GetConvMessageSeq(in *peerpb.GetConvMessageSeqReq) (*peerpb.GetConvMessageSeqResp, error) {
	convSeq, err := messagemodel.RedisSeq.GetConvMaxSeq(context.Background(), in.ConvId, in.ConvType)
	if err != nil {
		l.Errorf("get conv message seq error: %v", err)
		return &peerpb.GetConvMessageSeqResp{}, err
	}
	minSeq, err := messagemodel.RedisSeq.GetConvMessageMinSeq(context.Background(), in.ConvId, in.Header.UserId, in.ConvType)
	if err != nil {
		l.Errorf("get conv message seq error: %v", err)
		return &peerpb.GetConvMessageSeqResp{}, err
	}
	return &peerpb.GetConvMessageSeqResp{
		ConvId:   in.ConvId,
		ConvType: in.ConvType,
		MinSeq:   uint32(minSeq),
		MaxSeq:   uint32(convSeq.MessageMaxSeq),
	}, nil
}
