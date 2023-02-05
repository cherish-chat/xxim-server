package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtEmojiGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtEmojiGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtEmojiGroupLogic {
	return &UpdateAppMgmtEmojiGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtEmojiGroupLogic) UpdateAppMgmtEmojiGroup(in *pb.UpdateAppMgmtEmojiGroupReq) (*pb.UpdateAppMgmtEmojiGroupResp, error) {
	// 查询原模型
	model := &appmgmtmodel.EmojiGroup{}
	err := l.svcCtx.Mysql().Model(model).Where("name = ?", in.AppMgmtEmojiGroup.Name).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtEmojiGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["name"] = in.AppMgmtEmojiGroup.Name
		updateMap["coverId"] = in.AppMgmtEmojiGroup.CoverId
		updateMap["isEnable"] = in.AppMgmtEmojiGroup.IsEnable
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("name = ?", in.AppMgmtEmojiGroup.Name).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtEmojiGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtEmojiGroupResp{}, nil
}
