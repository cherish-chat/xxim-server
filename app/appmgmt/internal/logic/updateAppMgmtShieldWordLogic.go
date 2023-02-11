package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtShieldWordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtShieldWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtShieldWordLogic {
	return &UpdateAppMgmtShieldWordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtShieldWordLogic) UpdateAppMgmtShieldWord(in *pb.UpdateAppMgmtShieldWordReq) (*pb.UpdateAppMgmtShieldWordResp, error) {
	// 查询原模型
	model := &appmgmtmodel.ShieldWord{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtShieldWord.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtShieldWordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["word"] = in.AppMgmtShieldWord.Word
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtShieldWord.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtShieldWordResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtShieldWordResp{}, nil
}
