package groupmodel

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type (
	Group struct {
		// 群id 从10001开始 redis incr
		Id     string `bson:"_id" json:"id" gorm:"column:id;primary_key;type:char(32);not null"`
		Name   string `bson:"name" json:"name" gorm:"column:name;type:varchar(32);not null;index;default:''"`
		Avatar string `bson:"avatar" json:"avatar"` // 群头像
		// 群主
		Owner string `bson:"owner" json:"owner" gorm:"column:owner;type:char(32);not null;index"`
		// 所有管理员
		Managers xorm.SliceString `bson:"managers" json:"managers" gorm:"column:managers;type:json;"`
		// 创建时间
		CreateTime int64 `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint;not null;index"`
		// 解散时间
		DismissTime int64 `bson:"dismissTime" json:"dismissTime" gorm:"column:dismissTime;type:bigint;not null;index"`
		// 群描述
		Description string `bson:"description" json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`
		// 群成员人数
		MemberCount int `bson:"memberCount" json:"memberCount" gorm:"column:memberCount;type:int;not null;default:0;index"`

		// GroupSetting
		// 全体禁言开关
		AllMute bool `bson:"allMute" json:"allMute" gorm:"column:allMute;type:tinyint(1);not null;default:0;index"`
		// 发言频率限制
		SpeakLimit int32 `bson:"speakLimit,omitempty" json:"speakLimit" gorm:"column:speakLimit;type:int;not null;default:0;"`
		// 群成员人数上限
		MaxMember int32 `bson:"maxMember,omitempty" json:"maxMember" gorm:"column:maxMember;type:int;not null;default:0;"`
		// 成员权限选项
		// 群成员是否可以发起临时会话
		MemberCanStartTempChat bool `bson:"memberCanStartTempChat" json:"memberCanStartTempChat" gorm:"column:memberCanStartTempChat;type:tinyint(1);not null;default:0;"`
		// 群成员是否可以邀请好友加入群
		MemberCanInviteFriend bool `bson:"memberCanInviteFriend" json:"memberCanInviteFriend" gorm:"column:memberCanInviteFriend;type:tinyint(1);not null;default:0;"`
		// 新成员可见的历史消息条数
		NewMemberHistoryMsgCount int32 `bson:"newMemberHistoryMsgCount,omitempty" json:"newMemberHistoryMsgCount" gorm:"column:newMemberHistoryMsgCount;type:int;not null;default:0;"`
		// 是否开启匿名聊天
		AnonymousChat   bool            `bson:"anonymousChat" json:"anonymousChat" gorm:"column:anonymousChat;type:tinyint(1);not null;default:0;"`
		JoinGroupOption JoinGroupOption `bson:"joinGroupOption" json:"joinGroupOption" gorm:"column:joinGroupOption;type:json;"`
		AdminRemark     string          `bson:"adminRemark" json:"adminRemark" gorm:"column:adminRemark;type:varchar(255);not null;default:''"`
	}
	JoinGroupOption struct {
		Type int `bson:"type" json:"type"`
		// 验证信息
		// 问题
		Question string `bson:"question" json:"question"`
		// 答案
		Answer string `bson:"answer" json:"answer"`
	}
)

func (m *Group) TableName() string {
	return "group"
}

func (m *Group) GroupBaseInfo() *pb.GroupBaseInfo {
	return &pb.GroupBaseInfo{
		Id:     m.Id,
		Name:   m.Name,
		Avatar: m.Avatar,
	}
}

func (m *Group) Bytes() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *Group) ToPB() *pb.GroupModel {
	return &pb.GroupModel{
		Id:                       m.Id,
		Name:                     m.Name,
		Avatar:                   m.Avatar,
		Owner:                    m.Owner,
		Managers:                 m.Managers,
		CreateTime:               m.CreateTime,
		CreateTimeStr:            utils.TimeFormat(m.CreateTime),
		DismissTime:              m.DismissTime,
		DismissTimeStr:           utils.TimeFormat(m.DismissTime),
		Description:              m.Description,
		MemberCount:              int32(m.MemberCount),
		AllMute:                  m.AllMute,
		SpeakLimit:               m.SpeakLimit,
		MaxMember:                m.MaxMember,
		MemberCanStartTempChat:   m.MemberCanStartTempChat,
		MemberCanInviteFriend:    m.MemberCanInviteFriend,
		NewMemberHistoryMsgCount: m.NewMemberHistoryMsgCount,
		AnonymousChat:            m.AnonymousChat,
		JoinGroupOption:          m.JoinGroupOption.ToPB(),
		AdminRemark:              m.AdminRemark,
	}
}

func (m JoinGroupOption) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JoinGroupOption) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

func (m JoinGroupOption) ToPB() *pb.GroupModel_JoinGroupOption {
	return &pb.GroupModel_JoinGroupOption{
		Type:     int32(m.Type),
		Question: m.Question,
		Answer:   m.Answer,
	}
}

func GroupFromBytes(data []byte) *Group {
	group := &Group{}
	err := json.Unmarshal(data, group)
	if err != nil {
		return nil
	}
	return group
}

func ListGroupByIdsFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, ids []string) ([]*Group, error) {
	var groups []*Group
	err := tx.Where("id in (?)", ids).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	// 缓存到redis
	foundMap := make(map[string]bool)
	for _, group := range groups {
		err := rc.SetexCtx(ctx, rediskey.GroupKey(group.Id), string(group.Bytes()), rediskey.GroupKeyExpire())
		if err != nil {
			logx.Errorf("redis setex error: %v", err)
		}
		foundMap[group.Id] = true
	}
	// not found
	for _, id := range ids {
		if _, found := foundMap[id]; !found {
			// 存入占位符
			err := rc.SetexCtx(ctx, rediskey.GroupKey(id), xredis.NotFound, rediskey.GroupKeyExpire())
			if err != nil {
				logx.Errorf("redis setex error: %v", err)
			}
		}
	}
	return groups, nil
}

func ListGroupByIdsFromRedis(ctx context.Context, tx *gorm.DB, rc *redis.Redis, ids []string) ([]*Group, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// mget
	groups := make([]*Group, 0)
	keys := make([]string, 0)
	for _, id := range ids {
		keys = append(keys, rediskey.GroupKey(id))
	}
	val, err := rc.MgetCtx(ctx, keys...)
	if err != nil {
		return nil, err
	}
	foundMap := make(map[string]bool)
	realNotFoundMap := make(map[string]bool)
	for _, v := range val {
		// 是否为占位符
		if v == xredis.NotFound {
			// 真的不存在
			realNotFoundMap[v] = true
			continue
		}
		if v == "" {
			continue
		}
		// 反序列化
		group := GroupFromBytes([]byte(v))
		groups = append(groups, group)
		foundMap[group.Id] = true
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, found := foundMap[id]; !found {
			if _, realNotFound := realNotFoundMap[id]; !realNotFound {
				notFoundIds = append(notFoundIds, id)
			}
		}
	}
	// 从mysql中查询
	if len(notFoundIds) > 0 {
		mysqlGroups, err := ListGroupByIdsFromMysql(ctx, tx, rc, notFoundIds)
		if err != nil {
			return nil, err
		}
		for _, group := range mysqlGroups {
			groups = append(groups, group)
		}
	}
	// 返回
	return groups, nil
}

func CleanGroupCache(ctx context.Context, rc *redis.Redis, ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	keys := make([]string, 0)
	for _, id := range ids {
		keys = append(keys, rediskey.GroupKey(id))
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err
}
