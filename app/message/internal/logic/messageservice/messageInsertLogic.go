package messageservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInsertLogic {
	return &MessageInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MessageInsert 插入消息
func (l *MessageInsertLogic) MessageInsert(in *pb.MessageInsertReq) (*pb.MessageInsertResp, error) {
	// todo: add your logic here and delete this line

	return &pb.MessageInsertResp{}, nil
}
