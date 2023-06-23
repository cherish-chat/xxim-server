package handler

import (
	"context"
	"encoding/json"
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
		AppId:        r.URL.Query().Get("appId"),
		UserId:       r.URL.Query().Get("userId"),
		UserToken:    r.URL.Query().Get("userToken"),
		ClientIp:     utils.Http.GetClientIP(r),
		InstallId:    r.URL.Query().Get("installId"),
		Platform:     pb.PlatformFromString(r.URL.Query().Get("platform")),
		GatewayPodIp: utils.GetPodIp(),
		DeviceModel:  r.URL.Query().Get("deviceModel"),
		OsVersion:    r.URL.Query().Get("osVersion"),
		AppVersion:   r.URL.Query().Get("appVersion"),
		Language:     pb.LanguageFromString(r.URL.Query().Get("language")),
		ConnectTime:  utils.Time.Now13(),
		Encoding:     pb.EncodingFromString(r.URL.Query().Get("encoding")),
		Extra:        r.URL.Query().Get("extra"),
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
	connectionId := utils.Snowflake.Int64()
	defer func() {
		logger.Debugf("removing subscriber: %d", connectionId)
		err := gatewayservicelogic.WsManager.RemoveSubscriber(header, connectionId, websocket.StatusNormalClosure, "finished")
		if err != nil {
			logger.Errorf("failed to remove subscriber: %v", err)
			return
		} else {
			logger.Debugf("removed subscriber: %d", connectionId)
		}
	}()
	connection, err := gatewayservicelogic.WsManager.AddSubscriber(ctx, header, c, connectionId)
	if err != nil {
		logger.Errorf("failed to add subscriber: %v", err)
		c.Close(websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeServerInternalError), err.Error())
		cancelFunc()
		return
	}
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

func (h *WsHandler) onReceive(ctx context.Context, connection *gatewayservicelogic.WsConnection, typ websocket.MessageType, msg []byte) (pb.ResponseCode, error) {
	apiRequest := &pb.GatewayApiRequest{}
	if connection.Header.Encoding == pb.EncodingProto_JSON {
		err := json.Unmarshal(msg, apiRequest)
		if err != nil {
			return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", err)
		}
	} else if connection.Header.Encoding == pb.EncodingProto_PROTOBUF {
		err := proto.Unmarshal(msg, apiRequest)
		if err != nil {
			return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", err)
		}
	} else {
		return pb.ResponseCode_INVALID_DATA, fmt.Errorf("handle message error: %v", "unsupported encoding")
	}
	apiRequest.Header = connection.Header
	route, ok := wsRouteMap[apiRequest.Path]
	tracer := otel.Tracer(common.TraceName)
	propagator := otel.GetTextMapPropagator()
	spanName := apiRequest.Path
	carrier := propagation.MapCarrier{
		"appId":        apiRequest.Header.AppId,
		"userId":       apiRequest.Header.UserId,
		"clientIp":     apiRequest.Header.ClientIp,
		"installId":    apiRequest.Header.InstallId,
		"platform":     apiRequest.Header.Platform.String(),
		"gatewayPodIp": apiRequest.Header.GatewayPodIp,
		"deviceModel":  apiRequest.Header.DeviceModel,
		"osVersion":    apiRequest.Header.OsVersion,
		"appVersion":   apiRequest.Header.AppVersion,
		"language":     apiRequest.Header.Language.String(),
		"connectTime":  utils.Number.Int64ToString(apiRequest.Header.ConnectTime),
		"encoding":     apiRequest.Header.Encoding.ContentType(),
		"extra":        apiRequest.Header.Extra,
	}
	spanCtx := propagator.Extract(ctx, carrier)
	kvs := []attribute.KeyValue{attribute.Int64("connection.Id", connection.Id)}
	for k, v := range carrier {
		kvs = append(kvs, attribute.String(k, v))
	}
	spanCtx, span := tracer.Start(spanCtx, spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(
			kvs...,
		),
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
		err := connection.Connection.Write(ctx, websocket.MessageBinary, responseBody)
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
