package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/utils/xaes"
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
		{
			// 解密
			if c.ConnParam.AesKey != nil && c.ConnParam.AesIv != nil {
				// aes解密
				var err error
				msg, err = xaes.Decrypt([]byte(*c.ConnParam.AesIv), []byte(*c.ConnParam.AesKey), msg)
				if err != nil {
					l.Errorf("【疑似攻击】userId: %s, ip: %s, ip2region: %s", c.ConnParam.UserId, c.ConnParam.Ips, ip2region.Ip2Region(c.ConnParam.Ips).String())
					c.Conn.Close(int(websocket.StatusPolicyViolation), "protocol error")
					return
				}
			}
		}
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
			var data []byte
			if respBody != nil {
				code = respBody.Code
				data = respBody.Data
			}
			bodyData, _ = proto.Marshal(&pb.ResponseBody{
				ReqId:  body.ReqId,
				Method: body.Method,
				Code:   code,
				Data:   data,
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
