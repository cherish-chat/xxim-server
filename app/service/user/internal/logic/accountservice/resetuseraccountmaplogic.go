package accountservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

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
func (l *ResetUserAccountMapLogic) ResetUserAccountMap(in *peerpb.ResetUserAccountMapReq) (*peerpb.ResetUserAccountMapResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.ResetUserAccountMapResp{}, nil
}
