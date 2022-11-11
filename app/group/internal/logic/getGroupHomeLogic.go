package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupHomeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupHomeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupHomeLogic {
	return &GetGroupHomeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupHome 获取群聊首页
func (l *GetGroupHomeLogic) GetGroupHome(in *pb.GetGroupHomeReq) (*pb.GetGroupHomeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupHomeResp{}, nil
}
