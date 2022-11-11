package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferGroupOwnerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferGroupOwnerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferGroupOwnerLogic {
	return &TransferGroupOwnerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TransferGroupOwner 转让群主
func (l *TransferGroupOwnerLogic) TransferGroupOwner(in *pb.TransferGroupOwnerReq) (*pb.TransferGroupOwnerResp, error) {
	// todo: add your logic here and delete this line

	return &pb.TransferGroupOwnerResp{}, nil
}
