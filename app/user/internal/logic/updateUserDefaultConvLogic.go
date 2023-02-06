package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserDefaultConvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserDefaultConvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserDefaultConvLogic {
	return &UpdateUserDefaultConvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserDefaultConvLogic) UpdateUserDefaultConv(in *pb.UpdateUserDefaultConvReq) (*pb.UpdateUserDefaultConvResp, error) {
	// 查询原模型
	model := &usermodel.DefaultConv{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserDefaultConv.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateUserDefaultConvResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["convType"] = in.UserDefaultConv.ConvType
		updateMap["filterType"] = in.UserDefaultConv.FilterType
		updateMap["invitationCode"] = in.UserDefaultConv.InvitationCode
		updateMap["convId"] = in.UserDefaultConv.ConvId
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserDefaultConv.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateUserDefaultConvResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateUserDefaultConvResp{}, nil
}
