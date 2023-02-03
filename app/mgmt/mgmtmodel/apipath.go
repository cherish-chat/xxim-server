package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type ApiPath struct {
	Id         string `gorm:"column:id;primarykey;comment:'主键'"`
	Title      string `gorm:"column:title;not null;default:'';comment:'api标题''"`
	Path       string `gorm:"column:path;not null;default:'';comment:'api路径'"`
	Remark     string `gorm:"column:remark;not null;default:'';comment:'备注信息'"`
	CreateTime int64  `gorm:"column:createTime;not null;comment:'创建时间'"`
	UpdateTime int64  `gorm:"column:updateTime;not null;comment:'更新时间'"`
}

func (m *ApiPath) TableName() string {
	return MGMT_TABLE_PREFIX + "apipath"
}

func (m *ApiPath) ToPB() *pb.MSApiPath {
	return &pb.MSApiPath{
		Id:           m.Id,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
		CreatedBy:    "",
		UpdatedAt:    m.UpdateTime,
		UpdatedAtStr: utils.TimeFormat(m.UpdateTime),
		UpdatedBy:    "",
		Title:        m.Title,
		Path:         m.Path,
	}
}

var defaultApiPaths = []*ApiPath{
	{
		Id:         "1",
		Title:      "管理员:添加",
		Path:       "/api/ms/add/admin",
		Remark:     "",
		CreateTime: 0,
		UpdateTime: 0,
	},
}

func initApiPath(tx *gorm.DB) {
	for _, model := range defaultApiPaths {
		insertIfNotFound(tx, model.Id, model)
	}
}
