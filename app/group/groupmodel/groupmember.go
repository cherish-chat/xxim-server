package groupmodel

type (
	GroupMember struct {
		// 群id
		GroupId string `bson:"groupId" json:"groupId" gorm:"column:groupId;type:char(32);not null;index:group_user,unique;comment:群id;index;"`
		// 用户id
		UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);not null;index:group_user,unique;comment:用户id;index;"`
		// 加入时间
		CreateTime int64 `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint;not null;index;comment:加入时间"`
	}
)

func (m *GroupMember) TableName() string {
	return "group_member"
}
