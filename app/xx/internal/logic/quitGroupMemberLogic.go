package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuitGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuitGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitGroupMemberLogic {
	return &QuitGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuitGroupMember 退出群组
func (l *QuitGroupMemberLogic) QuitGroupMember(in *pb.QuitGroupMemberReq) (*pb.QuitGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QuitGroupMemberResp{}, nil
}
