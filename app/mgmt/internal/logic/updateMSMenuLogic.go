package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSMenuLogic {
	return &UpdateMSMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSMenuLogic) UpdateMSMenu(in *pb.UpdateMSMenuReq) (*pb.UpdateMSMenuResp, error) {
	// 查询原模型
	model := &mgmtmodel.Menu{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Menu.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSMenuResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	updateMap["updateTime"] = time.Now().UnixMilli()
	if in.Menu.Pid != "" {
		updateMap["pid"] = in.Menu.Pid
	}
	if in.Menu.MenuType != "" {
		updateMap["menuType"] = in.Menu.MenuType
	}
	if in.Menu.MenuName != "" {
		updateMap["menuName"] = in.Menu.MenuName
	}
	if in.Menu.MenuIcon != "" {
		updateMap["menuIcon"] = in.Menu.MenuIcon
	}
	updateMap["menuSort"] = in.Menu.MenuSort
	updateMap["perms"] = in.Menu.Perms
	updateMap["paths"] = in.Menu.Paths
	updateMap["component"] = in.Menu.Component
	updateMap["selected"] = in.Menu.Selected
	updateMap["params"] = in.Menu.Params
	updateMap["isCache"] = in.Menu.IsCache
	updateMap["isShow"] = in.Menu.IsShow
	updateMap["isDisable"] = in.Menu.IsDisable
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.Menu.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSMenuResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSMenuResp{}, nil
}
