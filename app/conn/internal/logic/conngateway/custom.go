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
	GetEvent() pb.ActiveEvent
}

func OnReceiveCustom[REQ IReq, RESP IResp](
	ctx context.Context,
	c *types.UserConn,
	body IBody,
	req REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
) (*pb.ResponseBody, error) {
	err := proto.Unmarshal(body.GetData(), req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("%s unmarshal error: %s", req.Path(), err.Error())
		return nil, err
	}
	var resp RESP
	xtrace.StartFuncSpan(ctx, req.Path(), func(ctx context.Context) {
		req.SetCommonReq(&pb.CommonReq{
			UserId:      c.ConnParam.UserId,
			Token:       c.ConnParam.Token,
			DeviceModel: c.ConnParam.DeviceModel,
			DeviceId:    c.ConnParam.DeviceId,
			OsVersion:   c.ConnParam.OsVersion,
			Platform:    c.ConnParam.Platform,
			AppVersion:  c.ConnParam.AppVersion,
			Language:    c.ConnParam.Language,
			Ip:          c.ConnParam.Ips,
		})
		resp, err = do(ctx, req)
	}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
		"req-id": body.GetReqId(),
		"event":  body.GetEvent().String(),
	}))
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("%s error: %s", req.Path(), err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	// 请求日志
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(ctx), "log", func(ctx context.Context) {
		reqLog(c, body, req, resp, err)
	}, nil)
	return &pb.ResponseBody{
		Event: body.GetEvent(),
		ReqId: body.GetReqId(),
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}
