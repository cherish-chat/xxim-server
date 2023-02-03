package mgmtmodel

type LoginRecordInfo LoginInfo
type LoginRecord struct {
	Id     string `bson:"_id" json:"id" gorm:"column:id;type:char(32);primary_key"`
	UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);index"`
	LoginRecordInfo
}

func (m *LoginRecord) TableName() string {
	return MGMT_TABLE_PREFIX + "login_record"
}
