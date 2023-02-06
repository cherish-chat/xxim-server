package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserIpBlackListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserIpBlackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserIpBlackListLogic {
	return &UpdateUserIpBlackListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserIpBlackListLogic) UpdateUserIpBlackList(in *pb.UpdateUserIpBlackListReq) (*pb.UpdateUserIpBlackListResp, error) {
	// 查询原模型
	model := &usermodel.IpBlackList{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserIpList.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateUserIpBlackListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["platform"] = in.UserIpList.Platform
		updateMap["startIp"] = in.UserIpList.StartIp
		updateMap["endIp"] = in.UserIpList.EndIp
		updateMap["remark"] = in.UserIpList.Remark
		updateMap["userId"] = in.UserIpList.UserId
		updateMap["isEnable"] = in.UserIpList.IsEnable
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserIpList.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateUserIpBlackListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateUserIpBlackListResp{}, nil
}
