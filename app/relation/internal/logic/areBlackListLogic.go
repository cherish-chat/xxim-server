package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"

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
	blacklist, err := relationmodel.AreMyBlacklist(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&relationmodel.Friend{}), in.A, in.BList)
	if err != nil {
		l.Errorf("AreMyFriend failed, err: %v", err)
		return &pb.AreBlackListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AreBlackListResp{BlackList: blacklist}, nil
}
