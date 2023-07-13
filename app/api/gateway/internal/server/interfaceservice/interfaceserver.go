package server

import (
	"context"
	"fmt"
	"github.com/cherish-chat/cherish-cloud-proto/signalingpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/handler"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"nhooyr.io/websocket"
	"time"
)

type CustomInterfaceServiceServer struct {
	svcCtx *svc.ServiceContext
	engine *gin.Engine
}

func NewCustomInterfaceServiceServer(svcCtx *svc.ServiceContext) *CustomInterfaceServiceServer {
	s := &CustomInterfaceServiceServer{svcCtx: svcCtx, engine: gin.Default()}
	connectionmanager.InitConnectionLogic(svcCtx)
	return s
}

func (s *CustomInterfaceServiceServer) Start() {
	handler.SetupRoutes(s.svcCtx, s.engine)

	if s.svcCtx.Config.Gateway.Mode == "tcp" {
		panic("tcp mode not support, please use p2p mode")
	} else if s.svcCtx.Config.Gateway.Mode == "ws" {
		panic("ws mode not support, please use p2p mode")
	} else {
		if s.svcCtx.Config.Gateway.SignalingServer == "" {
			panic("signaling server not config")
		}
	}
	background := context.Background()
	reconnectCount := 0
	for {
		// 连接到ws服务器
		conn, response, err := websocket.Dial(
			background,
			fmt.Sprintf("%s/ws?appId=%s&appSecret=%s",
				s.svcCtx.Config.Gateway.SignalingServer,
				s.svcCtx.Config.Gateway.AppId,
				s.svcCtx.Config.Gateway.AppSecret,
			),
			nil)
		if err != nil {
			logx.Errorf("websocket.Dial error: %v", err)
			sleep(&reconnectCount)
			continue
		}
		if response.StatusCode != 101 {
			logx.Errorf("websocket.Dial StatusCode error: %d", response.StatusCode)
			sleep(&reconnectCount)
			continue
		}
		// 连接成功，开始读写
		reconnectCount = 0
		shutdownChan := make(chan struct{})
		go func() {
			ticker := time.NewTicker(time.Second * time.Duration(15))
			defer ticker.Stop()
			for {
				select {
				case <-shutdownChan:
					logx.Infof("shutdownChan")
					return
				case <-ticker.C:
					err := conn.Write(background, websocket.MessageText, []byte("ping"))
					if err != nil {
						logx.Errorf("conn.Write error: %v", err)
						continue
					}
				}
			}
		}()
		for {
			typ, data, err := conn.Read(background)
			if err != nil {
				logx.Errorf("conn.Read error: %v", err)
				shutdownChan <- struct{}{}
				sleep(&reconnectCount)
				break
			}
			logx.Infof("conn.Read data: %s", data)
			s.onReceive(conn, typ, data)
		}
	}
}

func (s *CustomInterfaceServiceServer) onReceive(conn *websocket.Conn, typ websocket.MessageType, data []byte) {
	switch typ {
	case websocket.MessageText:
	case websocket.MessageBinary:
		message := &signalingpb.MessageForPeerServer{}
		err := proto.Unmarshal(data, message)
		if err != nil {
			logx.Errorf("proto.Unmarshal error: %v", err)
			return
		}
		switch message.Type {
		case signalingpb.MessageForPeerServer_Request:
			request := message.Request
			switch request.Type {
			case signalingpb.RequestForPeerServer_GetAnswer:
				in := &signalingpb.GetAnswerReq{}
				err := proto.Unmarshal(request.Payload, in)
				if err != nil {
					logx.Errorf("proto.Unmarshal error: %v", err)
					return
				}

				offerHandler := handler.NewOfferHandler(s.svcCtx)
				answer, err := offerHandler.Offer(&webrtc.SessionDescription{
					Type: webrtc.SDPTypeOffer,
					SDP:  in.Sdp,
				})
				if err != nil {
					logx.Errorf("onCallOffer error: %v", err)
					return
				}
				answerData, _ := proto.Marshal(&signalingpb.GetAnswerResp{
					Sdp: answer.SDP,
				})
				responseData, _ := proto.Marshal(&signalingpb.MessageForPeerServer{
					Type: signalingpb.MessageForPeerServer_Response,
					Response: &signalingpb.ResponseForPeerServer{
						Type:      signalingpb.ResponseForPeerServer_GetAnswer,
						RequestId: request.RequestId,
						Payload:   answerData,
					},
				})
				// 返回
				conn.Write(context.Background(), websocket.MessageBinary, responseData)
			}
		case signalingpb.MessageForPeerServer_Response:
			// TODO
		}
	}
}

func sleep(reconnectCount *int) {
	*reconnectCount++
	if *reconnectCount > 10 {
		*reconnectCount = 10
	}
	time.Sleep(time.Duration(*reconnectCount) * time.Second)
}
