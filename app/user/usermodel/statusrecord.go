package usermodel

type StatusRecord struct {
	Id string `gorm:"column:id;type:varchar(32);primary_key" json:"id"`
	// 用户id
	UserId string `gorm:"column:userId;type:varchar(32);not null;default:''" json:"userId"`
	// 谁封禁的
	Operator string `gorm:"column:operator;type:varchar(32);not null;default:''" json:"operator"`
	// 封禁时间
	DisableTime int64 `gorm:"column:disableTime;type:bigint(20);not null;default:0" json:"disableTime"`
	// 解封时间
	UnblockTime int64 `gorm:"column:unblockTime;type:bigint(20);not null;default:0" json:"unblockTime"`
	// 取消时间
	CancelTime int64 `gorm:"column:cancelTime;type:bigint(20);not null;default:0" json:"cancelTime"`
	// 是否封禁ip
	DisableIp string `gorm:"column:disableIp;type:varchar(32);not null;default:''" json:"disableIp"`
}

func (m *StatusRecord) TableName() string {
	return TABLE_PREFIX + "status_record"
}
