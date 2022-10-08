package dbmodel

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/qiniu/qmgo"
)

type (
	User struct {
		Id           string       `bson:"_id" json:"id"`
		Nickname     string       `bson:"nickname" json:"nickname"`
		Avatar       string       `bson:"avatar" json:"avatar"`
		Xb           string       `bson:"xb" json:"xb"`
		Birthday     string       `bson:"birthday" json:"birthday"`
		Signature    string       `bson:"signature" json:"signature"`
		Tags         []string     `bson:"tags" json:"tags"`
		Password     string       `bson:"password" json:"password"`
		IsRobot      bool         `bson:"isRobot" json:"isRobot"`
		IsGuest      bool         `bson:"isGuest" json:"isGuest"`
		IsAdmin      bool         `bson:"isAdmin" json:"isAdmin"`
		IsOfficial   bool         `bson:"isOfficial" json:"isOfficial"`
		UnbanTime    string       `bson:"unbanTime" json:"unbanTime"`
		AdminRemark  string       `bson:"adminRemark" json:"adminRemark"`
		RegistryInfo RegistryInfo `bson:"registryInfo" json:"registryInfo"`
		Ex           UserEx       `bson:"ex" json:"ex"`
	}
	UserEx       struct{}
	RegistryInfo struct {
		Platform       string   `bson:"platform" json:"platform"`
		Time           string   `bson:"time" json:"time"`
		DeviceModel    string   `bson:"deviceModel" json:"deviceModel"`
		DeviceId       string   `bson:"deviceId" json:"deviceId"`
		Ips            []string `bson:"ips" json:"ips"`
		RegisterSource string   `bson:"registerSource" json:"registerSource"`
		Salt           string   `bson:"salt" json:"salt"`
		IpCountry      string   `bson:"ipCountry" json:"ipCountry"`
		IpProvince     string   `bson:"ipProvince" json:"ipProvince"`
		IpCity         string   `bson:"ipCity" json:"ipCity"`
		IpDistrict     string   `bson:"ipDistrict" json:"ipDistrict"`
	}
)

func (e *UserEx) Bytes() []byte {
	bytes, _ := json.Marshal(e)
	return bytes
}

func NewUserEx(bytes []byte) UserEx {
	ex := UserEx{}
	_ = json.Unmarshal(bytes, &ex)
	return ex
}

func InitUser(c *qmgo.Collection) {

}

func (x *User) ToPbUser() *pb.UserData {
	return &pb.UserData{
		Id:        x.Id,
		Nickname:  x.Nickname,
		Avatar:    x.Avatar,
		Xb:        x.Xb,
		Birthday:  x.Birthday,
		Signature: x.Signature,
		Tags:      x.Tags,
		Password:  x.Password,
		RegisterInfo: &pb.UserData_RegisterInfo{
			Platform:       x.RegistryInfo.Platform,
			Time:           x.RegistryInfo.Time,
			Ip:             utils.Conditional(len(x.RegistryInfo.Ips) > 0, x.RegistryInfo.Ips[0], ""),
			Device:         x.RegistryInfo.DeviceModel,
			DeviceId:       x.RegistryInfo.DeviceId,
			RegisterSource: x.RegistryInfo.RegisterSource,
			Salt:           x.RegistryInfo.Salt,
			IpCountry:      x.RegistryInfo.IpCountry,
			IpProvince:     x.RegistryInfo.IpProvince,
			IpCity:         x.RegistryInfo.IpCity,
			IpDistrict:     x.RegistryInfo.IpDistrict,
		},
		IsRobot:     x.IsRobot,
		IsGuest:     x.IsGuest,
		IsAdmin:     x.IsAdmin,
		IsOfficial:  x.IsOfficial,
		UnbanTime:   x.UnbanTime,
		AdminRemark: x.AdminRemark,
		Ex:          x.Ex.Bytes(),
	}
}
