package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherish-chat/xxim-server/common"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/types"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"sync"
	"time"
)

type IClient interface {
	Request(path string, req any, resp any) error
	GatewayGetUserConnection(req *pb.GatewayGetUserConnectionReq) (resp *pb.GatewayGetUserConnectionResp, err error)
	GatewayBatchGetUserConnection(req *pb.GatewayBatchGetUserConnectionReq) (resp *pb.GatewayBatchGetUserConnectionResp, err error)
	GatewayGetConnectionByFilter(req *pb.GatewayGetConnectionByFilterReq) (resp *pb.GatewayGetConnectionByFilterResp, err error)
	GatewayWriteDataToWs(req *pb.GatewayWriteDataToWsReq) (resp *pb.GatewayWriteDataToWsResp, err error)
	GatewayKickWs(req *pb.GatewayKickWsReq) (resp *pb.GatewayKickWsResp, err error)
	UserRegister(req *pb.UserRegisterReq) (resp *pb.UserRegisterResp, err error)
	UserAccessToken(req *pb.UserAccessTokenReq) (resp *pb.UserAccessTokenResp, err error)
}

type HttpClient struct {
	httpClient          *http.Client
	latestEndpointIndex int
	Config              *Config
}

func NewHttpClient(config *Config) (*HttpClient, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return &HttpClient{
		httpClient: http.DefaultClient,
		Config:     config,
	}, nil
}

type WsClient struct {
	wsClient    *websocket.Conn
	httpClient  *HttpClient
	Config      *Config
	responseMap sync.Map // key: requestId, value: chan *types.GatewayApiResponse
}

func NewWsClient(config *Config) (*WsClient, error) {
	httpClient, err := NewHttpClient(config)
	if err != nil {
		return nil, err
	}
	endpoint := config.Endpoints[0]
	url := fmt.Sprintf("%s%s", endpoint, "/ws?")
	params := map[string]string{
		"appId":       config.AppId,
		"userId":      httpClient.getUserId(),
		"userToken":   httpClient.getUserToken(),
		"installId":   config.InstallId,
		"platform":    config.Platform.ToString(),
		"deviceModel": config.DeviceModel,
		"osVersion":   config.OsVersion,
		"appVersion":  common.Version,
		"language":    config.Language.ToString(),
		"encoding":    config.Encoding.ToString(),
	}
	for k, v := range params {
		url += fmt.Sprintf("%s=%s&", k, v)
	}
	url = url[:len(url)-1]
	wsClient, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial error: %v", err)
	}
	w := &WsClient{
		Config:     config,
		wsClient:   wsClient,
		httpClient: httpClient,
	}
	go w.loopRead()
	go w.heartbeat()
	return w, nil
}

func (c *WsClient) loopRead() {
	for {
		_, message, err := c.wsClient.Read(context.Background())
		if err != nil {
			err = utils.Error.DeepUnwrap(err)
			closeError, ok := err.(websocket.CloseError)
			if ok {
				logx.Errorf("read message error: %v, code: %d, reason: %s", err, closeError.Code, closeError.Reason)
				time.Sleep(time.Second * 1)
				os.Exit(1)
				return
			}
			logx.Errorf("read message error: %v", err)
			continue
		}
		var resp pb.GatewayApiResponse
		if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
			err = proto.Unmarshal(message, &resp)
		} else {
			err = json.Unmarshal(message, &resp)
		}
		if err != nil {
			logx.Errorf("unmarshal message error: %v", err)
			continue
		}
		ch, ok := c.responseMap.Load(resp.RequestId)
		if !ok {
			logx.Infof("response not found, data: %s", string(message))
			continue
		}
		ch.(chan *pb.GatewayApiResponse) <- &resp
	}
}

func (c *WsClient) heartbeat() {
	ticker := time.NewTicker(c.Config.KeepAliveSecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go func() {
				logx.Debugf("send heartbeat")
				e := c.Request("/v1/gateway/keepAlive", &pb.GatewayKeepAliveReq{}, &pb.GatewayKeepAliveResp{})
				if e != nil {
					logx.Errorf("heartbeat error: %v", e)
				} else {
					logx.Debugf("heartbeat success")
				}
			}()
		}
	}
}

var (
	ErrInvalidRequestType = errors.New("invalid request type, req must implement types.ReqInterface")
	ErrRequestTimeout     = errors.New("request timeout")
)

func (c *HttpClient) GetUrl(path string) string {
	endpoint := c.Config.Endpoints[c.latestEndpointIndex]
	return fmt.Sprintf("%s/api%s", endpoint, path)
}

func (c *WsClient) GetUrl(path string) string {
	return path
}

