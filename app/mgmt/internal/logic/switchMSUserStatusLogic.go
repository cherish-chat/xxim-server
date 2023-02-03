package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchMSUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSwitchMSUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchMSUserStatusLogic {
	return &SwitchMSUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SwitchMSUserStatusLogic) SwitchMSUserStatus(in *pb.SwitchMSUserStatusReq) (*pb.SwitchMSUserStatusResp, error) {
	// 查询原模型
	model := &mgmtmodel.User{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.SwitchMSUserStatusResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	updateMap["isDisable"] = !model.IsDisable
	err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).Updates(updateMap).Error
	if err != nil {
		l.Errorf("更新用户失败: %v", err)
		return &pb.SwitchMSUserStatusResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.SwitchMSUserStatusResp{}, nil
}
