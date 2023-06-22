package infoservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserModelByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserModelByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserModelByIdLogic {
	return &GetUserModelByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserModelById 获取用户模型
func (l *GetUserModelByIdLogic) GetUserModelById(in *pb.GetUserModelByIdReq) (*pb.GetUserModelByIdResp, error) {
	userDest := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{
		"userId": in.UserId,
	}).One(userDest)
	if err != nil {
		l.Errorf("get user model by id error: %v", err)
		return &pb.GetUserModelByIdResp{}, err
	}
	userSettingsMap := make(map[string]*usermodel.UserSetting)
	userSettings := make([]*usermodel.UserSetting, 0)
	if in.GetOpt().GetWithUserSettings() && len(in.GetOpt().GetUserSettingKeys()) > 0 {
		keys := in.GetOpt().GetUserSettingKeys()
		err := l.svcCtx.UserSettingCollection.Find(l.ctx, bson.M{
			"userId": in.UserId,
			"k": bson.M{
				"$in": keys,
			},
		}).All(&userSettings)
		if err != nil {
			l.Errorf("get user settings error: %v", err)
			return &pb.GetUserModelByIdResp{}, err
		}
		for _, setting := range userSettings {
			userSettingsMap[setting.K] = setting
		}

		notExistsKeys := make([]string, 0)
		for _, key := range keys {
			if _, ok := userSettingsMap[key]; !ok {
				notExistsKeys = append(notExistsKeys, key)
			}
		}

		if len(notExistsKeys) > 0 {
			var notExistsUserSettings []*usermodel.UserSetting
			err := l.svcCtx.UserSettingCollection.Find(l.ctx, bson.M{
				"userId": "",
				"k": bson.M{
					"$in": notExistsKeys,
				},
			}).All(&notExistsUserSettings)
			if err != nil {
				l.Errorf("get user settings error: %v", err)
				return &pb.GetUserModelByIdResp{}, err
			}
			for _, setting := range notExistsUserSettings {
				userSettingsMap[setting.K] = setting
			}
		}

		notExistsKeys = make([]string, 0)
		for _, key := range keys {
			if _, ok := userSettingsMap[key]; !ok {
				notExistsKeys = append(notExistsKeys, key)
			}
		}

		if len(notExistsKeys) > 0 {
			l.Errorf("[WARNING] get user settings, not exists keys: %v", utils.AnyString(notExistsKeys))
			for _, key := range notExistsKeys {
				userSettingsMap[key] = &usermodel.UserSetting{
					K: key,
					V: "",
				}
			}
		}
	}
	return &pb.GetUserModelByIdResp{
		UserModelJson:    utils.Json.MarshalToBytes(userDest),
		UserSettingsJson: utils.Json.MarshalToBytes(userSettingsMap),
	}, nil
}
