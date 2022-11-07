package usermodel

import (
	"github.com/qiniu/qmgo"
)

type (
	User struct {
		Id           string `bson:"_id" json:"id"`
		Password     string `bson:"password" json:"password"`
		PasswordSalt string `bson:"passwordSalt" json:"passwordSalt"`
		Nickname     string `bson:"nickname" json:"nickname"`
		Avatar       string `bson:"avatar" json:"avatar"`
		// 注册信息
		RegInfo *LoginInfo `bson:"regInfo" json:"regInfo"`
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
)

func (m *User) CollectionName() string {
	return "user"
}

func (m *User) Indexes(c *qmgo.Collection) error {
	//TODO implement me
	return nil
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
