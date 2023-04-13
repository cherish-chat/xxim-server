package handler

import (
	"github.com/cherish-chat/xxim-server/app/xos/internal/config"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func Cors(config config.CorsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
