package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MsgInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMsgInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MsgInsertLogic {
	return &MsgInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MsgInsert 插入消息
func (l *MsgInsertLogic) MsgInsert(in *pb.MsgInsertReq) (*pb.MsgInsertResp, error) {
	// todo: add your logic here and delete this line

	return &pb.MsgInsertResp{}, nil
}
