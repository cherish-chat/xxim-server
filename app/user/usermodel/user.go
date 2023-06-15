package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User 用户 数据库模型
type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	// 账户信息
	// Username 用户名 唯一
	UserId string `bson:"userId" json:"userId"`
	// RegisterTime 注册时间
	RegisterTime primitive.DateTime `bson:"registerTime" json:"registerTime"`
	// DestroyTime 注销时间
	DestroyTime primitive.DateTime `bson:"destroyTime,omitempty" json:"destroyTime,omitempty"`
	// AccountMap 账户object
	AccountMap bson.M `bson:"accountMap" json:"accountMap"`

	// 基本信息
	// Nickname 昵称
	Nickname string `bson:"nickname" json:"nickname"`
	// Avatar 头像
	Avatar string `bson:"avatar" json:"avatar"`
	// ProfileMap 个人资料
	ProfileMap bson.M `bson:"profileMap" json:"profileMap"`

	// ExtraMap 扩展信息
	ExtraMap bson.M `bson:"extraMap" json:"extraMap"`
}

func (m *User) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"-registerTime"},
		IndexOptions: options.Index().SetName("registerTime"),
	}, {
		Key:          []string{"userId"},
		IndexOptions: options.Index().SetName("userId").SetUnique(true),
	}, {
		Key:          []string{"-nickname"},
		IndexOptions: options.Index().SetName("nickname"),
	}, {
		Key:          []string{"-destroyTime"},
		IndexOptions: options.Index().SetName("destroyTime"),
	}, {
		Key:          []string{"accountMap." + pb.AccountTypeUsername},
		IndexOptions: options.Index().SetName("username"),
	}, {
		Key:          []string{"accountMap." + pb.AccountTypeEmail},
		IndexOptions: options.Index().SetName("email"),
	}, {
		Key:          []string{"accountMap." + pb.AccountTypePhone, "accountMap." + pb.AccountTypePhoneCode},
		IndexOptions: options.Index().SetName("phone"),
	}}
}

func (m *User) GetAccountMap() utils.SSM {
	return utils.NewSSMFromBsonM(m.AccountMap)
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
