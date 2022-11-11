package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupSettingLogic {
	return &SetGroupSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetGroupSetting 设置群设置
func (l *SetGroupSettingLogic) SetGroupSetting(in *pb.SetGroupSettingReq) (*pb.SetGroupSettingResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetGroupSettingResp{}, nil
}
