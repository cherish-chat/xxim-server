package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Album 相册实体
type Album struct {
	ID         uint   `gorm:"primarykey;comment:'主键ID';AUTO_INCREMENT;"`
	Cid        uint   `gorm:"not null;default:0;comment:'类目ID';index;"`
	Aid        string `gorm:"not null;default:0;comment:'管理ID';index;"`
	Type       int    `gorm:"not null;default:10;comment:'文件类型: [10=图片, 20=视频]';index;"`
	Name       string `gorm:"not null;default:'';comment:'文件名称';index;"`
	Url        string `gorm:"not null;comment:'文件路径'"`
	Ext        string `gorm:"not null;default:'';comment:'文件扩展';index;"`
	Size       int64  `gorm:"not null;default:0;comment:文件大小"`
	CreateTime int64  `gorm:"column:createTime;autoCreateTime;not null;comment:'创建时间';index;"`
	UpdateTime int64  `gorm:"column:updateTime;autoUpdateTime;not null;comment:'更新时间';index;"`
	DeleteTime int64  `gorm:"column:deleteTime;not null;default:0;comment:'删除时间';index;"`
}

func (m *Album) TableName() string {
	return MGMT_TABLE_PREFIX + "album"
}

func (m *Album) ToPB() *pb.MSAlbum {
	return &pb.MSAlbum{
		Id:            int32(m.ID),
		Cid:           strconv.Itoa(int(m.Cid)),
		Aid:           m.Aid,
		Type:          int32(m.Type),
		Name:          m.Name,
		Url:           m.Url,
		Ext:           m.Ext,
		Size:          m.Size,
		CreateTime:    m.CreateTime,
		UpdateTime:    m.UpdateTime,
		DeleteTime:    m.DeleteTime,
		CreateTimeStr: utils.TimeFormat(m.CreateTime),
		UpdateTimeStr: utils.TimeFormat(m.UpdateTime),
		DeleteTimeStr: utils.TimeFormat(m.DeleteTime),
	}
}

// AlbumCate 相册分类实体
type AlbumCate struct {
	ID         uint   `gorm:"primarykey;comment:'主键ID';AUTO_INCREMENT;"`
	Pid        uint   `gorm:"not null;default:0;comment:'父级ID';index;"`
	Type       int    `gorm:"not null;default:10;comment:'文件类型: [10=图片, 20=视频]';index;"`
	Name       string `gorm:"not null;default:'';comment:'分类名称';index;"`
	CreateTime int64  `gorm:"column:createTime;autoCreateTime;not null;comment:'创建时间';index;"`
	UpdateTime int64  `gorm:"column:updateTime;autoUpdateTime;not null;comment:'更新时间';"`
	DeleteTime int64  `gorm:"column:deleteTime;not null;default:0;comment:'删除时间';index;"`
}

func (m *AlbumCate) TableName() string {
	return MGMT_TABLE_PREFIX + "album_cate"
}

func (m *AlbumCate) ToPb() *pb.MSAlbumCate {
	return &pb.MSAlbumCate{
		Id:         int32(m.ID),
		Pid:        strconv.Itoa(int(m.Pid)),
		Type:       int32(m.Type),
		Name:       m.Name,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
		DeleteTime: m.DeleteTime,
	}
}

func initAlbumCate(tx *gorm.DB) {
	// name=默认相册
	insertIfNotFound(tx, "1", &AlbumCate{
		ID:         1,
		Pid:        0,
		Type:       10,
		Name:       "未命名相册",
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
		DeleteTime: 0,
	})
}
