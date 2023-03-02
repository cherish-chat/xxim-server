package groupmodel

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"math"
)

type (
	RoleType    int8 // 角色类型 // 0:普通成员 1:管理员 2:群主
	GroupMember struct {
		// 群id
		GroupId string `bson:"groupId" json:"groupId" gorm:"column:groupId;type:char(32);not null;index:group_user,unique;comment:群id;index;"`
		// 用户id
		UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);not null;index:group_user,unique;comment:用户id;index;"`
		// 加入时间
		CreateTime int64 `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint;not null;index;comment:加入时间"`
		// 角色
		Role RoleType `bson:"role" json:"role" gorm:"column:role;type:int;not null;default:0;comment:角色;index;"` // 0:普通成员 1:管理员 2:群主
		// 我设置的我的备注
		Remark string `bson:"remark" json:"remark" gorm:"column:remark;type:varchar(255);not null;default:'';comment:我设置的我的备注;index;"`
		// 我设置的群的备注
		GroupRemark string `bson:"groupRemark" json:"groupRemark" gorm:"column:groupRemark;type:varchar(255);not null;default:'';comment:我设置的群的备注"`
		// 解禁时间
		UnbanTime int64 `bson:"unbanTime" json:"unbanTime" gorm:"column:unbanTime;type:bigint;not null;default:0;comment:解禁时间"`
		// 是否免打扰
		NoDisturb bool `bson:"noDisturb" json:"noDisturb" gorm:"column:noDisturb;type:tinyint(1);not null;index;default:0;comment:是否免打扰"`
		// 预览
		Preview bool `bson:"preview" json:"preview" gorm:"column:preview;type:tinyint(1);not null;default:0;comment:预览"`
		// 置顶
		IsTop bool `bson:"isTop" json:"isTop" gorm:"column:isTop;type:tinyint(1);not null;default:0;comment:置顶"`
	}
)

const (
	RoleType_MEMBER  RoleType = 0 // 普通成员
	RoleType_MANAGER RoleType = 1 // 管理员
	RoleType_OWNER   RoleType = 2 // 群主
)

func (m *GroupMember) TableName() string {
	return "group_member"
}

func (m *GroupMember) Bytes() []byte {
	buf, _ := json.Marshal(m)
	return buf
}

func (m *GroupMember) Pb() *pb.GroupMemberInfo {
	return &pb.GroupMemberInfo{
		GroupId:     m.GroupId,
		MemberId:    m.UserId,
		Remark:      m.Remark,
		GroupRemark: m.GroupRemark,
		Top:         false,
		NoDisturb:   m.NoDisturb,
		Preview:     m.Preview,
		DisturbMore: 0,
		ChatBg:      "",
		Role:        pb.GroupRole(m.Role),
		UnbanTime:   m.UnbanTime,
	}
}

func GroupMemberFromBytes(bytes []byte) *GroupMember {
	var groupMember GroupMember
	_ = json.Unmarshal(bytes, &groupMember)
	return &groupMember
}

func ListGroupsByUserIdFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, userId string) ([]string, error) {
	var groupMembers []*GroupMember
	err := tx.WithContext(ctx).Model(&GroupMember{}).Where("userId = ?", userId).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	// 放到redis zset中
	var groupIds []string
	var members []interface{}
	for _, group := range groupMembers {
		members = append(members, group.CreateTime, group.GroupId)
		groupIds = append(groupIds, group.GroupId)
	}
	err = xredis.ZAddsEx(rc, ctx, rediskey.GroupMemberListByUserId(userId), rediskey.GroupMemberListByUserIdExpire(), members)
	if err != nil {
		logx.Errorf("zadd group member list error: %v", err)
	}
	return groupIds, nil
}

