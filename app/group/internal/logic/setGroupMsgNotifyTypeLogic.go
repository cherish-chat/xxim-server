package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupMsgNotifyTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupMsgNotifyTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMsgNotifyTypeLogic {
	return &SetGroupMsgNotifyTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetGroupMsgNotifyType 设置群消息通知选项
func (l *SetGroupMsgNotifyTypeLogic) SetGroupMsgNotifyType(in *pb.SetGroupMsgNotifyTypeReq) (*pb.SetGroupMsgNotifyTypeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetGroupMsgNotifyTypeResp{}, nil
}
