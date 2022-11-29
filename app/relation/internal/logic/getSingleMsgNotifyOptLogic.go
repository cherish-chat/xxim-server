package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSingleMsgNotifyOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSingleMsgNotifyOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSingleMsgNotifyOptLogic {
	return &GetSingleMsgNotifyOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSingleMsgNotifyOptLogic) GetSingleMsgNotifyOpt(in *pb.GetSingleMsgNotifyOptReq) (*pb.GetSingleMsgNotifyOptResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetSingleMsgNotifyOptResp{
		CommonResp: pb.NewSuccessResp(),
		Opt: &pb.MsgNotifyOpt{
			NoDisturb: false,
			Preview:   false,
			Sound:     false,
			SoundName: "",
			Vibrate:   false,
		},
	}, nil
}
