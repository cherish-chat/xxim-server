package groupservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountJoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountJoinGroupLogic {
	return &CountJoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CountJoinGroup 统计加入的群组数量
func (l *CountJoinGroupLogic) CountJoinGroup(in *pb.CountJoinGroupReq) (*pb.CountJoinGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CountJoinGroupResp{}, nil
}
