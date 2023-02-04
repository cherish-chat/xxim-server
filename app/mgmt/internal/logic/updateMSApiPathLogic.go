package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSApiPathLogic {
	return &UpdateMSApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSApiPathLogic) UpdateMSApiPath(in *pb.UpdateMSApiPathReq) (*pb.UpdateMSApiPathResp, error) {
	// 查询原模型
	model := &mgmtmodel.ApiPath{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.ApiPath.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSApiPathResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.ApiPath.Title != "" {
		updateMap["title"] = in.ApiPath.Title
	}
	if in.ApiPath.Path != "" {
		updateMap["path"] = in.ApiPath.Path
	}
	updateMap["logEnable"] = in.ApiPath.LogEnable
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.ApiPath.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSApiPathResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSApiPathResp{}, nil
}
