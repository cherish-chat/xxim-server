package usermodel

type InvitationCode struct {
	Id string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
}
