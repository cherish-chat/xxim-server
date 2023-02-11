package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtNoticeLogic {
	return &UpdateAppMgmtNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtNoticeLogic) UpdateAppMgmtNotice(in *pb.UpdateAppMgmtNoticeReq) (*pb.UpdateAppMgmtNoticeResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Notice{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtNotice.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtNoticeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["position"] = in.AppMgmtNotice.Position
		updateMap["platform"] = in.AppMgmtNotice.Platform
		updateMap["title"] = in.AppMgmtNotice.Title
		updateMap["image"] = in.AppMgmtNotice.Image
		updateMap["content"] = in.AppMgmtNotice.Content
		updateMap["sort"] = in.AppMgmtNotice.Sort
		updateMap["isEnable"] = in.AppMgmtNotice.IsEnable
		updateMap["startTime"] = in.AppMgmtNotice.StartTime
		updateMap["endTime"] = in.AppMgmtNotice.EndTime
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtNotice.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtNoticeResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtNoticeResp{}, nil
}
