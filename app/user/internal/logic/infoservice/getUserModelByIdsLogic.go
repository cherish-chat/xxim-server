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

type GetUserModelByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserModelByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserModelByIdsLogic {
	return &GetUserModelByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserModelByIds 批量获取用户模型
func (l *GetUserModelByIdsLogic) GetUserModelByIds(in *pb.GetUserModelByIdsReq) (*pb.GetUserModelByIdsResp, error) {
	userDests := make([]*usermodel.User, 0)
	err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{
		"userId": bson.M{
			"$in": in.UserIds,
		},
	}).All(&userDests)
	if err != nil {
		l.Errorf("get user model by id error: %v", err)
		return &pb.GetUserModelByIdsResp{}, err
	}
	userSettingsMap := make(map[string]map[string]*usermodel.UserSetting)
	userSettings := make([]*usermodel.UserSetting, 0)
	if in.GetOpt().GetWithUserSettings() && len(in.GetOpt().GetUserSettingKeys()) > 0 {
		keys := in.GetOpt().GetUserSettingKeys()
		err := l.svcCtx.UserSettingCollection.Find(l.ctx, bson.M{
			"userId": bson.M{
				"$in": in.UserIds,
			},
			"k": bson.M{
				"$in": keys,
			},
		}).All(&userSettings)
		if err != nil {
			l.Errorf("get user settings error: %v", err)
			return &pb.GetUserModelByIdsResp{}, err
		}
		for _, setting := range userSettings {
			userId := setting.UserId
			if _, ok := userSettingsMap[userId]; !ok {
				userSettingsMap[userId] = make(map[string]*usermodel.UserSetting)
			}
			userSettingsMap[userId][setting.K] = setting
		}

		notExistsKeys := make([]string, 0)
		for _, userid := range in.UserIds {
			if userSettingMap, ok := userSettingsMap[userid]; !ok {
				notExistsKeys = append(notExistsKeys, keys...)
				userSettingsMap[userid] = make(map[string]*usermodel.UserSetting)
			} else {
				for _, key := range keys {
					if _, ok := userSettingMap[key]; !ok {
						notExistsKeys = append(notExistsKeys, key)
					}
				}
			}
		}

		notExistsKeys = utils.AnySet(notExistsKeys)

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
				return &pb.GetUserModelByIdsResp{}, err
			}
			for _, setting := range notExistsUserSettings {
				for userId, m := range userSettingsMap {
					if _, ok := m[setting.K]; !ok {
						userSettingsMap[userId][setting.K] = setting
					}
				}
			}
		}

		for _, userid := range in.UserIds {
			userSettingMap := userSettingsMap[userid]
			for _, key := range keys {
				if _, ok := userSettingMap[key]; !ok {
					userSettingsMap[userid][key] = &usermodel.UserSetting{
						K: key,
						V: "",
					}
				}
			}
		}
	}

	userDestsResp := make(map[string][]byte)
	for _, dest := range userDests {
		userDestsResp[dest.UserId] = utils.Json.MarshalToBytes(dest)
	}
	userSettingsResp := make(map[string][]byte)
	for userId, settingMap := range userSettingsMap {
		userSettingsResp[userId] = utils.Json.MarshalToBytes(settingMap)
	}
	return &pb.GetUserModelByIdsResp{
		UserModelJsons:    userDestsResp,
		UserSettingsJsons: userSettingsResp,
	}, nil
}
