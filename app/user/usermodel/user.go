package usermodel

import (
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User 用户 数据库模型
type User struct {
	UserId primitive.ObjectID `bson:"_id" json:"userId"`

	// 账户信息
	// Username 用户名 唯一
	Username string `bson:"username" json:"username"`
	// Password 密码
	Password string `bson:"password" json:"password"`
	// Salt 密码盐
	Salt string `bson:"salt" json:"salt"`
	// PhoneCountryCode 手机国家码
	PhoneCountryCode int32 `bson:"phoneCountryCode" json:"phoneCountryCode"`
	// Phone 手机号   // + 国家码 唯一
	Phone string `bson:"phone" json:"phone"`
	// DestroyTime 注销时间
	DestroyTime primitive.DateTime `bson:"destroyTime" json:"destroyTime"`
	// RegisterTime 注册时间
	RegisterTime primitive.DateTime `bson:"registerTime" json:"registerTime"`

	// 基本信息
	// Nickname 昵称
	Nickname string `bson:"nickname" json:"nickname"`
	// Avatar 头像
	Avatar string `bson:"avatar" json:"avatar"`
}

func (m *User) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"-registerTime"},
		IndexOptions: options.Index().SetName("registerTime"),
	}, {
		Key:          []string{"username"},
		IndexOptions: options.Index().SetName("username").SetUnique(true),
	}, {
		Key:          []string{"phone", "phoneCountryCode"},
		IndexOptions: options.Index().SetName("phone").SetUnique(true),
	}, {
		Key:          []string{"-nickname"},
		IndexOptions: options.Index().SetName("nickname"),
	}, {
		Key:          []string{"-destroyTime"},
		IndexOptions: options.Index().SetName("destroyTime"),
	}}
}

type xUserModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var UserModel *xUserModel

func InitUserModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	UserModel = &xUserModel{
		coll: coll,
		rc:   rc,
	}
}
