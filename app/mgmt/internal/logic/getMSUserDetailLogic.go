package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSUserDetailLogic {
	return &GetMSUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSUserDetailLogic) GetMSUserDetail(in *pb.GetMSUserDetailReq) (*pb.GetMSUserDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMSUserDetailResp{}, nil
}
