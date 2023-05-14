package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostLogic {
	return &PostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostLogic) Post(in *pb.GatewayPostReq) (*pb.GatewayPostResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayPostResp{}, nil
}
