package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"time"
)

type ApiPath struct {
	Id         string `gorm:"column:id;primarykey;comment:'主键'"`
	Title      string `gorm:"column:title;not null;default:'';comment:'api标题''"`
	Path       string `gorm:"column:path;not null;default:'';comment:'api路径'"`
	Remark     string `gorm:"column:remark;not null;default:'';comment:'备注信息'"`
	LogEnable  bool   `gorm:"column:logEnable;not null;default:0;comment:'是否记录日志'"`
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
		LogEnable:    m.LogEnable,
	}
}

var defaultApiPaths = []*ApiPath{{
	Id:         "login",
	Title:      "登录",
	Path:       "/api/ms/login",
	LogEnable:  true,
	CreateTime: time.Now().UnixMilli(),
	UpdateTime: time.Now().UnixMilli(),
}}

func init() {
	initDefaultApiPaths("ms", map[string]string{
		"apipath":     "服务器api",
		"admin":       "管理员",
		"ipwhitelist": "ip白名单",
		"menu":        "管理系统菜单",
		"role":        "角色",
	})
	initDefaultApiPaths("server", map[string]string{
		"config": "服务器配置",
	})
	initDefaultApiPaths("appmgmt", map[string]string{
		"config": "app配置",
	})
}

func initDefaultApiPaths(group string, services map[string]string) {
	now := time.Now().UnixMilli()
	for k, v := range services {
		defaultApiPaths = append(defaultApiPaths, &ApiPath{
			Id:         utils.GenId(),
			Title:      v + ":添加",
			Path:       "/api/" + group + "/add/" + k,
			Remark:     "",
			LogEnable:  true,
			CreateTime: now,
			UpdateTime: now,
		})
		defaultApiPaths = append(defaultApiPaths, &ApiPath{
			Id:         utils.GenId(),
			Title:      v + ":删除",
			Path:       "/api/" + group + "/delete/" + k,
			Remark:     "",
			LogEnable:  true,
			CreateTime: now,
			UpdateTime: now,
		})
		defaultApiPaths = append(defaultApiPaths, &ApiPath{
			Id:         utils.GenId(),
			Title:      v + ":编辑",
			Path:       "/api/" + group + "/update/" + k,
			Remark:     "",
			LogEnable:  true,
			CreateTime: now,
			UpdateTime: now,
		})
	}
}

func initApiPath(tx *gorm.DB) {
	for _, model := range defaultApiPaths {
		upsert(tx, model.Id, model)
	}
}
