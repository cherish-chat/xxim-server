package groupmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"os"
)

// Group 群组 数据库模型
type Group struct {
	//GroupId 群id 主键
	GroupId int `json:"groupId" gorm:"column:groupId;type:bigint(20);not null;primaryKey" bson:"_id"`
	//GroupName 群名称
	GroupName string `json:"groupName" gorm:"column:groupName;type:varchar(32);not null;index;" bson:"groupName"`
	//GroupAvatar 群头像
	GroupAvatar string `json:"groupAvatar" gorm:"column:groupAvatar;type:varchar(255);not null" bson:"groupAvatar"`
	//GroupInfo 自定义群信息
	//如果没有查询需要，不要SELECT这个字段
	GroupInfo bson.M `json:"groupInfo" gorm:"column:groupInfo;type:text;not null;" bson:"groupInfo"`

	//OwnerUserId 群主id
	OwnerUserId string `json:"ownerUserId" gorm:"column:ownerUserId;type:char(32);not null" bson:"ownerUserId"`
	//ManagerUserIds 管理员id列表 逗号分隔
	//如果没有查询需要，不要SELECT这个字段，因为这个字段可能会很大，群管理员上限是1900人，因为65535/33=1985.90
	ManagerUserIds []string `json:"managerUserIds" gorm:"column:managerUserIds;type:text;not null" bson:"managerUserIds"`
	//CreatedAt 创建时间 13位时间戳
	CreateTime int64 `json:"createTime" gorm:"column:createTime;type:bigint;not null;index;" bson:"createTime"`
	//UpdatedAt 更新时间 13位时间戳
	UpdateTime int64 `json:"updateTime" gorm:"column:updateTime;type:bigint;not null" bson:"updateTime"`
	//DismissTime 解散时间 13位时间戳
	DismissTime int64 `json:"dismissTime" gorm:"column:dismissTime;type:bigint;not null;default:0;index;" bson:"dismissTime"`
	//MemberCount 成员数量
	MemberCount int `json:"memberCount" gorm:"column:memberCount;type:int;not null;default:0;" bson:"memberCount"`

	//RemarkForAdmin 管理员设置的备注
	RemarkForAdmin string `json:"remarkForAdmin" gorm:"column:remarkForAdmin;type:varchar(32);not null;default:'';" bson:"remarkForAdmin"`
}

// TableName 表名
func (m *Group) TableName() string {
	return "group"
}

// GetIndexes 获取索引
func (m *Group) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key: []string{"groupName"},
	}, {
		Key: []string{"-dismissTime"},
	}, {
		Key: []string{"-createTime"},
	}}
}

type xGroupModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var GroupModel *xGroupModel

func InitGroupModel(coll *qmgo.QmgoClient, rc *redis.Redis, minGroupId int) {
	GroupModel = &xGroupModel{
		coll: coll,
		rc:   rc,
	}
	// 查询最大的群组id
	group := &Group{}
	err := coll.Find(context.Background(), bson.M{}).Sort("-_id").Limit(1).One(group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 说明没有数据
			group.GroupId = minGroupId
		} else {
			logx.Errorf(`db.Model(&Group{}).Order("groupId desc").First(&group) error: %v`, err)
			os.Exit(1)
			return
		}
	}
	// 保存到redis
	err = rc.Set(xcache.IncrKeyGroupId, utils.AnyString(group.GroupId))
	if err != nil {
		logx.Errorf(`rc.Set(xcache.IncrKeyGroupId, utils.AnyString(group.GroupId)) error: %v`, err)
		os.Exit(1)
		return
	}
}

func (x *xGroupModel) GenerateGroupId() int {
	// 从redis中获取
	groupId, err := x.rc.Incr(xcache.IncrKeyGroupId)
	if err != nil {
		logx.Errorf("x.rc.Incr(xcache.IncrKeyGroupId) error: %v", err)
		return 0
	}
	return int(groupId)
}

func (x *xGroupModel) Insert(ctx context.Context, tx gorm.DB, group *Group) error {
	if group.GroupId == 0 {
		group.GroupId = x.GenerateGroupId()
	}
	return tx.WithContext(ctx).Create(group).Error
}
