package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLatestConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLatestConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLatestConnLogic {
	return &GetUserLatestConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLatestConnLogic) GetUserLatestConn(in *pb.GetUserLatestConnReq) (*pb.GetUserLatestConnResp, error) {
	var batchResp *pb.BatchGetUserLatestConnResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "BatchGetUserLatestConn", func(ctx context.Context) {
		batchResp, err = NewBatchGetUserLatestConnLogic(ctx, l.svcCtx).BatchGetUserLatestConn(&pb.BatchGetUserLatestConnReq{UserIds: []string{in.UserId}})
	})
	if err != nil {
		l.Errorf("GetUserLatestConnLogic GetUserLatestConn err: %v", err)
		return &pb.GetUserLatestConnResp{}, err
	}
	if len(batchResp.UserLatestConns) == 0 {
		return &pb.GetUserLatestConnResp{}, nil
	}
	return batchResp.UserLatestConns[0], nil
}
