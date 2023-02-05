package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtVpnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtVpnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtVpnLogic {
	return &UpdateAppMgmtVpnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtVpnLogic) UpdateAppMgmtVpn(in *pb.UpdateAppMgmtVpnReq) (*pb.UpdateAppMgmtVpnResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Vpn{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtVpn.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtVpnResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["platform"] = in.AppMgmtVpn.Platform
		updateMap["type"] = in.AppMgmtVpn.Type
		updateMap["name"] = in.AppMgmtVpn.Name
		updateMap["ip"] = in.AppMgmtVpn.Ip
		updateMap["port"] = in.AppMgmtVpn.Port
		updateMap["username"] = in.AppMgmtVpn.Username
		updateMap["password"] = in.AppMgmtVpn.Password
		updateMap["secretKey"] = in.AppMgmtVpn.SecretKey
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtVpn.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtVpnResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtVpnResp{}, nil
}
