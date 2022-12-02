package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSetMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchSetMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSetMinSeqLogic {
	return &BatchSetMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchSetMinSeq 批量设置用户某会话的minseq
func (l *BatchSetMinSeqLogic) BatchSetMinSeq(in *pb.BatchSetMinSeqReq) (*pb.BatchSetMinSeqResp, error) {
	if len(in.UserIdList) == 0 {
		return &pb.BatchSetMinSeqResp{}, nil
	}
	// hmset key field value [field value ...]
	minSeqMap := make(map[string]string)
	for _, userId := range in.UserIdList {
		minSeqMap[rediskey.HKConvMinSeq(userId)] = in.MinSeq
	}
	err := l.svcCtx.Redis().HmsetCtx(l.ctx, rediskey.ConvKv(in.ConvId), minSeqMap)
	if err != nil {
		l.Errorf("BatchSetMinSeq HmsetCtx error: %v", err)
		return &pb.BatchSetMinSeqResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.BatchSetMinSeqResp{}, nil
}
