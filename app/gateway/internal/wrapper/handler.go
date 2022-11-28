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
		requester := &pb.CommonReq{}
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		err := proto.Unmarshal(body, requester)
		if err != nil {
			RequestValidateErr(w, err.Error())
			return
		}
		req := config.NewRequest()
		err = proto.Unmarshal(requester.Data, req)
		if err != nil {
			RequestValidateErr(w, err.Error())
			return
		}
		requester.Ip = xhttp.GetRequestIP(r)
		requester.UserAgent = r.UserAgent()
		resp := logic.NewAuthLogic(r, svcCtx).Auth(requester)
		if resp.Code != pb.CommonResp_Success {
			AuthError(w, resp)
			return
		}
		requester.Data = nil
		req.SetCommonReq(requester)
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
		Success(w, data, response.GetCommonResp())
	}
}

func Success(w http.ResponseWriter, data []byte, commonResp *pb.CommonResp) {
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

func RequestValidateErr(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	resp, _ := proto.Marshal(&pb.CommonResp{
		Code: pb.CommonResp_RequestError, // httpCode: 400
		Msg:  utils.AnyPtr(msg),
		Data: nil,
	})
	_, _ = w.Write(resp)
}

func AuthError(w http.ResponseWriter, resp *pb.CommonResp) {
	// 401
	w.WriteHeader(http.StatusUnauthorized) // httpCode: 401
	respBytes, _ := proto.Marshal(resp)
	_, _ = w.Write(respBytes)
}
