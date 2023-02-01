package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSUserListLogic {
	return &GetAllMSUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSUserListLogic) GetAllMSUserList(in *pb.GetAllMSUserListReq) (*pb.GetAllMSUserListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAllMSUserListResp{}, nil
}
