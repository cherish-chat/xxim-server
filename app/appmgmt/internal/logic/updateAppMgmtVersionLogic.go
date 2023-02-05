package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtVersionLogic {
	return &UpdateAppMgmtVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtVersionLogic) UpdateAppMgmtVersion(in *pb.UpdateAppMgmtVersionReq) (*pb.UpdateAppMgmtVersionResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Version{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtVersion.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtVersionResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["version"] = in.AppMgmtVersion.Version
		updateMap["platform"] = in.AppMgmtVersion.Platform
		updateMap["type"] = in.AppMgmtVersion.Type
		updateMap["content"] = in.AppMgmtVersion.Content
		updateMap["content"] = in.AppMgmtVersion.Content
		updateMap["downloadUrl"] = in.AppMgmtVersion.DownloadUrl
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtVersion.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtVersionResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtVersionResp{}, nil
}
