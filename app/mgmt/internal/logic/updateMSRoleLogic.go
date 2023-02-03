package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSRoleLogic {
	return &UpdateMSRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSRoleLogic) UpdateMSRole(in *pb.UpdateMSRoleReq) (*pb.UpdateMSRoleResp, error) {
	// 查询原模型
	model := &mgmtmodel.Role{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Role.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSRoleResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.Role.Name != "" {
		updateMap["name"] = in.Role.Name
	}
	if in.Role.Remark != "" {
		updateMap["remark"] = in.Role.Remark
	}
	updateMap["sort"] = in.Role.Sort
	updateMap["menuIds"] = strings.Join(in.Role.MenuIds, ",")
	updateMap["apiPathIds"] = strings.Join(in.Role.ApiPathIds, ",")
	if in.Role.IsDisable {
		if !model.IsDisable {
			// 封禁
			updateMap["isDisable"] = true
		}
	} else {
		if model.IsDisable {
			// 解封
			updateMap["isDisable"] = false
		}
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.Role.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSRoleResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSRoleResp{}, nil
}
