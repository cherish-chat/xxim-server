package middleware

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Cors(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	if !svcCtx.Config.Http.Cors.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	} else {
		config := svcCtx.Config.Http.Cors
		return func(c *gin.Context) {
			method := c.Request.Method
			c.Header("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ","))
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
			c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
			if method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
		}
	}
}
