package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelfUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfUserInfoLogic {
	return &GetSelfUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetSelfUserInfo 获取自己的用户信息
func (l *GetSelfUserInfoLogic) GetSelfUserInfo(in *peerpb.GetSelfUserInfoReq) (*peerpb.GetSelfUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GetSelfUserInfoResp{}, nil
}
