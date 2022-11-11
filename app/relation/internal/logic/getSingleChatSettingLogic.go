package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSingleChatSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSingleChatSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSingleChatSettingLogic {
	return &GetSingleChatSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSingleChatSettingLogic) GetSingleChatSetting(in *pb.GetSingleChatSettingReq) (*pb.GetSingleChatSettingResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetSingleChatSettingResp{}, nil
}