func (c *HttpClient) Request(path string, req any, resp any) error {
	var body io.Reader
	if req != nil {
		var data []byte
		var err error
		if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
			var ok bool
			var message types.ReqInterface
			message, ok = req.(types.ReqInterface)
			if !ok {
				return ErrInvalidRequestType
			}
			data, err = proto.Marshal(message)
			if err != nil {
				return fmt.Errorf("req marshal error: %v", err)
			}
			data, _ = proto.Marshal(&pb.GatewayApiRequest{
				Header: &pb.RequestHeader{
					AppId:        c.Config.AppId,
					UserId:       c.getUserId(),
					UserToken:    c.getUserToken(),
					ClientIp:     "", //客户端不需要设置 由服务端设置
					InstallId:    c.Config.InstallId,
					Platform:     *c.Config.Platform,
					GatewayPodIp: "",
					DeviceModel:  c.Config.DeviceModel,
					OsVersion:    c.Config.OsVersion,
					AppVersion:   common.Version,
					Language:     *c.Config.Language,
					ConnectTime:  0,
					Encoding:     pb.EncodingProto_PROTOBUF,
					Extra:        c.Config.CustomHeader,
				},
				Body: data,
			})
		} else {
			data, err = json.Marshal(req)
			if err != nil {
				return fmt.Errorf("req marshal error: %v", err)
			}
			data, _ = json.Marshal(&pb.GatewayApiRequest{
				Header: &pb.RequestHeader{
					AppId:        c.Config.AppId,
					UserId:       c.getUserId(),
					UserToken:    c.getUserToken(),
					ClientIp:     "", //客户端不需要设置 由服务端设置
					InstallId:    c.Config.InstallId,
					Platform:     *c.Config.Platform,
					GatewayPodIp: "",
					DeviceModel:  c.Config.DeviceModel,
					OsVersion:    c.Config.OsVersion,
					AppVersion:   common.Version,
					Language:     *c.Config.Language,
					ConnectTime:  0,
					Encoding:     pb.EncodingProto_JSON,
					Extra:        c.Config.CustomHeader,
				},
				Body: data,
			})
		}
		body = bytes.NewReader(data)
	}
	request, err := http.NewRequest("POST", c.GetUrl(path), body)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}
	// set content type
	if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
		request.Header.Set("Content-Type", "application/x-protobuf")
	} else {
		request.Header.Set("Content-Type", "application/json")
	}
	// do request
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("do request error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code error: %d", response.StatusCode)
	}
	defer response.Body.Close()
	// read response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %v", err)
	}
	// unmarshal response body
	if resp != nil {
		if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
			message, ok := resp.(proto.Message)
			if !ok {
				return fmt.Errorf("invalid response type, resp must implement proto.Message")
			}
			err = proto.Unmarshal(data, message)
			if err != nil {
				return fmt.Errorf("resp unmarshal error: %v", err)
			}
		} else {
			err = json.Unmarshal(data, resp)
			if err != nil {
				return fmt.Errorf("resp unmarshal error: %v", err)
			}
		}
	}
	return nil
}

func (c *WsClient) Request(path string, req any, resp any) error {
	requestId := utils.Snowflake.String()
	var body []byte
	if req != nil {
		if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
			pb, ok := req.(proto.Message)
			if !ok {
				return ErrInvalidRequestType
			}
			data, err := proto.Marshal(pb)
			if err != nil {
				return fmt.Errorf("req marshal error: %v", err)
			}
			body = data
		} else {
			data, err := json.Marshal(req)
			if err != nil {
				return fmt.Errorf("req marshal error: %v", err)
			}
			body = data
		}
	}
	apiRequest := &pb.GatewayApiRequest{
		RequestId: requestId,
		Path:      path,
		Body:      body,
	}
	var data []byte
	var err error
	if *c.Config.Encoding == pb.EncodingProto_JSON {
		data, err = json.Marshal(apiRequest)
		if err != nil {
			return fmt.Errorf("apiRequest marshal error: %v", err)
		}
	} else {
		data, err = proto.Marshal(apiRequest)
		if err != nil {
			return fmt.Errorf("apiRequest marshal error: %v", err)
		}
	}
	ch := c.waitResponse(requestId)
	err = c.wsClient.Write(context.Background(), websocket.MessageBinary, data)
	if err != nil {
		return fmt.Errorf("write error: %v", err)
	}
	select {
	case <-time.After(c.Config.RequestTimeout):
		return ErrRequestTimeout
	case response := <-ch:
		if response.GetHeader().GetCode() != pb.ResponseCode_SUCCESS {
			return fmt.Errorf("response error: %v", response.GetHeader().GetCode())
		}
		if resp != nil {
			getBody := response.GetBody()
			if len(getBody) == 0 {
				return nil
			}
			if *c.Config.Encoding == pb.EncodingProto_PROTOBUF {
				message, ok := resp.(proto.Message)
				if !ok {
					return fmt.Errorf("invalid response type, resp must implement proto.Message")
				}
				err = proto.Unmarshal(getBody, message)
				if err != nil {
					return fmt.Errorf("resp unmarshal error: %v", err)
				}
			} else {
				err = json.Unmarshal(getBody, resp)
				if err != nil {
					return fmt.Errorf("resp unmarshal error: %v", err)
				}
			}
		}
	}
	return nil
}

func (c *WsClient) waitResponse(requestId string) chan *pb.GatewayApiResponse {
	ch := make(chan *pb.GatewayApiResponse)
	c.responseMap.Store(requestId, ch)
	return ch
}

func (c *HttpClient) getUserId() string {
	switch c.Config.Account.AuthType {
	case AuthType_Password:
		// TODO: 使用 username 和 password 登录, 获取 userId 和 userToken, 存储到本地
		return ""
	default:
		logx.Errorf("invalid auth type: %v", c.Config.Account.AuthType)
		os.Exit(1)
		return ""
	}
}

func (c *HttpClient) getUserToken() string {
	switch c.Config.Account.AuthType {
	case AuthType_Password:
		// TODO: 使用 username 和 password 登录, 获取 userId 和 userToken, 存储到本地
		return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4NDI1Mjk5ODEsImp0aSI6IjMifQ.nbeMXsmCN3xXAFahLA76tdy5iKBdLIyqo3VfH9VXGSQ"
	default:
		logx.Errorf("invalid auth type: %v", c.Config.Account.AuthType)
		os.Exit(1)
		return ""
	}
}
