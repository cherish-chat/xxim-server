package xhttp

import (
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	AllowOrigins = []string{
		"https://xxim.cherish.chat",
		"http://xxim.cherish.chat",
		"https://enterprise.cherish.chat",
		"http://enterprise.cherish.chat",
	}
)

// Cors gin跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 判断 Origin 请求头是否在允许的列表中
		origin := c.Request.Header.Get("Origin")
		// 如果是 localhost 则允许跨域
		if strings.HasPrefix(origin, "http://localhost") {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			// 如果不是 localhost 则判断是否在允许的列表中
			for _, allowOrigin := range AllowOrigins {
				if allowOrigin == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
