package wrapper

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xhttp"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

func WrapHandler[REQ IReq, RESP IResp](
	svcCtx *svc.ServiceContext,
	config Config[REQ, RESP],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := config.NewRequest()
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		err := proto.Unmarshal(body, req)
		if err != nil {
			requestValidateErr(w, err.Error())
			return
		}
		requester := req.GetRequester()
		if requester == nil {
			requestValidateErr(w, "requester is nil")
			return
		}
		requester.Ip = xhttp.GetRequestIP(r)
		requester.Ua = r.UserAgent()
		requester.Token = r.Header.Get("Token")
		requester.OsVersion = r.Header.Get("OsVersion")
		requester.AppVersion = r.Header.Get("AppVersion")
		requester.DeviceId = r.Header.Get("DeviceId")
		requester.Platform = r.Header.Get("Platform")
		requester.DeviceModel = r.Header.Get("DeviceModel")
		requester.Language = r.Header.Get("Language")
		requester.Id = r.Header.Get("UserId")
		resp := logic.NewAuthLogic(r, svcCtx).Auth(requester)
		if resp.Code != pb.CommonResp_Success {
			authError(w, resp)
			return
		}
		response, err := config.Do(r.Context(), req)
		if err != nil {
			internalErr(w, err)
			return
		}
		data, err := proto.Marshal(response)
		if err != nil {
			internalErr(w, err)
			return
		}
		success(w, data, response.GetCommonResp())
	}
}

func success(w http.ResponseWriter, data []byte, commonResp *pb.CommonResp) {
	w.WriteHeader(http.StatusOK) // httpCode: 200
	if commonResp == nil {
		commonResp = pb.NewSuccessResp()
	}
	commonResp.Data = data
	resp, _ := proto.Marshal(commonResp)
	_, _ = w.Write(resp)
}

func internalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	resp, _ := proto.Marshal(&pb.CommonResp{
		Code: pb.CommonResp_InternalError, // httpCode: 500
		Msg:  utils.AnyPtr(err.Error()),
		Data: nil,
	})
	_, _ = w.Write(resp)
}

func requestValidateErr(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	resp, _ := proto.Marshal(&pb.CommonResp{
		Code: pb.CommonResp_RequestError, // httpCode: 400
		Msg:  utils.AnyPtr(msg),
		Data: nil,
	})
	_, _ = w.Write(resp)
}

func authError(w http.ResponseWriter, resp *pb.CommonResp) {
	// 401
	w.WriteHeader(http.StatusUnauthorized) // httpCode: 401
	respBytes, _ := proto.Marshal(resp)
	_, _ = w.Write(respBytes)
}
