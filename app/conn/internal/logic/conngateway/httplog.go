package conngateway

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

func HttpReqLog[REQ IReq, RESP IResp](ctx *gin.Context, method string, req REQ, resp RESP, err error) {
	reqStr := utils.AnyToString(req)
	respStr := utils.AnyToString(resp)
	if err != nil {
		logx.WithContext(ctx.Request.Context()).Errorf("method: %s, req: %s, resp: %s, error: %v", method, reqStr, respStr, err)
	} else {
		logx.WithContext(ctx.Request.Context()).Debugf("method: %s, req: %s, resp: %s", method, reqStr, respStr)
	}
}
