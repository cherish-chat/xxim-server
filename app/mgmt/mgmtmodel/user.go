package mgmtmodel

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type User struct {
	Id           string `bson:"_id" json:"id" gorm:"column:id;primary_key;type:char(32);"`
	Password     string `bson:"password" json:"password" gorm:"column:password;type:char(64);"`
	PasswordSalt string `bson:"passwordSalt" json:"passwordSalt" gorm:"column:passwordSalt;type:char(64);"`
	Nickname     string `bson:"nickname" json:"nickname" gorm:"column:nickname;type:varchar(64);index;"`
	Avatar       string `bson:"avatar" json:"avatar" gorm:"column:avatar;type:varchar(255);"`
	// 角色id
	RoleId     string     `bson:"roleId" json:"roleId" gorm:"column:roleId;type:char(32);index;default:'';"`
	IsDisable  bool       `bson:"isDisable" json:"isDisable" gorm:"column:isDisable;type:tinyint(1);default:0;"`
	CreateTime int64      `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint(13);index;"`
	RegInfo    *LoginInfo `bson:"regInfo" json:"regInfo" gorm:"column:regInfo;type:json;"`
}

type LoginInfo struct {
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
	// user-agent
	UserAgent string `bson:"userAgent" json:"userAgent" gorm:"column:userAgent;type:varchar(255);"`
}

func (m LoginInfo) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *LoginInfo) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

func (m *User) TableName() string {
	return MGMT_TABLE_PREFIX + "user"
}

func (m *User) ToPBMSUser(record *LoginRecord) *pb.MSUser {
	return &pb.MSUser{
		Id:            m.Id,
		CreatedAt:     m.CreateTime,
		CreatedAtStr:  utils.TimeFormat(m.CreateTime),
		Username:      m.Nickname,
		Password:      m.Password,
		Nickname:      m.Nickname,
		Avatar:        m.Avatar,
		Role:          m.RoleId,
		IsDisable:     m.IsDisable,
		LastLoginIp:   record.Ip,
		LastLoginTime: utils.TimeFormat(record.Time),
	}
}
