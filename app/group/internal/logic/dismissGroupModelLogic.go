package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DismissGroupModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDismissGroupModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupModelLogic {
	return &DismissGroupModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DismissGroupModel 解散群组
func (l *DismissGroupModelLogic) DismissGroupModel(in *pb.DismissGroupModelReq) (*pb.DismissGroupModelResp, error) {
	var e error
	logic := NewKickGroupMemberLogic(l.ctx, l.svcCtx)
	for _, id := range in.Ids {
		xtrace.StartFuncSpan(l.ctx, "KickGroupMember.DismissRecoverGroup", func(ctx context.Context) {
			_, err := logic.DismissRecoverGroup(&pb.KickGroupMemberReq{
				CommonReq: in.CommonReq,
				GroupId:   id,
			})
			if err != nil {
				l.Errorf("DismissGroupModel KickGroupMember error: %v", err)
				e = err
			}
		})
	}
	if e != nil {
		return &pb.DismissGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, e
	}
	return &pb.DismissGroupModelResp{}, nil
}
