package usermodel

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	User struct {
		Id           string `bson:"_id" json:"id"`
		Password     string `bson:"password" json:"password"`
		PasswordSalt string `bson:"passwordSalt" json:"passwordSalt"`
		Nickname     string `bson:"nickname" json:"nickname"`
		Avatar       string `bson:"avatar" json:"avatar"`
		// 注册信息
		RegInfo  *LoginInfo       `bson:"regInfo" json:"regInfo"`
		Xb       pb.XB            `bson:"xb" json:"xb"`
		Birthday *pb.BirthdayInfo `bson:"birthday,omitempty" json:"birthday,omitempty"`
		// 其他信息
		InfoMap   M         `bson:"infoMap" json:"infoMap"`
		LevelInfo LevelInfo `bson:"levelInfo" json:"levelInfo"`
	}
	LoginInfo struct {
		Time        int64  `bson:"time" json:"time"` // 13位时间戳
		Ip          string `bson:"ip" json:"ip"`
		IpCountry   string `bson:"ipCountry" json:"ipCountry"`     // 中国
		IpProvince  string `bson:"ipProvince" json:"ipProvince"`   // 北京市
		IpCity      string `bson:"ipCity" json:"ipCity"`           // 北京市
		IpISP       string `bson:"ipService" json:"ipService"`     // 电信
		AppVersion  string `bson:"appVersion" json:"appVersion"`   // 1.0.0
		Ua          string `bson:"ua" json:"ua"`                   // user-agent
		OsVersion   string `bson:"osVersion" json:"osVersion"`     // 10.0.0
		Platform    string `bson:"platform" json:"platform"`       // iphone/ipad/android/pc/mac/linux/windows
		DeviceId    string `bson:"deviceId" json:"deviceId"`       // 设备id
		DeviceModel string `bson:"deviceModel" json:"deviceModel"` // 设备型号
	}
	LevelInfo struct {
		Level        int32 `bson:"level" json:"level"`
		Exp          int32 `bson:"exp" json:"exp"`
		NextLevelExp int32 `bson:"nextLevelExp" json:"nextLevelExp"`
	}
)

func (m LevelInfo) Pb() *pb.LevelInfo {
	return &pb.LevelInfo{
		Level:        m.Level,
		Exp:          m.Exp,
		NextLevelExp: m.NextLevelExp,
	}
}

func (m *User) CollectionName() string {
	return "user"
}

func (m *User) Indexes(c *qmgo.Collection) error {
	//TODO implement me
	return nil
}

func (m *User) GetId() string {
	return m.Id
}

func (m *User) Marshal() []byte {
	return utils.AnyToBytes(m)
}

func (m *User) ExpireSeconds() int {
	return xredis.ExpireMinutes(5)
}

func (m *User) NotFound(id string) {
	m.Id = id
}

func (m *User) String() string {
	return utils.AnyToString(m)
}

func (m *User) BaseInfo() *pb.UserBaseInfo {
	return &pb.UserBaseInfo{
		Id:       m.Id,
		Nickname: m.Nickname,
		Avatar:   m.Avatar,
		Xb:       m.Xb,
		Birthday: m.Birthday,
		IpRegion: nil,
	}
}

func UserFromBytes(bytes []byte) *User {
	v := &User{}
	_ = json.Unmarshal(bytes, v)
	return v
}

type (
	UserTmp struct {
		UserId       string `bson:"userId" json:"userId"`
		Password     string `bson:"password" json:"password"`
		PasswordSalt string `bson:"passwordSalt" json:"passwordSalt"`
		// 注册信息
		RegInfo *LoginInfo `bson:"regInfo" json:"regInfo"`
	}
)

func (m *UserTmp) CollectionName() string {
	return "user_tmp"
}

func (m *UserTmp) Indexes(c *qmgo.Collection) error {
	//TODO implement me
	return nil
}

// GetUsersByIds 批量获取用户信息
func GetUsersByIds(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, ids []string) ([]*User, error) {
	users, err := getUsersByIdsFromRedis(ctx, rc, ids)
	if err != nil {
		return userFromMongo(ctx, rc, c, ids)
	}
	// 判断是否有缺失
	userMap := make(map[string]*User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	var missIds []string
	for _, id := range ids {
		if _, ok := userMap[id]; !ok {
			missIds = append(missIds, id)
		}
	}
	if len(missIds) > 0 {
		missUsers, err := userFromMongo(ctx, rc, c, missIds)
		if err != nil {
			return nil, err
		}
		users = append(users, missUsers...)
	}
	return users, nil
}

func userFromMongo(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, ids []string) ([]*User, error) {
	users := make([]*User, 0)
	err := error(nil)
	xtrace.StartFuncSpan(ctx, "FindUserByIds", func(ctx context.Context) {
		err = c.Find(ctx, bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		}).All(&users)
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("find users by ids %v error: %s", ids, err.Error())
		return users, err
	}
	// 存入 redis
	userMap := make(map[string]*User)
	for _, user := range users {
		userMap[user.Id] = user
		key := rediskey.UserKey(user.Id)
		err = rc.SetexCtx(ctx, key, utils.AnyToString(user), user.ExpireSeconds())
		if err != nil {
			logx.WithContext(ctx).Errorf("set user %s to redis error: %s", user.Id, err.Error())
			continue
		}
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := userMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		for _, id := range notFoundIds {
			key := rediskey.UserKey(id)
			err = rc.SetexCtx(ctx, key, xredis.NotFound, xredis.ExpireMinutes(5))
			if err != nil {
				logx.WithContext(ctx).Errorf("set user %s to redis error: %s", id, err.Error())
				continue
			}
		}
	}
	return users, nil
}

func getUsersByIdsFromRedis(ctx context.Context, rc *redis.Redis, ids []string) ([]*User, error) {
	users := make([]*User, 0)
	vals, err := rc.MgetCtx(ctx, utils.UpdateSlice(ids, func(id string) string {
		return rediskey.UserKey(id)
	})...)
	if err != nil {
		logx.WithContext(ctx).Errorf("get users by ids %v from redis error: %s", ids, err.Error())
		return users, err
	}
	for i, val := range vals {
		user := &User{}
		if val == xredis.NotFound {
			id := ids[i]
			user.NotFound(id)
		} else {
			err = json.Unmarshal([]byte(val), user)
			if err != nil {
				logx.WithContext(ctx).Errorf("convert user error: %s", err.Error())
				continue
			}
		}
		users = append(users, user)
	}
	return users, nil
}

func FlushUserCache(ctx context.Context, rc *redis.Redis, ids []string) error {
	var err error
	if len(ids) > 0 {
		xtrace.StartFuncSpan(ctx, "DeleteCache", func(ctx context.Context) {
			redisKeys := utils.UpdateSlice(ids, func(v string) string {
				return rediskey.UserKey(v)
			})
			_, err = rc.DelCtx(ctx, redisKeys...)
		})
	}
	return err
}
