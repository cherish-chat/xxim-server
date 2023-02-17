package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"strconv"
	"strings"
	"time"
)

var enableLogApiPathsInited bool
var enableLogPathTitleMap = make(map[string]string)

func initEnableLogApiPaths(tx *gorm.DB) []*mgmtmodel.ApiPath {
	var models []*mgmtmodel.ApiPath
	tx.Model(&mgmtmodel.ApiPath{}).Where("logEnable = ?", true).Find(&models)
	for _, model := range models {
		enableLogPathTitleMap[model.Path] = model.Title
	}
	return models
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Log(tx *gorm.DB) gin.HandlerFunc {
	if !enableLogApiPathsInited {
		initEnableLogApiPaths(tx)
	}
	return func(ctx *gin.Context) {
		// 只有POST请求才需要记录
		if ctx.Request.Method != "POST" {
			ctx.Next()
			return
		}
		// 只有 json 请求才需要记录
		if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
			ctx.Next()
			return
		}
		ip := handler.GetClientIP(ctx)
		path := ctx.Request.URL.Path
		title, ok := enableLogPathTitleMap[path]
		if !ok {
			ctx.Next()
			return
		}
		var (
			bodyBuf []byte
			reqTime = time.Now()
		)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		body := ctx.Request.Body
		defer body.Close()
		bodyBuf, _ = io.ReadAll(body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
		defer func() {
			c := ctx.Copy()
			respTime := time.Now()
			go func() {
				var (
					success   bool
					resultMsg string
					resp      string
				)

				respBody := blw.body.Bytes()
				responseBody := &handler.ResponseBody{}
				err := json.Unmarshal(respBody, responseBody)
				if err == nil {
					success = responseBody.Code == pb.CommonResp_Success
					resultMsg = responseBody.Msg
					resp = utils.AnyToString(responseBody.Data)
				} else {
					success = false
					if blw.ResponseWriter.Status() != 200 {
						resultMsg = "请求失败code:" + strconv.Itoa(blw.ResponseWriter.Status())
					} else {
						resultMsg = "解析返回结果失败"
					}
					resp = string(respBody)
				}
				model := &mgmtmodel.OperationLog{
					Id:             utils.GenId(),
					ReqTime:        reqTime.UnixMilli(),
					ReqTimeStr:     reqTime.Format("2006-01-02 15:04:05.000"),
					RespTime:       respTime.UnixMilli(),
					RespTimeStr:    respTime.Format("2006-01-02 15:04:05.000"),
					OperationType:  getOperationType(path),
					OperationTitle: title,
					ReqPath:        path,
					ReqParams:      utils.AnyToString(bodyBuf),
					ResultSuccess:  c.Writer.Status() == 200 && success,
					ReqResultMsg:   resultMsg,
					Resp:           resp,
					ReqIp:          ip,
					IpSource:       ip2region.Ip2Region(ip).String(),
					ReqCost:        respTime.UnixMilli() - reqTime.UnixMilli(),
					Operator:       c.GetHeader("UserId"),
				}
				tx.Create(model)
			}()
		}()
		ctx.Next()
	}
}

func getOperationType(path string) string {
	if strings.Contains(path, "/add") {
		return "add"
	}
	if strings.Contains(path, "/update") {
		return "update"
	}
	if strings.Contains(path, "/delete") {
		return "delete"
	}
	if strings.Contains(path, "/list") {
		return "list"
	}
	return "other"
}
