package handler

import (
	"context"
	"errors"
	"fmt"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"io"
	"math"
	"nhooyr.io/websocket"
	"strings"
)

type WsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewWsHandler(svcCtx *svc.ServiceContext) *WsHandler {
	return &WsHandler{
		svcCtx: svcCtx,
	}
}

func (h *WsHandler) Upgrade(ginCtx *gin.Context) {
	r := ginCtx.Request
	w := ginCtx.Writer
	logger := logx.WithContext(r.Context())
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	header := &pb.RequestHeader{
		ClientIp: utils.Http.GetClientIP(r),
	}
	compressionMode := websocket.CompressionNoContextTakeover
	// https://github.com/nhooyr/websocket/issues/218
	// 如果是Safari浏览器，不压缩
	if strings.Contains(r.UserAgent(), "Safari") {
		compressionMode = websocket.CompressionDisabled
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:         nil,
		InsecureSkipVerify:   true,
		OriginPatterns:       nil,
		CompressionMode:      compressionMode,
		CompressionThreshold: 0,
	})
	if err != nil {
		// 如果是 / 说明是健康检查
		if r.URL.Path == "/" {
			return
		}
		logger.Errorf("failed to accept websocket connection: %v", err)
		return
	}
	c.SetReadLimit(math.MaxInt32)
	beforeConnectResp, err := h.svcCtx.CallbackService.UserBeforeConnect(r.Context(), &pb.UserBeforeConnectReq{Header: header})
	if err != nil {
		logger.Errorf("beforeConnect error: %v", err)
		c.Close(websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeServerInternalError), err.Error())
		return
	}
	if !beforeConnectResp.Success {
		c.Close(websocket.StatusCode(beforeConnectResp.CloseCode), beforeConnectResp.CloseReason)
		return
	}

	header.UserId = beforeConnectResp.UserId

	defer c.Close(websocket.StatusInternalError, "")

	ctx, cancelFunc := context.WithCancel(r.Context())
	connection := gatewayservicelogic.NewWebsocketConnect(ctx, header, c)
	defer func() {
		gatewayservicelogic.ConnectionLogic.OnDisconnect(connection)
	}()
	go func() {
		// 读取消息
		defer cancelFunc()
		for {
			logger.Debugf("start read")
			typ, msg, err := c.Read(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// 正常关闭
				} else if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway {
					// 正常关闭
					logx.Infof("websocket closed: %v", err)
				} else if strings.Contains(err.Error(), "connection reset by peer") {
					// 网络断开
					logx.Infof("websocket closed: %v", err)
				} else if strings.Contains(err.Error(), "corrupt input") {
					// 输入数据错误
					logx.Infof("websocket closed: %v", err)
				} else {
					logx.Errorf("failed to read message: %v", err)
				}
				return
			}
			go func() {
				_, _ = h.onReceive(ctx, connection, typ, msg)
			}()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (h *WsHandler) onReceive(ctx context.Context, connection *gatewayservicelogic.Connection, typ websocket.MessageType, msg []byte) (pb.ResponseCode, error) {
	var aesKey []byte
	var aesIv []byte
	var isEncrypt bool

	connection.PublicKeyLock.RLock()
	{
		if len(connection.SharedSecret) == 0 {
			// 不加密
			isEncrypt = false
		} else {
			// 加密
			isEncrypt = true
			aesKey = connection.SharedSecret[:]
			aesIv = connection.SharedSecret[8:24]
		}
	}
	connection.PublicKeyLock.RUnlock()

	if isEncrypt {
		var err error
		msg, err = utils.Aes.Decrypt(aesKey, aesIv, msg)
		if err != nil {
			logx.Errorf("decrypt message error: %v", err)
			return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", err)
		}
	}

	apiRequest := &pb.GatewayApiRequest{}
	err := proto.Unmarshal(msg, apiRequest)
	if err != nil {
		return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", err)
	}
	apiRequest.Header = connection.GetHeader()
	route, ok := universalRouteMap[apiRequest.Path]
	tracer := otel.Tracer(common.TraceName)
	propagator := otel.GetTextMapPropagator()
	spanName := apiRequest.Path
	carrier := propagation.MapCarrier{
		"appId":       apiRequest.Header.AppId,
		"userId":      apiRequest.Header.UserId,
		"clientIp":    apiRequest.Header.ClientIp,
		"installId":   apiRequest.Header.InstallId,
		"platform":    apiRequest.Header.Platform.String(),
		"deviceModel": apiRequest.Header.DeviceModel,
		"osVersion":   apiRequest.Header.OsVersion,
		"appVersion":  apiRequest.Header.AppVersion,
		"connectTime": connection.ConnectedTime.Format("2006-01-02 15:04:05"),
		"extra":       apiRequest.Header.Extra,
	}
	spanCtx := propagator.Extract(ctx, carrier)
	spanCtx, span := tracer.Start(spanCtx, spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	)
	defer span.End()
	propagator.Inject(spanCtx, carrier)
	if !ok {
		// 404
		logx.Errorf("path 404 not found: %s", apiRequest.Path)
		span.SetStatus(codes.Error, "path"+apiRequest.Path+"404 not found")
		return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", "path 404 not found")
	}
	code, responseBody, err := route(spanCtx, connection, apiRequest)
	if len(responseBody) > 0 {
		// 发送消息
		err := connection.SendMessage(ctx, responseBody)
		if err != nil {
			logx.Infof("failed to write message: %v", err)
		}
	}
	span.SetAttributes(attribute.Int("responseBody.length", len(responseBody)))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return code, err
}
