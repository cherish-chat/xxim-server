package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditGroupInfoLogic {
	return &EditGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditGroupInfo 编辑群信息
func (l *EditGroupInfoLogic) EditGroupInfo(in *pb.EditGroupInfoReq) (*pb.EditGroupInfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.EditGroupInfoResp{}, nil
}
