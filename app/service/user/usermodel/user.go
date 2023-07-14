package usermodel

import (
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//AccountStatus 账户状态

const (
	AccountStatusNormal = "0"
)

//AccountRole 账户角色

const (
	AccountRoleUser  = "user"
	AccountRoleRobot = "robot"
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
	AccountMap bson.M `bson:"accountMap,omitempty" json:"accountMap"`

	// 基本信息
	// Nickname 昵称
	Nickname string `bson:"nickname" json:"nickname"`
	// Avatar 头像
	Avatar string `bson:"avatar" json:"avatar"`
	// ProfileMap 个人资料
	ProfileMap bson.M `bson:"profileMap,omitempty" json:"profileMap"`
	// CountMap 计数信息
	CountMap bson.M `bson:"countMap,omitempty" json:"countMap"`

	// ExtraMap 扩展信息
	ExtraMap bson.M `bson:"extraMap,omitempty" json:"extraMap"`
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
		Key:          []string{"accountMap." + peerpb.AccountType_Username.String()},
		IndexOptions: options.Index().SetName("username"),
	}, {
		Key:          []string{"accountMap." + peerpb.AccountType_Email.String()},
		IndexOptions: options.Index().SetName("email"),
	}, {
		Key:          []string{"accountMap." + peerpb.AccountType_Phone.String(), "accountMap." + peerpb.AccountType_PhoneCountryCode.String()},
		IndexOptions: options.Index().SetName("phone"),
	}}
}

func (m *User) GetAccountMap() utils.SSM {
	return utils.NewSSMFromBsonM(m.AccountMap)
}

type UserCountMap struct {
	// FriendCount 好友数量
	FriendCount uint32
	// JoinGroupCount 加入群组数量
	JoinGroupCount uint32
	// CreateGroupCount 创建群组数量
	CreateGroupCount uint32
}

func (m *User) GetCountMap() UserCountMap {
	countMap := UserCountMap{}
	c, ok := m.CountMap[peerpb.UpdateUserCountMapReq_friendCount.String()]
	if ok {
		countMap.FriendCount = uint32(utils.Number.Any2Int64(c))
	}
	c, ok = m.CountMap[peerpb.UpdateUserCountMapReq_joinGroupCount.String()]
	if ok {
		countMap.JoinGroupCount = uint32(utils.Number.Any2Int64(c))
	}
	c, ok = m.CountMap[peerpb.UpdateUserCountMapReq_createGroupCount.String()]
	if ok {
		countMap.CreateGroupCount = uint32(utils.Number.Any2Int64(c))
	}
	return countMap
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
