package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConvOnlineCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConvOnlineCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConvOnlineCountLogic {
	return &GetConvOnlineCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetConvOnlineCount 获取一个会话里所有的在线用户
func (l *GetConvOnlineCountLogic) GetConvOnlineCount(in *pb.GetConvOnlineCountReq) (*pb.GetConvOnlineCountResp, error) {
	//l.svcCtx.ImService().BatchGetUserLatestConn(l.ctx, &pb.BatchGetUserLatestConnReq{UserIds: memberIds})
	// todo
	return &pb.GetConvOnlineCountResp{}, nil
}
