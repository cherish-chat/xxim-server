package conngateway

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func reqLog[REQ IReq, RESP IResp](c *types.UserConn, body IBody, req REQ, resp RESP, err error) {
	reqStr := utils.AnyToString(req)
	respStr := utils.AnyToString(resp)
	reqId := body.GetReqId()
	event := body.GetEvent().String()
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("reqId: %s, event: %s, req: %s, resp: %s, error: %v", reqId, event, reqStr, respStr, err)
	} else {
		logx.WithContext(c.Ctx).Infof("reqId: %s, event: %s, req: %s, resp: %s", reqId, event, reqStr, respStr)
	}
}
