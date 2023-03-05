package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtLinkLogic {
	return &UpdateAppMgmtLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtLinkLogic) UpdateAppMgmtLink(in *pb.UpdateAppMgmtLinkReq) (*pb.UpdateAppMgmtLinkResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Link{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtLink.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtLinkResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["sort"] = in.AppMgmtLink.Sort
		updateMap["url"] = in.AppMgmtLink.Url
		updateMap["name"] = in.AppMgmtLink.Name
		updateMap["icon"] = in.AppMgmtLink.Icon
		updateMap["isEnable"] = in.AppMgmtLink.IsEnable
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtLink.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtLinkResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtLinkResp{}, nil
}
