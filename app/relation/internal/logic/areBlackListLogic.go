package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreBlackListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreBlackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreBlackListLogic {
	return &AreBlackListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreBlackListLogic) AreBlackList(in *pb.AreBlackListReq) (*pb.AreBlackListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AreBlackListResp{}, nil
}
