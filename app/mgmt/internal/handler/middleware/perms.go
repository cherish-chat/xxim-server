package middleware

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

var allApiPaths []*mgmtmodel.ApiPath
var allApiPathsInited bool
var pathIdMap = make(map[string]string)

func getAllApiPaths(tx *gorm.DB) []*mgmtmodel.ApiPath {
	var models []*mgmtmodel.ApiPath
	tx.Model(&mgmtmodel.ApiPath{}).Find(&models)
	for _, model := range models {
		pathIdMap[model.Path] = model.Id
	}
	return models
}

func Perms(tx *gorm.DB) gin.HandlerFunc {
	if !allApiPathsInited {
		allApiPaths = getAllApiPaths(tx)
		allApiPathsInited = true
	}
	return func(c *gin.Context) {
		if _, ok := dontCheckTokenMap[c.Request.URL.Path]; ok {
			c.Next()
			return
		}
		// 当前path对应的id
		pathId, ok := pathIdMap[c.Request.URL.Path]
		if !ok {
			// 不需要权限
			c.Next()
			return
		}
		userId := c.GetHeader("userId")
		// 查询用户
		user := &mgmtmodel.User{}
		if err := tx.Where("id = ?", userId).First(user).Error; err != nil {
			c.AbortWithStatus(401)
			return
		}
		// 查询role
		role := &mgmtmodel.Role{}
		if err := tx.Where("id = ?", user.RoleId).First(role).Error; err != nil {
			c.AbortWithStatus(401)
			return
		}
		if utils.InSlice(strings.Split(role.ApiPathIds, ","), pathId) {
			// 有权
			c.Next()
			return
		}
		c.AbortWithStatus(401)
	}
}
