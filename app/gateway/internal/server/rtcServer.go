package server

import (
	"context"
	"encoding/json"
	"fmt"
	cloudxpb "github.com/cherish-chat/imcloudx-server/common/pb"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/pion/webrtc/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"nhooyr.io/websocket"
	"time"
)

type RtcServer struct {
	svcCtx *svc.ServiceContext
}

func NewRtcServer(svcCtx *svc.ServiceContext) *RtcServer {
	r := &RtcServer{svcCtx: svcCtx}
	return r
}

func (s *RtcServer) Start() {
	host := s.svcCtx.Config.Cloudx.Host
	if host == "" {
		return
	}
	schema := "ws"
	if s.svcCtx.Config.Cloudx.Ssl {
		schema = "wss"
	}
	background := context.Background()
	reconnectCount := 0
	for {
		// 连接到ws服务器
		conn, response, err := websocket.Dial(
			background,
			fmt.Sprintf("%s://%s:%d/ws?clientId=%s&clientSecret=%s",
				schema,
				host,
				s.svcCtx.Config.Cloudx.Port,
				s.svcCtx.Config.Cloudx.ClientId,
				s.svcCtx.Config.Cloudx.ClientSecret,
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
			ticker := time.NewTicker(time.Second * time.Duration(s.svcCtx.Config.Cloudx.KeepLiveSeconds))
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

func (s *RtcServer) onReceive(conn *websocket.Conn, typ websocket.MessageType, data []byte) {
	switch typ {
	case websocket.MessageText:
	case websocket.MessageBinary:
		message := &cloudxpb.WsReceiveMessageBinary{}
		err := proto.Unmarshal(data, message)
		if err != nil {
			logx.Errorf("proto.Unmarshal error: %v", err)
			return
		}
		switch message.Type {
		case cloudxpb.WsReceiveMessageBinaryType_ReceiveRequest:
			request := &cloudxpb.NodeReq{}
			err := proto.Unmarshal(message.Data, request)
			if err != nil {
				logx.Errorf("proto.Unmarshal error: %v", err)
				return
			}
			switch request.Method {
			case "/offer":
				offer := &webrtc.SessionDescription{}
				err := json.Unmarshal(request.Body, offer)
				if err != nil {
					logx.Errorf("json.Unmarshal error: %v", err)
					return
				}

				offerHandler := handler.NewOfferHandler(s.svcCtx)
				answer, err := offerHandler.Offer(offer)
				if err != nil {
					logx.Errorf("onCallOffer error: %v", err)
					return
				}
				answerData, _ := json.Marshal(answer)
				responseData, _ := proto.Marshal(&cloudxpb.NodeResp{
					AppId:     request.AppId,
					RequestId: request.RequestId,
					Headers:   request.Headers,
					Status:    200,
					Body:      answerData,
					ErrMsg:    "",
				})
				responseData, _ = proto.Marshal(&cloudxpb.WsReceiveMessageBinary{
					Type: cloudxpb.WsReceiveMessageBinaryType_ReceiveResponse,
					Data: responseData,
				})
				// 返回
				conn.Write(context.Background(), websocket.MessageBinary, responseData)
			}
		case cloudxpb.WsReceiveMessageBinaryType_ReceiveResponse:
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

// RandSeq generates a random string to serve as dummy data
func RandSeq(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
