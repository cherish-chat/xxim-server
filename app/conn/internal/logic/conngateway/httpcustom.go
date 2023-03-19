package conngateway

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
)

func OnHttpReceiveCustom[REQ IReq, RESP IResp](
	ctx *gin.Context,
	req REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
) (*pb.CommonResp, error) {
	logger := logx.WithContext(ctx.Request.Context())
	method := ctx.Request.RequestURI
	var err error
	var data []byte
	data, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		errFormat := fmt.Sprintf("read body error: %s", err.Error())
		logger.Error(errFormat)
		return &pb.CommonResp{
			Code: pb.CommonResp_RequestError,
			Data: nil,
			Msg:  utils.AnyPtr(errFormat),
		}, nil
	}
	commonReq := &pb.CommonReq{}
	err = proto.Unmarshal(data, commonReq)
	if err != nil {
		errFormat := fmt.Sprintf("unmarshal error: %s", err.Error())
		logger.Error(errFormat)
		return &pb.CommonResp{
			Code: pb.CommonResp_RequestError,
			Data: nil,
			Msg:  utils.AnyPtr(errFormat),
		}, nil
	}
	if len(commonReq.Data) != 0 {
		err = proto.Unmarshal(commonReq.Data, req)
		if err != nil {
			errFormat := fmt.Sprintf("unmarshal error: %s", err.Error())
			logger.Error(errFormat)
			return &pb.CommonResp{
				Code: pb.CommonResp_RequestError,
				Data: nil,
				Msg:  utils.AnyPtr(errFormat),
			}, nil
		}
	}
	var beforeRequestResp *pb.BeforeRequestResp
	// BeforeRequest
	{
		xtrace.StartFuncSpan(ctx, method+"/BeforeRequest", func(ctx context.Context) {
			beforeRequestResp, err = svcCtx.ImService().BeforeRequest(ctx, &pb.BeforeRequestReq{CommonReq: commonReq, Method: method})
		})
		if err != nil {
			// 判断是不是 status.Error(codes.Unauthenticated, "ip被封禁")
			statusError, ok := status.FromError(err)
			if ok && statusError.Code() == codes.Unauthenticated {
				return pb.NewAuthErrorResp(statusError.Message()), nil
			}
			logger.Errorf("BeforeRequest err: %v", err)
			return pb.NewInternalErrorResp(err.Error()), nil
		} else {
			if beforeRequestResp.GetCommonResp().GetCode() != pb.CommonResp_Success {
				logger.Errorf("BeforeRequest err: %v", beforeRequestResp.GetCommonResp().GetMsg())
				return beforeRequestResp.GetCommonResp(), nil
			}
		}
	}
	var resp RESP
	xtrace.StartFuncSpan(ctx.Request.Context(), method, func(ctx context.Context) {
		req.SetCommonReq(commonReq)
		resp, err = do(ctx, req)
	})
	if err != nil {
		logger.Errorf("%s error: %s", method, err.Error())
	}
	var respBuff []byte
	if resp.GetCommonResp().GetCode() == pb.CommonResp_Success {
		respBuff, _ = proto.Marshal(resp)
	} else {
		respBuff, _ = proto.Marshal(resp.GetCommonResp())
	}
	// 请求日志
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(ctx), "log", func(c context.Context) {
		HttpReqLog(ctx, method, req, resp, err)
	}, nil)
	return &pb.CommonResp{
		Code: resp.GetCommonResp().GetCode(),
		Data: respBuff,
	}, err
}
