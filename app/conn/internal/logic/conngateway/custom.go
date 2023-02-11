package conngateway

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type IBody interface {
	GetData() []byte
	GetReqId() string
}

func OnReceiveCustom[REQ IReq, RESP IResp](
	ctx context.Context,
	method string,
	c *types.UserConn,
	body IBody,
	req REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
	callback func(ctx context.Context, resp RESP, c *types.UserConn),
) (*pb.ResponseBody, error) {
	logger := logx.WithContext(ctx)
	err := proto.Unmarshal(body.GetData(), req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("%s unmarshal error: %s", method, err.Error())
		return nil, err
	}
	commonReq := &pb.CommonReq{
		UserId:      c.ConnParam.UserId,
		Token:       c.ConnParam.Token,
		DeviceModel: c.ConnParam.DeviceModel,
		DeviceId:    c.ConnParam.DeviceId,
		OsVersion:   c.ConnParam.OsVersion,
		Platform:    c.ConnParam.Platform,
		AppVersion:  c.ConnParam.AppVersion,
		Language:    c.ConnParam.Language,
		Ip:          c.ConnParam.Ips,
	}
	var beforeRequestResp *pb.BeforeRequestResp
	// BeforeRequest
	{
		xtrace.StartFuncSpan(ctx, method+"/BeforeRequest", func(ctx context.Context) {
			beforeRequestResp, err = svcCtx.ImService().BeforeRequest(ctx, &pb.BeforeRequestReq{CommonReq: commonReq, Method: method})
			if err != nil {
				logger.Errorf("BeforeRequest err: %v", err)
				return
			}
		})
		if err != nil {
			logger.Errorf("BeforeRequest err: %v", err)
			return &pb.ResponseBody{
				ReqId:  body.GetReqId(),
				Method: method,
				Code:   pb.ResponseBody_AuthError,
			}, nil
		} else {
			if beforeRequestResp.GetCommonResp().GetCode() != pb.CommonResp_Success {
				data, _ := proto.Marshal(beforeRequestResp.CommonResp)
				logger.Errorf("BeforeRequest err: %v", beforeRequestResp.GetCommonResp().GetMsg())
				return &pb.ResponseBody{
					ReqId:  body.GetReqId(),
					Method: method,
					Code:   pb.ResponseBody_Code(beforeRequestResp.GetCommonResp().GetCode()),
					Data:   data,
				}, nil
			}
		}
	}
	var resp RESP
	xtrace.StartFuncSpan(ctx, method, func(ctx context.Context) {
		req.SetCommonReq(commonReq)
		resp, err = do(ctx, req)
	}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
		"req-id": body.GetReqId(),
	}))
	if err != nil {
		logger.Errorf("%s error: %s", method, err.Error())
	} else {
		if callback != nil {
			callback(ctx, resp, c)
		}
	}
	var respBuff []byte
	if resp.GetCommonResp().GetCode() == pb.CommonResp_Success {
		respBuff, _ = proto.Marshal(resp)
	} else {
		respBuff, _ = proto.Marshal(resp.GetCommonResp())
	}
	// 请求日志
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(ctx), "log", func(ctx context.Context) {
		ReqLog(c, method, body, req, resp, err)
	}, nil)
	return &pb.ResponseBody{
		Method: method,
		ReqId:  body.GetReqId(),
		Code:   pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:   respBuff,
	}, err
}
