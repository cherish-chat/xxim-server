package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettingsLogic {
	return &GetUserSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserSettingsLogic) GetUserSettings(in *pb.GetUserSettingsReq) (*pb.GetUserSettingsResp, error) {
	var models []*usermodel.UserSetting
	var err error
	if len(in.Keys) > 0 {
		err = l.svcCtx.Mysql().Model(&usermodel.UserSetting{}).Where("userId = ? and `key` in (?)", in.CommonReq.UserId, in.Keys).Find(&models).Error
	} else {
		err = l.svcCtx.Mysql().Model(&usermodel.UserSetting{}).Where("userId = ?", in.CommonReq.UserId).Find(&models).Error
	}
	if err != nil {
		l.Errorf("get user settings failed, err: %v", err)
		return &pb.GetUserSettingsResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	var resp = make(map[int32]*pb.UserSetting)
	for _, model := range models {
		resp[int32(model.Key)] = &pb.UserSetting{
			Key:   model.Key,
			Value: model.Value,
		}
	}
	var notInitKeys []int32
	for _, key := range in.Keys {
		if _, ok := resp[int32(key)]; !ok {
			notInitKeys = append(notInitKeys, int32(key))
		}
	}
	if len(notInitKeys) > 0 {
		// 获取默认设置
		var defaultModels []*usermodel.UserSetting
		err = l.svcCtx.Mysql().Model(&usermodel.UserSetting{}).Where("userId = ? and `key` in (?)", "", notInitKeys).Find(&defaultModels).Error
		if err != nil {
			l.Errorf("get user settings failed, err: %v", err)
			return &pb.GetUserSettingsResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
		for _, model := range defaultModels {
			resp[int32(model.Key)] = &pb.UserSetting{
				Key:   model.Key,
				Value: model.Value,
			}
		}
	}
	for _, key := range in.Keys {
		if _, ok := resp[int32(key)]; !ok {
			resp[int32(key)] = &pb.UserSetting{
				Key:   key,
				Value: "",
			}
		}
	}
	return &pb.GetUserSettingsResp{Settings: resp}, nil
}
