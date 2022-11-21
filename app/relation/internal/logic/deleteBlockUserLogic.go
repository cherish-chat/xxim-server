package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBlockUserLogic {
	return &DeleteBlockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBlockUserLogic) DeleteBlockUser(in *pb.DeleteBlockUserReq) (*pb.DeleteBlockUserResp, error) {
	blacklist := &relationmodel.Blacklist{
		UserId:      in.CommonReq.Id,
		BlacklistId: in.UserId,
	}
	err := l.svcCtx.Mysql().Where("userId = ? and blacklistId = ?", blacklist.UserId, blacklist.BlacklistId).Delete(blacklist).Error
	if err != nil {
		l.Errorf("Upsert failed, err: %v", err)
		return &pb.DeleteBlockUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 刷新缓存
	err = relationmodel.FlushBlacklistList(l.ctx, l.svcCtx.Redis(), in.CommonReq.Id)
	if err != nil {
		l.Errorf("FlushBlacklistList failed, err: %v", err)
		return &pb.DeleteBlockUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		_, _ = relationmodel.GetMyBlacklistList(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.CommonReq.Id)
	}, nil)
	return &pb.DeleteBlockUserResp{}, nil
}
