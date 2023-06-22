package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountCreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountCreateGroupLogic {
	return &CountCreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CountCreateGroup 统计创建的群组数量
func (l *CountCreateGroupLogic) CountCreateGroup(in *pb.CountCreateGroupReq) (*pb.CountCreateGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CountCreateGroupResp{}, nil
}
