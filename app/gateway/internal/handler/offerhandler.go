package handler

import (
	"context"
	"fmt"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"time"
)

func NewOfferHandler(svcCtx *svc.ServiceContext) *OfferHandler {
	return &OfferHandler{svcCtx: svcCtx}
}

type OfferHandler struct {
	svcCtx *svc.ServiceContext
}

func (h *OfferHandler) Offer(in *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: h.svcCtx.Config.Cloudx.StunUrls,
			},
		},
	}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		logx.Errorf("webrtc.NewPeerConnection error: %v", err)
		return nil, err
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	ctx, cancelFunction := context.WithCancel(context.Background())
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			cancelFunction()
			logx.Infof("Peer Connection State has changed: %s", connectionState.String())
		} else {
			logx.Infof("Peer Connection State has changed: %s", connectionState.String())
		}
	})
	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		id := d.ID()
		connectionId := utils.Snowflake.Int64()
		logx.Infof("New DataChannel %s %d", d.Label(), id)

		now := time.Now()
		connectTime := &now
		var connection *gatewayservicelogic.UniversalConnection
		// Register channel opening handling
		d.OnOpen(func() {
			logx.Infof("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds", d.Label(), id)
			now := time.Now()
			*connectTime = now
			connection, _ = gatewayservicelogic.WsManager.AddSubscriber(ctx, &pb.RequestHeader{
				AppId:        h.svcCtx.Config.Cloudx.AppId,
				UserId:       "",
				UserToken:    "",
				ClientIp:     "",
				InstallId:    "",
				Platform:     0,
				GatewayPodIp: utils.GetPodIp(),
				DeviceModel:  "",
				OsVersion:    "",
				AppVersion:   "",
				Language:     0,
				ConnectTime:  (*connectTime).UnixMilli(),
				Encoding:     pb.EncodingProto_PROTOBUF,
				Extra:        "",
			}, gatewayservicelogic.NewRtcForConnection(d), connectionId)
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			logx.Infof("Message from DataChannel '%s': '%s'", d.Label(), string(msg.Data))
			for {
				if connection != nil {
					break
				}
				time.Sleep(time.Millisecond * 10)
			}
			code, response, err := h.onRequest(ctx, connection, connectTime, msg.Data)
			if err != nil {
				logx.Errorf("h.onRequest error: %v", err)
				h.returnResponse(ctx, connection, code, response, err)
				return
			}
			h.returnResponse(ctx, connection, code, response, nil)
			_ = response
			_ = code
		})

		d.OnClose(func() {
			cancelFunction()
			_ = gatewayservicelogic.WsManager.RemoveSubscriber(connection.GetHeader(), connectionId, pb.WebsocketCustomCloseCode(1000), "DataChannel Close")
			logx.Infof("DataChannel '%s'-'%d' closed", d.Label(), id)
		})

	})

	err = peerConnection.SetRemoteDescription(*in)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// 返回answer
	return &answer, nil
}

func (h *OfferHandler) returnResponse(ctx context.Context, connection *gatewayservicelogic.UniversalConnection, code pb.ResponseCode, response []byte, err error) {
	connection.Connection.Write(ctx, response)
}

func (h *OfferHandler) onRequest(ctx context.Context, connection *gatewayservicelogic.UniversalConnection, connectTime *time.Time, msg []byte) (pb.ResponseCode, []byte, error) {
	apiRequest := &pb.GatewayApiRequest{}
	err := proto.Unmarshal(msg, apiRequest)
	if err != nil {
		data, _ := proto.Marshal(&pb.GatewayWriteDataContent{
			DataType: pb.GatewayWriteDataType_Response,
			Response: &pb.GatewayApiResponse{
				Header:    i18n.NewInvalidDataError(err.Error()),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			},
			Message: nil,
			Notice:  nil,
		})
		return pb.ResponseCode_INVALID_DATA, data, fmt.Errorf("handle message error: %v", err)
	}
	if apiRequest.Header == nil || apiRequest.Header.AppId == "" {
		apiRequest.Header = connection.GetHeader()
	} else if apiRequest.Header.Extra != "" {
		m := make(map[string]any)
		err := utils.Json.Unmarshal([]byte(apiRequest.Header.Extra), &m)
		if err == nil {
			if _, resetHeader := m["resetHeader"]; resetHeader {
				apiRequest.Header = connection.ReSetHeader(apiRequest.Header)
			}
		}
	}
	route, ok := universalRouteMap[apiRequest.Path]
	tracer := otel.Tracer(common.TraceName)
	propagator := otel.GetTextMapPropagator()
	spanName := apiRequest.Path
	carrier := propagation.MapCarrier{
		"appId":  apiRequest.Header.AppId,
		"userId": apiRequest.Header.UserId,
		//"clientIp":     apiRequest.Header.ClientIp, // 获取不到
		"installId":    apiRequest.Header.InstallId,
		"platform":     apiRequest.Header.Platform.String(),
		"gatewayPodIp": utils.GetPodIp(),
		"deviceModel":  apiRequest.Header.DeviceModel,
		"osVersion":    apiRequest.Header.OsVersion,
		"appVersion":   apiRequest.Header.AppVersion,
		"language":     apiRequest.Header.Language.String(),
		"connectTime":  utils.Number.Int64ToString((*connectTime).UnixMilli()),
		//"encoding":     apiRequest.Header.Encoding.ContentType(), // 只能 protobuf
		"extra": apiRequest.Header.Extra,
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
		data, _ := proto.Marshal(&pb.GatewayWriteDataContent{
			DataType: pb.GatewayWriteDataType_Response,
			Response: &pb.GatewayApiResponse{
				Header:    i18n.NewInvalidMethodError(),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			},
			Message: nil,
			Notice:  nil,
		})
		return pb.ResponseCode_INVALID_METHOD, data, fmt.Errorf("handle message error: %v", "path 404 not found")
	}

	beforeConnectResp, err := h.svcCtx.CallbackService.UserBeforeRequest(ctx, &pb.UserBeforeRequestReq{
		Header: apiRequest.Header,
		Path:   apiRequest.Path,
	})
	if err != nil {
		logx.Errorf("beforeConnect error: %v", err)
		data, _ := proto.Marshal(&pb.GatewayWriteDataContent{
			DataType: pb.GatewayWriteDataType_Response,
			Response: &pb.GatewayApiResponse{
				Header:    i18n.NewAuthError(pb.AuthErrorTypeInvalid, err.Error()),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			},
			Message: nil,
			Notice:  nil,
		})
		return pb.ResponseCode_SERVER_ERROR, data, fmt.Errorf("handle message error: %v", err)
	}
	if beforeConnectResp.GetHeader().GetCode() != pb.ResponseCode_SUCCESS {
		data, _ := proto.Marshal(&pb.GatewayWriteDataContent{
			DataType: pb.GatewayWriteDataType_Response,
			Response: &pb.GatewayApiResponse{
				Header:    i18n.NewAuthError(pb.AuthErrorTypeInvalid, beforeConnectResp.Header.Code.String()),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			},
			Message: nil,
			Notice:  nil,
		})
		return pb.ResponseCode_UNAUTHORIZED, data, fmt.Errorf("handle message error: %v", "authentication failed")
	} else {
		apiRequest.Header.UserId = beforeConnectResp.GetUserId()
	}
	return route(ctx, connection, apiRequest)
}
