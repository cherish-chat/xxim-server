package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"io"
)

var dontCheckTokenMap = map[string]bool{
	"/api/ms/login":  true,
	"/api/ms/health": true,
	"/api/ms/config": true,
}

func Auth(rc *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := dontCheckTokenMap[c.Request.URL.Path]; ok {
			dontCheckToken(c)
			return
		}
		// 请求头中的token、ip
		token := c.GetHeader("token")
		userId := c.GetHeader("userId")
		tokenCode, _ := xjwt.VerifyTokenAdmin(c.Request.Context(), rc, userId, token)
		if xjwt.VerifyTokenCodeOK != tokenCode {
			c.AbortWithStatus(401)
			return
		} else {
			userAgent := c.GetHeader("user-agent")
			ip := handler.GetClientIP(c)
			commonReq := &pb.CommonReq{
				UserId:      userId,
				Token:       token,
				DeviceModel: "web",
				DeviceId:    "web",
				OsVersion:   "web",
				Platform:    "web",
				AppVersion:  "latest",
				Language:    "zh",
				Ip:          ip,
				UserAgent:   userAgent,
			}
			body := `{}`
			if c.Request.Body != nil {
				bodyBytes, _ := c.GetRawData()
				body = string(bodyBytes)
			}
			jsonBody := make(map[string]interface{})
			if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
				// 不做处理
			} else {
				jsonBody["commonReq"] = commonReq
			}
			// 重新放到请求体中
			if newBody, err := json.Marshal(jsonBody); err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
			}
			c.Next()
		}
	}
}

func dontCheckToken(c *gin.Context) {
	userAgent := c.GetHeader("user-agent")
	ip := handler.GetClientIP(c)
	commonReq := &pb.CommonReq{
		DeviceModel: "web",
		DeviceId:    "web",
		OsVersion:   "web",
		Platform:    "web",
		AppVersion:  "latest",
		Language:    "zh",
		Ip:          ip,
		UserAgent:   userAgent,
	}
	body := `{}`
	if c.Request.Body != nil {
		bodyBytes, _ := c.GetRawData()
		body = string(bodyBytes)
	}
	jsonBody := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		// 不做处理
	} else {
		jsonBody["commonReq"] = commonReq
	}
	// 重新放到请求体中
	if newBody, err := json.Marshal(jsonBody); err == nil {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
	}
	c.Next()
}