func ListGroupsByUserIdFromRedis(ctx context.Context, tx *gorm.DB, rc *redis.Redis, userId string) ([]string, error) {
	// redis key 是否存在
	exists, err := rc.ExistsCtx(ctx, rediskey.GroupMemberListByUserId(userId))
	if err != nil {
		logx.Errorf("redis key exists error: %v", err)
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	if !exists {
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	// 从redis中获取
	val, err := rc.ZrangebyscoreWithScoresCtx(ctx, rediskey.GroupMemberListByUserId(userId), 0, math.MaxInt64)
	if err != nil {
		logx.Errorf("zrangebyscore error: %v", err)
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	var groupIds []string
	for _, v := range val {
		groupIds = append(groupIds, v.Key)
	}
	return groupIds, nil
}

func FlushGroupsByUserIdCache(ctx context.Context, rc *redis.Redis, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}
	var keys []string
	for _, userId := range userIds {
		keys = append(keys, rediskey.GroupMemberListByUserId(userId))
	}
	var err error
	xtrace.RunWithTrace(xtrace.TraceIdFromContext(ctx), "FlushGroupsByUserIdCache", func(ctx context.Context) {
		_, err = rc.DelCtx(ctx, keys...)
	}, nil)
	return err
}

func IsGroupMember(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId, userId string) (bool, error) {
	// redis key 是否存在
	exists, err := rc.ExistsCtx(ctx, rediskey.GroupMemberListByUserId(userId))
	if err != nil {
		logx.Errorf("redis key exists error: %v", err)
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	if !exists {
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	// 从redis中获取
	val, err := rc.ZscoreCtx(ctx, rediskey.GroupMemberListByUserId(userId), groupId)
	if err != nil {
		logx.Errorf("zscore error: %v", err)
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	return val > 0, nil
}

func IsGroupMemberFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId string, userId string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&GroupMember{}).Where("groupId = ? and userId = ?", groupId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ListGroupMemberFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId string, userIds []string) ([]*GroupMember, error) {
	var groupMembers []*GroupMember
	err := tx.WithContext(ctx).Model(&GroupMember{}).Where("groupId = ? and userId in (?)", groupId, userIds).
		Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	// 放到redis string中
	foundMap := make(map[string]bool)
	for _, groupMember := range groupMembers {
		err = rc.SetexCtx(ctx, rediskey.GroupMemberKey(groupId, groupMember.UserId), string(groupMember.Bytes()), rediskey.GroupMemberExpire())
		if err != nil {
			logx.Errorf("set group member error: %v", err)
		}
		foundMap[groupMember.UserId] = true
	}
	// 未找到的放到redis中
	for _, userId := range userIds {
		if _, found := foundMap[userId]; !found {
			err = rc.SetexCtx(ctx, rediskey.GroupMemberKey(groupId, userId), xredis.NotFound, rediskey.GroupMemberExpire())
			if err != nil {
				logx.Errorf("set group member error: %v", err)
			}
		}
	}
	return groupMembers, nil
}

func ListGroupMemberFromRedis(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId string, ids []string) ([]*GroupMember, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// mget
	models := make([]*GroupMember, 0)
	keys := make([]string, 0)
	for _, id := range ids {
		keys = append(keys, rediskey.GroupMemberKey(groupId, id))
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
		model := GroupMemberFromBytes([]byte(v))
		models = append(models, model)
		foundMap[model.UserId] = true
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
		mysqlGroups, err := ListGroupMemberFromMysql(ctx, tx, rc, groupId, notFoundIds)
		if err != nil {
			return nil, err
		}
		for _, group := range mysqlGroups {
			models = append(models, group)
		}
	}
	// 返回
	return models, nil
}

func FlushGroupMemberCache(ctx context.Context, rc *redis.Redis, groupId string, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}
	var keys []string
	for _, userId := range userIds {
		keys = append(keys, rediskey.GroupMemberKey(groupId, userId))
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err
}

func FlushGroupMemberListCache(ctx context.Context, rc *redis.Redis, groupId string) error {
	// 查询 keys
	keyListKey := rediskey.GroupMemberSearchKeyList(groupId)
	var keys = []string{keyListKey}
	// hgetall keyListKey
	val, err := rc.HgetallCtx(ctx, keyListKey)
	if err != nil {
		return err
	}
	for _, v := range val {
		keys = append(keys, v)
	}
	// del keys
	_, err = rc.DelCtx(ctx, keys...)
	return err
}
