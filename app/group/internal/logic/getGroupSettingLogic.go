package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupSettingLogic {
	return &GetGroupSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupSetting 获取群设置
func (l *GetGroupSettingLogic) GetGroupSetting(in *pb.GetGroupSettingReq) (*pb.GetGroupSettingResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupSettingResp{}, nil
}
