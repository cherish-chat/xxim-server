package handler

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

type commonRespGetter interface {
	GetCommonResp() *pb.CommonResp
}

type ResponseBody struct {
	Data any                `json:"data"`
	Code pb.CommonResp_Code `json:"code"`
	Msg  string             `json:"msg"`
	Show bool               `json:"show"`
}

func ReturnOk(ctx *gin.Context, data any) {
	body := &ResponseBody{
		Data: data,
		Code: 0,
		Msg:  "",
		Show: false,
	}
	if iCommonResp, ok := data.(commonRespGetter); ok {
		commonResp := iCommonResp.GetCommonResp()
		if commonResp == nil {
			commonResp = pb.NewSuccessResp()
		}
		body.Code = commonResp.Code
		body.Msg = utils.If(commonResp.Msg != nil, *commonResp.Msg, "")
	}
	ctx.JSON(200, body)
}

func ReturnOkShow(ctx *gin.Context, data any) {
	body := &ResponseBody{
		Data: data,
		Code: 0,
		Msg:  "",
		Show: true,
	}
	if iCommonResp, ok := data.(commonRespGetter); ok {
		commonResp := iCommonResp.GetCommonResp()
		body.Code = commonResp.Code
		body.Msg = utils.If(commonResp.Msg != nil, *commonResp.Msg, "")
	}
	ctx.JSON(200, body)
}

func GetClientIP(ctx *gin.Context) string {
	// 先取 X-Real-IP
	ip := ctx.Request.Header.Get("X-Real-IP")
	if ip == "" {
		// 取 X-Forwarded-For
		ip = ctx.Request.Header.Get("X-Forwarded-For")
		if ip == "" {
			// 取 RemoteAddr
			ip = ctx.Request.RemoteAddr
		}
	}
	return strings.Split(ip, ",")[0]
}
