package usermodel

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type (
	User struct {
		Id           string `bson:"_id" json:"id" gorm:"column:id;primary_key;type:char(32);"`
		Password     string `bson:"password" json:"password" gorm:"column:password;type:char(64);"`
		PasswordSalt string `bson:"passwordSalt" json:"passwordSalt" gorm:"column:password_salt;type:char(64);"`
		Nickname     string `bson:"nickname" json:"nickname" gorm:"column:nickname;type:varchar(64);index;"`
		Avatar       string `bson:"avatar" json:"avatar" gorm:"column:avatar;type:varchar(255);"`
		// 注册信息
		RegInfo  *LoginInfo       `bson:"regInfo" json:"regInfo" gorm:"column:reg_info;type:json;"`
		Xb       pb.XB            `bson:"xb" json:"xb" gorm:"column:xb;type:tinyint(1);index;"`
		Birthday *pb.BirthdayInfo `bson:"birthday,omitempty" json:"birthday,omitempty" gorm:"column:birthday;type:json;"`
		// 其他信息
		InfoMap   xorm.M    `bson:"infoMap" json:"infoMap" gorm:"column:info_map;type:json;"`
		LevelInfo LevelInfo `bson:"levelInfo" json:"levelInfo" gorm:"column:level_info;type:json;"`
	}
	LoginInfo struct {
		// 13位时间戳
		Time int64  `bson:"time" json:"time" gorm:"column:time;type:bigint(13);index;"`
		Ip   string `bson:"ip" json:"ip" gorm:"column:ip;type:varchar(64);"`
		// 中国
		IpCountry string `bson:"ipCountry" json:"ipCountry" gorm:"column:ipCountry;type:varchar(64);"`
		// 北京市
		IpProvince string `bson:"ipProvince" json:"ipProvince" gorm:"column:ipProvince;type:varchar(64);"`
		// 北京市
		IpCity string `bson:"ipCity" json:"ipCity" gorm:"column:ipCity;type:varchar(64);"`
		// 电信
		IpISP string `bson:"ipService" json:"ipService" gorm:"column:ipService;type:varchar(64);"`
		// 1.0.0
		AppVersion string `bson:"appVersion" json:"appVersion" gorm:"column:appVersion;type:varchar(64);"`
		// user-agent
		Ua string `bson:"ua" json:"ua" gorm:"column:ua;type:varchar(255);"`
		// 10.0.0
		OsVersion string `bson:"osVersion" json:"osVersion" gorm:"column:osVersion;type:varchar(64);"`
		// iphone/ipad/android/pc/mac/linux/windows
		Platform string `bson:"platform" json:"platform" gorm:"column:platform;type:varchar(64);"`
		// 设备id
		DeviceId string `bson:"deviceId" json:"deviceId" gorm:"column:deviceId;type:varchar(64);"`
		// 设备型号
		DeviceModel string `bson:"deviceModel" json:"deviceModel" gorm:"column:deviceModel;type:varchar(64);"`
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

func (m *User) TableName() string {
	return "user"
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

func (m *User) GetId() string {
	return m.Id
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
		UserId       string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);primary_key"`
		Password     string `bson:"password" json:"password" gorm:"column:password;type:char(64)"`
		PasswordSalt string `bson:"passwordSalt" json:"passwordSalt" gorm:"column:passwordSalt;type:char(64)"`
		// 注册信息
		RegInfo *LoginInfo `bson:"regInfo" json:"regInfo" gorm:"column:regInfo;type:json"`
	}
)

func (m *UserTmp) TableName() string {
	return "user_tmp"
}

// GetUsersByIds 批量获取用户信息
func GetUsersByIds(ctx context.Context, rc *redis.Redis, tx *gorm.DB, ids []string) ([]*User, error) {
	users, err := getUsersByIdsFromRedis(ctx, rc, ids)
	if err != nil {
		return userFromMysql(ctx, rc, tx, ids)
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
		missUsers, err := userFromMysql(ctx, rc, tx, missIds)
		if err != nil {
			return nil, err
		}
		users = append(users, missUsers...)
	}
	return users, nil
}

func userFromMysql(ctx context.Context, rc *redis.Redis, tx *gorm.DB, ids []string) ([]*User, error) {
	users := make([]*User, 0)
	err := error(nil)
	xtrace.StartFuncSpan(ctx, "FindUserByIds", func(ctx context.Context) {
		err = tx.Where("id in (?)", ids).Find(&users).Error
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

func (m LoginInfo) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *LoginInfo) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

func (m LevelInfo) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *LevelInfo) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}
