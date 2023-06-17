package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetUserAccountMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserAccountMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserAccountMapLogic {
	return &ResetUserAccountMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResetUserAccountMap 重置用户账号信息
func (l *ResetUserAccountMapLogic) ResetUserAccountMap(in *pb.ResetUserAccountMapReq) (*pb.ResetUserAccountMapResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ResetUserAccountMapResp{}, nil
}
