package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

//RichArticle 富文本文章
/*
	type AppMgmtRichArticle struct {
		// 文章id
		Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
		// 文章标题
		Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title"`
		// 富文本内容
		Content string `protobuf:"bytes,3,opt,name=content,proto3" json:"content"`
		// 内容类型
		ContentType string `protobuf:"bytes,4,opt,name=contentType,proto3" json:"contentType"` // example: text/html; text/markdown; text/plain; application/json
		// url地址
		Url string `protobuf:"bytes,5,opt,name=url,proto3" json:"url"`
		// 是否启用
		IsEnable bool `protobuf:"varint,6,opt,name=isEnable,proto3" json:"isEnable"`
		// 创建时间
		CreatedAt int64 `protobuf:"varint,7,opt,name=createdAt,proto3" json:"createdAt"`
		// 创建时间字符串
		CreatedAtStr string `protobuf:"bytes,8,opt,name=createdAtStr,proto3" json:"createdAtStr"`
		// 更新时间
		UpdatedAt int64 `protobuf:"varint,9,opt,name=updatedAt,proto3" json:"updatedAt"`
		// 更新时间字符串
		UpdatedAtStr string `protobuf:"bytes,10,opt,name=updatedAtStr,proto3" json:"updatedAtStr"`
		// 排序
		Sort int32 `protobuf:"varint,11,opt,name=sort,proto3" json:"sort"`
	}
*/
type RichArticle struct {
	// 文章id
	Id string `gorm:"column:id;type:char(32);primary_key;not null" json:"id"`
	// 文章标题
	Title string `gorm:"column:title;type:varchar(255);not null;index;" json:"title"`
	// 富文本内容
	Content string `gorm:"column:content;type:longtext;not null;" json:"content"`
	// 内容类型
	ContentType string `gorm:"column:contentType;type:varchar(255);not null;" json:"contentType"` // example: text/html; text/markdown; text/plain; application/json
	// url地址
	Url string `gorm:"column:url;type:varchar(255);not null;" json:"url"`
	// 是否启用
	IsEnable bool `gorm:"column:isEnable;type:tinyint(1);not null;default:1;index;" json:"isEnable"`
	// 创建时间
	CreatedAt int64 `gorm:"column:createdAt;type:bigint(20);not null;index;" json:"createdAt"`
	// 更新时间
	UpdatedAt int64 `gorm:"column:updatedAt;type:bigint(20);not null;index;" json:"updatedAt"`
	// 排序
	Sort int32 `gorm:"column:sort;type:int(11);not null;default:0;index;" json:"sort"`
}

// TableName sets the insert table name for this struct type
func (a *RichArticle) TableName() string {
	return APPMGR_TABLE_PREFIX + "rich_article"
}

func (a *RichArticle) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *RichArticle) ToPB() *pb.AppMgmtRichArticle {
	return &pb.AppMgmtRichArticle{
		Id:           a.Id,
		Title:        a.Title,
		Content:      a.Content,
		ContentType:  a.ContentType,
		Url:          a.Url,
		IsEnable:     a.IsEnable,
		CreatedAt:    a.CreatedAt,
		CreatedAtStr: utils.TimeFormat(a.CreatedAt),
		UpdatedAt:    a.UpdatedAt,
		UpdatedAtStr: utils.TimeFormat(a.UpdatedAt),
		Sort:         a.Sort,
	}
}
