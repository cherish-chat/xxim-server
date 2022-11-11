package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetSingleMsgNotifyOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetSingleMsgNotifyOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSingleMsgNotifyOptLogic {
	return &SetSingleMsgNotifyOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetSingleMsgNotifyOptLogic) SetSingleMsgNotifyOpt(in *pb.SetSingleMsgNotifyOptReq) (*pb.SetSingleMsgNotifyOptResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetSingleMsgNotifyOptResp{}, nil
}
