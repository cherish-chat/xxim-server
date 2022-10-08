package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupLogic {
	return &GetGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroup 获取群组
func (l *GetGroupLogic) GetGroup(in *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupResp{}, nil
}
