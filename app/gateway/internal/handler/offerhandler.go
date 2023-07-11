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
	"github.com/pion/sdp/v2"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
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
	sd := &sdp.SessionDescription{}
	if err := sd.Unmarshal([]byte(in.SDP)); err != nil {
		logx.Errorf("sd.Unmarshal error: %v", err)
		return nil, err
	}
	header := &pb.RequestHeader{ClientIp: utils.Sdp.GetClientIp(sd)}
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

	err = peerConnection.SetRemoteDescription(*in)
	if err != nil {
		logx.Errorf("peerConnection.SetRemoteDescription error: %v", err)
		return nil, err
	}

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		id := d.ID()
		logx.Infof("New DataChannel %s %d", d.Label(), id)

		connection := gatewayservicelogic.NewP2pConnection(ctx, header, d)
		// Register channel opening handling
		d.OnOpen(func() {
			logx.Infof("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds", d.Label(), id)
			gatewayservicelogic.ConnectionLogic.OnConnect(connection)
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
			code, response, err := h.onRequest(ctx, connection, msg.Data)
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
			gatewayservicelogic.ConnectionLogic.OnDisconnect(connection)
			_ = peerConnection.Close()
			logx.Infof("DataChannel '%s'-'%d' closed", d.Label(), id)
		})

	})

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		logx.Errorf("peerConnection.CreateAnswer error: %v", err)
		return nil, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		logx.Errorf("peerConnection.SetLocalDescription error: %v", err)
		return nil, err
	}

	// 返回answer
	return &answer, nil
}

func (h *OfferHandler) returnResponse(ctx context.Context, connection *gatewayservicelogic.Connection, code pb.ResponseCode, response []byte, err error) {
	connection.SendMessage(ctx, response)
}

func (h *OfferHandler) onRequest(ctx context.Context, connection *gatewayservicelogic.Connection, msg []byte) (pb.ResponseCode, []byte, error) {
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
			data, _ := proto.Marshal(&pb.GatewayWriteDataContent{
				DataType: pb.GatewayWriteDataType_Response,
				Response: &pb.GatewayApiResponse{
					Header:    i18n.NewInvalidDataError(err.Error()),
					RequestId: "",
					Path:      "",
					Body:      nil,
				},
				Message: nil,
				Notice:  nil,
			})
			return pb.ResponseCode_INVALID_DATA, data, fmt.Errorf("handle message error: %v", err)
		}
	}

	apiRequest := &pb.GatewayApiRequest{}
	err := proto.Unmarshal(msg, apiRequest)
	if err != nil {
		logx.Errorf("proto.Unmarshal error: %v", err)
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
	apiRequest.Header = connection.GetHeader()
	route, ok := universalRouteMap[apiRequest.Path]
	tracer := otel.Tracer(common.TraceName)
	propagator := otel.GetTextMapPropagator()
	spanName := apiRequest.Path
	carrier := propagation.MapCarrier{
		"appId":        apiRequest.Header.AppId,
		"userId":       apiRequest.Header.UserId,
		"clientIp":     apiRequest.Header.ClientIp,
		"installId":    apiRequest.Header.InstallId,
		"platform":     apiRequest.Header.Platform.String(),
		"gatewayPodIp": utils.GetPodIp(),
		"deviceModel":  apiRequest.Header.DeviceModel,
		"osVersion":    apiRequest.Header.OsVersion,
		"appVersion":   apiRequest.Header.AppVersion,
		"connectTime":  connection.ConnectedTime.Format("2006-01-02 15:04:05"),
		"extra":        apiRequest.Header.Extra,
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
	return route(ctx, connection, apiRequest)
}
