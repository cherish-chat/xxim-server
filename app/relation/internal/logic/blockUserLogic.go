package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BlockUserLogic) BlockUser(in *pb.BlockUserReq) (*pb.BlockUserResp, error) {
	blacklist := &relationmodel.Blacklist{
		UserId:      in.Requester.Id,
		BlacklistId: in.UserId,
		CreateTime:  time.Now().UnixMilli(),
	}
	_, err := l.svcCtx.Mongo().Collection(blacklist).Upsert(l.ctx, bson.M{
		"userId":      blacklist.UserId,
		"blacklistId": blacklist.BlacklistId,
	}, blacklist)
	if err != nil {
		l.Errorf("Upsert failed, err: %v", err)
		return &pb.BlockUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 刷新缓存
	err = relationmodel.FlushBlacklistList(l.ctx, l.svcCtx.Redis(), in.Requester.Id)
	if err != nil {
		l.Errorf("FlushBlacklistList failed, err: %v", err)
		return &pb.BlockUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		_, _ = relationmodel.GetMyBlacklistList(ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(blacklist), in.Requester.Id)
	}, nil)
	return &pb.BlockUserResp{}, nil
}
