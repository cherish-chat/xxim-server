package middleware

import (
	"bytes"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

func ApiLog(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	if len(svcCtx.Config.Http.ApiLog.Apis) == 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	} else {
		return func(c *gin.Context) {
			matched := false
			for _, api := range svcCtx.Config.Http.ApiLog.Apis {
				// 格式: GET r'^/api/v1/user/.*' 表示所有以 /api/v1/user/ 开头的 GET 请求都会被记录
				// 取出 Method 和 PathRegex
				apiSplit := strings.Split(api, " ")
				if len(apiSplit) != 2 {
					logx.Errorf("config.Http.ApiLog.Apis is invalid: %s", api)
					os.Exit(1)
				}
				method := strings.ToUpper(apiSplit[0])
				// 取出 Path
				pathRegex := apiSplit[1]
				// 判断是否匹配
				if method == c.Request.Method && pathRegexMatch(pathRegex, c.Request.URL.Path) {
					// 匹配，记录日志
					matched = true
					break
				}
			}
			if matched {
				// 记录日志
				var (
					bodyBuf  []byte
					reqTime  = time.Now()
					clientIp = utils.Http.GetClientIP(c.Request)
				)
				blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
				c.Writer = blw
				body := c.Request.Body
				defer body.Close()
				bodyBuf, _ = io.ReadAll(body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
				defer func() {
					c := c.Copy()
					respTime := time.Now()
					go func() {
						respBody := blw.body.Bytes()
						resp := string(respBody)
						logx.Debugf("clientIp: %s, method: %s, path: %s, requestTime: %s, responseTime: %s, request: %s, response: %s",
							clientIp, c.Request.Method, c.Request.URL.Path, reqTime.Format("2006-01-02 15:04:05.000"), respTime.Format("2006-01-02 15:04:05.000"), string(bodyBuf), resp,
						)
						// TODO: 推送到消息队列?
					}()
				}()
				c.Next()
			}
		}
	}
}

func pathRegexMatch(pathRegex string, path string) (matched bool) {
	// 如果 pathRegex 以 r 开头, 则认为是正则表达式，如果以 / 开头，则认为是完全匹配
	if strings.HasPrefix(pathRegex, "r") {
		// 正则表达式匹配
		// 去掉 r''
		pathRegex = strings.TrimPrefix(pathRegex, "r")
		pathRegex = strings.TrimPrefix(pathRegex, "'")
		pathRegex = strings.TrimPrefix(pathRegex, `"`)
		pathRegex = strings.TrimSuffix(pathRegex, "'")
		pathRegex = strings.TrimSuffix(pathRegex, `"`)
		compile := regexp.MustCompile(pathRegex)
		matched = compile.MatchString(path)
		return
	} else {
		// 完全匹配
		matched = pathRegex == path
		return
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
