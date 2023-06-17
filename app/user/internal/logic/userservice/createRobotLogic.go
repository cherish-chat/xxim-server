package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRobotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRobotLogic {
	return &CreateRobotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateRobot 创建机器人
func (l *CreateRobotLogic) CreateRobot(in *pb.CreateRobotReq) (*pb.CreateRobotResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateRobotResp{}, nil
}
