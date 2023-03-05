package appmgmtmodel

type Stats struct {
	Date string `gorm:"column:date;type:varchar(10);index:date_type_idx,unique;not null" json:"date"`
	Type string `gorm:"column:type;type:varchar(32);index:date_type_idx,unique;not null" json:"type"`
	Val  int    `gorm:"column:val;type:bigint(20);not null" json:"val"`
}

func (m *Stats) TableName() string {
	return APPMGR_TABLE_PREFIX + "stats"
}

const (
	StatsTypeNewUser          = "newUser"
	StatsTypeNewGroup         = "newGroup"
	StatsTypeTodayMsg         = "todayMsg"
	StatsTypeTodayMsgUser     = "todayMsgUser"
	StatsTypeTodayAliveGroup  = "todayAliveGroup"
	StatsTypeTodayAliveSingle = "todayAliveSingle"
	StatsTypeTodayAliveUser   = "todayAliveUser"
	StatsTypeTodayNewFriend   = "todayNewFriend"
)

var StatsNameMap = map[string]string{
	StatsTypeNewUser:          "新增用户",
	StatsTypeNewGroup:         "新增群组",
	StatsTypeTodayMsg:         "今日消息",
	StatsTypeTodayMsgUser:     "发消息用户",
	StatsTypeTodayAliveGroup:  "活跃群组",
	StatsTypeTodayAliveSingle: "活跃单聊",
	StatsTypeTodayAliveUser:   "活跃用户",
	StatsTypeTodayNewFriend:   "新增好友",
}
