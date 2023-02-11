package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils/xerr"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"strconv"
)

type customBody struct {
	Data  []byte
	ReqId string
}

func (c *customBody) GetData() []byte {
	return c.Data
}

func (c *customBody) GetReqId() string {
	return c.ReqId
}

func (l *ConnLogic) OnReceive(ctx context.Context, c *types.UserConn, typ int, msg []byte) {
	switch websocket.MessageType(typ) {
	case websocket.MessageBinary:
		// 接收到消息
		body := &pb.RequestBody{}
		var bodyData []byte
		err := proto.Unmarshal(msg, body)
		var respBody *pb.ResponseBody
		if err == nil {
			respBody, err = conngateway.OnReceive(body.Method, ctx, c, &customBody{
				Data:  body.Data,
				ReqId: body.ReqId,
			})
			bodyData, _ = proto.Marshal(respBody)
		}
		if err != nil {
			code := pb.ResponseBody_InternalError
			if errors.Is(err, xerr.InvalidParamError) {
				code = pb.ResponseBody_RequestError
				logx.WithContext(ctx).Infof("OnReceiveBody error: %s", err.Error())
			} else {
				logx.WithContext(ctx).Errorf("OnReceiveBody error: %s", err.Error())
			}
			if respBody != nil {
				code = respBody.Code
			}
			bodyData, _ = proto.Marshal(&pb.ResponseBody{
				ReqId:  body.ReqId,
				Method: body.Method,
				Code:   code,
				Data:   nil,
			})
		}
		data, _ := proto.Marshal(&pb.PushBody{
			Event: pb.PushEvent_PushResponseBody,
			Data:  bodyData,
		})
		xtrace.StartFuncSpan(ctx, "ReturnResponse", func(ctx context.Context) {
			err = l.SendMsgToConn(c, data)
		}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
			"dataLength": strconv.Itoa(len(data)),
		}))
		if err != nil {
			logx.WithContext(ctx).Infof("SendMsgToConn error: %s", err.Error())
		}
	default:
		// 无效的消息类型
		l.Infof("invalid message type: %d, msg: %s", typ, string(msg))
	}
}
