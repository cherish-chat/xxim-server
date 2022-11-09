package usermodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type UserSetting struct {
	UserId string            `json:"userId" bson:"userId"` // 如果userId为空，则表示所有用户的默认设置
	Key    pb.UserSettingKey `json:"key" bson:"key"`
	Value  string            `json:"value" bson:"value"`
}

func (m *UserSetting) CollectionName() string {
	return "user_setting"
}

func (m *UserSetting) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key: []string{"userId"},
	}, {
		Key:          []string{"userId", "key"},
		IndexOptions: opts.Index().SetUnique(true),
	}})
	return nil
}

func InitUserSetting(c *qmgo.Collection) {
	defaultUserSetting := func(k pb.UserSettingKey, v string) *UserSetting {
		return &UserSetting{
			UserId: "",
			Key:    k,
			Value:  v,
		}
	}
	ctx := context.Background()
	c.InsertOne(ctx, defaultUserSetting(pb.UserSettingKey_HowToAddFriend, "need_confirm"))
	c.InsertOne(ctx, defaultUserSetting(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Question, "你的名字是？"))
	c.InsertOne(ctx, defaultUserSetting(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Answer, "xxim"))
}
