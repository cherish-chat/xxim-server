package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherish-chat/xxim-server/common"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/types"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"os"
)

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

var (
	ErrInvalidRequestType = errors.New("invalid request type, req must implement types.ReqInterface")
)

func (c *HttpClient) getUrl(path string) string {
	endpoint := c.Config.Endpoints[c.latestEndpointIndex]
	return fmt.Sprintf("%s%s", endpoint, path)
}

func (c *HttpClient) Request(path string, req any, resp any) error {
	var body io.Reader
	if req != nil {
		var data []byte
		var err error
		if c.Config.ContentType == "protobuf" {
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
			data, _ = proto.Marshal(&types.GatewayApiRequest{
				Header: &types.RequestHeader{
					RequestId:    utils.Snowflake.String(),
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
					AesKey:       c.Config.AesKey,
					AesIv:        c.Config.AesIv,
					Encoding:     types.EncodingProto_PROTOBUF,
				},
				Body: data,
			})
		} else {
			data, err = json.Marshal(req)
			if err != nil {
				return fmt.Errorf("req marshal error: %v", err)
			}
			data, _ = json.Marshal(&types.GatewayApiRequest{
				Header: &types.RequestHeader{
					RequestId:    utils.Snowflake.String(),
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
					AesKey:       c.Config.AesKey,
					AesIv:        c.Config.AesIv,
					Encoding:     types.EncodingProto_JSON,
				},
				Body: data,
			})
		}
		// 是否开启了加密
		if c.Config.EnableEncrypted {
			data = utils.Aes.Encrypt(c.Config.AesKey, c.Config.AesIv, data)
		}
		body = bytes.NewReader(data)
	}
	request, err := http.NewRequest("POST", c.getUrl(path), body)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}
	// set content type
	if c.Config.ContentType == "protobuf" {
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
	// 是否开启了加密
	if c.Config.EnableEncrypted {
		data, err = utils.Aes.Decrypt(c.Config.AesKey, c.Config.AesIv, data)
		if err != nil {
			return fmt.Errorf("decrypt response body error: %v", err)
		}
	}
	// unmarshal response body
	if resp != nil {
		if c.Config.ContentType == "protobuf" {
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
		return ""
	default:
		logx.Errorf("invalid auth type: %v", c.Config.Account.AuthType)
		os.Exit(1)
		return ""
	}
}
