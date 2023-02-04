package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSIpWhiteListLogic {
	return &UpdateMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSIpWhiteListLogic) UpdateMSIpWhiteList(in *pb.UpdateMSIpWhiteListReq) (*pb.UpdateMSIpWhiteListResp, error) {
	// 查询原模型
	model := &mgmtmodel.MSIPWhitelist{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.IpWhiteList.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSIpWhiteListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.IpWhiteList.StartIp != "" {
		updateMap["startIp"] = in.IpWhiteList.StartIp
	}
	if in.IpWhiteList.EndIp != "" {
		updateMap["endIp"] = in.IpWhiteList.EndIp
	}
	if in.IpWhiteList.Remark != "" {
		updateMap["remark"] = in.IpWhiteList.Remark
	}
	updateMap["isEnable"] = in.IpWhiteList.IsEnable
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.IpWhiteList.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSIpWhiteListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSIpWhiteListResp{}, nil
}
