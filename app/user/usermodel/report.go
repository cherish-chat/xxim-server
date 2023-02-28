package usermodel

import "github.com/cherish-chat/xxim-server/common/xorm"

// ReportRecord 举报记录
type ReportRecord struct {
	Id string `gorm:"column:id;primary_key;type:char(32);"`
	// 举报人
	ReporterId string `gorm:"column:reporter_id;type:char(32);not null;index:idx_reporter_id;"`
	// 被举报人
	ReportedId string `gorm:"column:reported_id;type:char(32);not null;index:idx_reported_id;"`
	// 举报类型
	ReportType string `gorm:"column:report_type;type:char(32);not null;index:idx_report_type;"`
	// 举报内容
	ReportContent string `gorm:"column:report_content;type:varchar(255);not null;"`
	// 举报图片
	ReportImages xorm.SliceString `gorm:"column:report_images;type:JSON;not null;"`
	// 举报时间
	ReportTime int64 `gorm:"column:report_time;type:bigint(20);not null;index:idx_report_time;"`
	// 举报状态
	ReportStatus string `gorm:"column:report_status;type:char(32);not null;index:idx_report_status;"`
	// 处理时间
	HandleTime int64 `gorm:"column:handle_time;type:bigint(20);not null;"`
	// 处理人
	HandlerId string `gorm:"column:handler_id;type:char(32);not null;"`
}

func (m *ReportRecord) TableName() string {
	return TABLE_PREFIX + "report_record"
}
