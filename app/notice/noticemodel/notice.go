package noticemodel

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type (
	Notice struct {
		NoticeId    string       `gorm:"column:noticeId;type:char(96);primary_key;not null" json:"noticeId"`
		ConvId      string       `gorm:"column:convId;type:char(96);not null" json:"convId"`
		CreateTime  int64        `gorm:"column:createTime;type:bigint(13);not null" json:"createTime"`
		Title       string       `gorm:"column:title;type:varchar(255);not null" json:"title"`
		ContentType int32        `gorm:"column:contentType;type:int(11);not null" json:"contentType"`
		Content     []byte       `gorm:"column:content;type:blob;not null" json:"content"`
		Options     NoticeOption `gorm:"column:options;type:json;" json:"options"`
		Ext         []byte       `gorm:"column:ext;type:blob;" json:"ext"`
		UserId      string       `gorm:"column:userId;type:char(32);not null;index;default:'';" json:"userId"`
		IsBroadcast bool         `gorm:"column:isBroadcast;type:tinyint(1);not null;default:0" json:"isBroadcast"`
		// 推送失效
		PushInvalid bool `gorm:"column:pushInvalid;type:tinyint(1);not null;default:0;index;" json:"pushInvalid"`
	}
	NoticeOption struct {
		StorageForClient bool `gorm:"column:storageForClient;type:tinyint(1);not null" json:"storageForClient"`
		UpdateConvMsg    bool `gorm:"column:updateConvMsg;type:tinyint(1);not null" json:"updateConvMsg"`
		OnlinePushOnce   bool `gorm:"column:onlinePushOnce;type:tinyint(1);not null" json:"onlinePushOnce"`
	}
)

func (m *Notice) TableName() string {
	return "notice"
}

func (m NoticeOption) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *NoticeOption) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

func (m *Notice) Upsert(tx *gorm.DB) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "noticeId"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       m.Title,
			"contentType": m.ContentType,
			"content":     m.Content,
			"options":     m.Options,
			"ext":         m.Ext,
			"createTime":  m.CreateTime,
		}),
	}).Create(m).Error
}

func (m *Notice) ToProto() *pb.NoticeData {
	return &pb.NoticeData{
		ConvId:      m.ConvId,
		UnreadCount: 0,
		NoticeId:    m.NoticeId,
		CreateTime:  utils.AnyToString(m.CreateTime),
		Title:       m.Title,
		ContentType: m.ContentType,
		Content:     m.Content,
		Options: &pb.NoticeData_Options{
			StorageForClient: m.Options.StorageForClient,
			UpdateConvMsg:    m.Options.UpdateConvMsg,
			OnlinePushOnce:   m.Options.OnlinePushOnce,
		},
		Ext: m.Ext,
	}
}

func NoticeFromPB(data *pb.NoticeData, isBroadcast bool, userId string) *Notice {
	return &Notice{
		NoticeId:    utils.If(data.NoticeId != "", data.NoticeId, utils.GenId()),
		ConvId:      data.ConvId,
		CreateTime:  utils.If(data.CreateTime != "", utils.AnyToInt64(data.CreateTime), time.Now().UnixMilli()),
		Title:       data.Title,
		ContentType: data.ContentType,
		Content:     data.Content,
		Options: NoticeOption{
			StorageForClient: data.Options.StorageForClient,
			UpdateConvMsg:    data.Options.UpdateConvMsg,
			OnlinePushOnce:   data.Options.OnlinePushOnce,
		},
		IsBroadcast: isBroadcast,
		Ext:         data.Ext,
		UserId:      userId,
	}
}
