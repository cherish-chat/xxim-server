package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupNoticeLogic {
	return &CreateGroupNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateGroupNotice 创建群公告
func (l *CreateGroupNoticeLogic) CreateGroupNotice(in *pb.CreateGroupNoticeReq) (*pb.CreateGroupNoticeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateGroupNoticeResp{}, nil
}
