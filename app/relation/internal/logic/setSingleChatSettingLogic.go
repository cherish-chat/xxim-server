package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetSingleChatSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetSingleChatSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSingleChatSettingLogic {
	return &SetSingleChatSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetSingleChatSettingLogic) SetSingleChatSetting(in *pb.SetSingleChatSettingReq) (*pb.SetSingleChatSettingResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetSingleChatSettingResp{}, nil
}
