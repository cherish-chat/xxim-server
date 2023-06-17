package worldservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/world/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorldPostSubmitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWorldPostSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorldPostSubmitLogic {
	return &WorldPostSubmitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// WorldPostSubmit 世界圈帖子发布
func (l *WorldPostSubmitLogic) WorldPostSubmit(in *pb.WorldPostSubmitReq) (*pb.WorldPostSubmitResp, error) {
	// todo: add your logic here and delete this line

	return &pb.WorldPostSubmitResp{}, nil
}
