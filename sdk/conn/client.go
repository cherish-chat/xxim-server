package conn

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xaes"
	"github.com/cherish-chat/xxim-server/common/utils/xrsa"
	"github.com/cherish-chat/xxim-server/sdk/types"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type Client struct {
	Config       Config
	EventHandler types.EventHandler
	ctx          context.Context
	cancelFunc   context.CancelFunc
	ws           *websocket.Conn
	lock         sync.Mutex
	respMap      sync.Map
	ticker       *time.Ticker
	aesKey       []byte
	aesIv        []byte
}

func NewClient(config Config, eventHandler types.EventHandler) *Client {
	c := &Client{Config: config, EventHandler: eventHandler}
	// can cancel the context
	ctx := context.Background()
	ctxCancel, cancelFunc := context.WithCancel(ctx)
	c.ctx = ctxCancel
	c.cancelFunc = cancelFunc
	c.ticker = time.NewTicker(time.Second * 30)
	go c.timer()
	return c
}

func (c *Client) Close(code websocket.StatusCode, reason string) error {
	c.cancelFunc()
	// 关闭之前先判断ws是否存在
	ws := c.getWs()
	if ws != nil {
		c.EventHandler.BeforeClose(code, reason)
		defer c.EventHandler.AfterClose(code, reason)
		err := ws.Close(code, reason)
		c.setWs(nil)
		return err
	}
	c.respMap = sync.Map{}
	return nil
}

func (c *Client) ReConnect() error {
	c.EventHandler.BeforeReConnect()
	defer c.EventHandler.AfterReConnect()
	// 重连之前先关闭
	err := c.Close(websocket.StatusInternalError, "reconnect")
	if err != nil {
		return err
	}
	// 重连
	// 重新设置context
	ctx := context.Background()
	ctxCancel, cancelFunc := context.WithCancel(ctx)
	// 加锁
	c.lock.Lock()
	c.ctx = ctxCancel
	c.cancelFunc = cancelFunc
	c.lock.Unlock()
	return c.Connect()
}

func (c *Client) getWs() *websocket.Conn {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.ws
}

func (c *Client) setWs(ws *websocket.Conn) error {
	c.lock.Lock()
	c.ws = ws
	c.lock.Unlock()
	if ws == nil {
		return nil
	}
	// 读消息
	go c.readMessage()
	err := c.SetCxnParams()
	if err != nil {
		return err
	}
	err = c.SetUserParams()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Connect() error {
	c.EventHandler.BeforeConnect()
	defer c.EventHandler.AfterConnect()
	ws, _, err := websocket.Dial(c.ctx, c.Config.Addr, &websocket.DialOptions{
		HTTPHeader:      c.Config.GetHeader(),
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil {
		return err
	}
	err = c.setWs(ws)
	if err != nil {
		c.Close(websocket.StatusInternalError, "set ws error")
		return err
	}
	c.sync()
	return nil
}

func (c *Client) readMessage() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			ws := c.getWs()
			if ws == nil {
				return
			}
			typ, message, err := ws.Read(c.ctx)
			if err != nil {
				// close and return
				c.Close(websocket.StatusInternalError, "read message error")
				return
			}
			go c.EventHandler.OnMessage(typ, message)
			if typ == websocket.MessageBinary {
				// 解密
				if len(c.aesKey) > 0 && len(c.aesIv) > 0 {
					// aes解密
					var err error
					message, err = xaes.Decrypt(c.aesIv, c.aesKey, message)
					if err != nil {
						c.Close(websocket.StatusInternalError, "decrypt message error")
						return
					}
				}
				pushBody := &pb.PushBody{}
				err = proto.Unmarshal(message, pushBody)
				if err != nil {
					// close and return
					c.Close(websocket.StatusInternalError, "read message error")
					return
				}
				switch pushBody.Event {
				case pb.PushEvent_PushMsgDataList:
					go c.EventHandler.OnPushMsgDataList(pushBody)
				case pb.PushEvent_PushNoticeData:
					go c.EventHandler.OnPushNoticeData(pushBody)
				case pb.PushEvent_PushResponseBody:
					// 解析消息
					respBody := &pb.ResponseBody{}
					err = proto.Unmarshal(pushBody.Data, respBody)
					if err != nil {
						// close and return
						c.Close(websocket.StatusInternalError, "read message error")
						return
					}
					go c.EventHandler.OnPushResponseBody(pushBody)
					// 获取channel
					value, ok := c.respMap.Load(respBody.ReqId)
					if ok {
						ch := value.(chan *pb.ResponseBody)
						ch <- respBody
					}
				}
			}
		}
	}
}

func (c *Client) RequestX(
	method string,
	data proto.Message,
	response proto.Message,
) error {
	id := utils.GenId()
	dataBuff, _ := proto.Marshal(data)
	reqBody := &pb.RequestBody{
		ReqId:  id,
		Method: method,
		Data:   dataBuff,
	}
	dataBuff, _ = proto.Marshal(reqBody)
	ch := make(chan *pb.ResponseBody, 1)
	c.respMap.Store(id, ch)
	ws := c.getWs()
	if ws == nil {
		return errors.New("ws is nil")
	}
	err := c.write(method, ws, dataBuff)
	if err != nil {
		return err
	}
	// 5s ctx
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout")
		case valueBody := <-ch:
			c.respMap.Delete(id)
			if valueBody.Code == pb.ResponseBody_Success {
				err = proto.Unmarshal(valueBody.Data, response)
				if err != nil {
					return err
				} else {
					return nil
				}
			} else {
				return errors.New(valueBody.Code.String())
			}
		}
	}
}

func (c *Client) SetCxnParams() error {
	c.aesKey = nil
	c.aesIv = nil
	resp := &pb.SetCxnParamsResp{}
	var aesKeyEncrypted []byte
	var aesIvEncrypted []byte
	var err error
	aesKey := []byte(utils.GenId())
	aesIv := []byte(utils.GenId())
	if c.Config.RsaPublicKey != "" {
		bytes, err := xrsa.Encrypt(aesKey, []byte(c.Config.RsaPublicKey))
		if err != nil {
			logx.Errorf("set cxn params error: %s", err.Error())
			return err
		}
		aesKeyEncrypted = bytes
		c.aesKey = []byte(utils.Md5Bytes(aesKey))

		bytes, err = xrsa.Encrypt(aesIv, []byte(c.Config.RsaPublicKey))
		if err != nil {
			logx.Errorf("set cxn params error: %s", err.Error())
			return err
		}
		aesIvEncrypted = bytes
		c.aesIv = []byte(utils.Md5Bytes16(aesIv))
	}
	err = c.RequestX("/v1/conn/white/setCxnParams", &pb.SetCxnParamsReq{
		PackageId:   c.Config.DeviceConfig.PackageId,
		Platform:    c.Config.DeviceConfig.Platform,
		DeviceId:    c.Config.DeviceConfig.DeviceId,
		DeviceModel: c.Config.DeviceConfig.DeviceModel,
		OsVersion:   c.Config.DeviceConfig.OsVersion,
		AppVersion:  c.Config.DeviceConfig.AppVersion,
		Language:    c.Config.DeviceConfig.Language,
		NetworkUsed: c.Config.DeviceConfig.NetworkUsed,
		Ext:         c.Config.DeviceConfig.Ext,
		AesKey:      aesKeyEncrypted,
		AesIv:       aesIvEncrypted,
	}, resp)
	if err != nil {
		logx.Errorf("set cxn params error: %s", err.Error())
	}
	return err
}

func (c *Client) SetUserParams() error {
	token := c.Config.UserConfig.Token
	if token == "" {
		loginResp, err := c.LoginByPassword()
		if err != nil {
			return err
		}
		token = loginResp.Token
	}
	resp := &pb.SetUserParamsResp{}
	err := c.RequestX("/v1/conn/white/setUserParams", &pb.SetUserParamsReq{
		UserId: c.Config.UserConfig.UserId,
		Token:  token,
		Ext:    c.Config.UserConfig.Ext,
	}, resp)
	if err != nil {
		logx.Errorf("set user params error: %s", err.Error())
	}
	return err
}

func (c *Client) LoginByPassword() (*pb.LoginResp, error) {
	resp := &pb.LoginResp{}
	err := c.RequestX("/v1/user/white/login", &pb.LoginReq{
		Id:       c.Config.UserConfig.UserId,
		Password: c.Config.UserConfig.Password,
	}, resp)
	if err != nil {
		logx.Errorf("login error: %s", err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *Client) timer() {
	go c.EventHandler.OnTimer()
	go c.sync()
	for {
		select {
		case <-c.ticker.C:
			go c.EventHandler.OnTimer()
			go c.sync()
		}
	}
}

func (c *Client) sync() {
	getFriendListResp := &pb.GetFriendListResp{}
	err := c.RequestX("/v1/relation/getFriendList", &pb.GetFriendListReq{
		Page: &pb.Page{Page: 1, Size: 100000},
		Opt:  pb.GetFriendListReq_OnlyId,
	}, getFriendListResp)
	if err != nil {
		logx.Errorf("getFriendList error: %s", err.Error())
		return
	}
	getGroupListResp := &pb.GetMyGroupListResp{}
	err = c.RequestX("/v1/group/getMyGroupList", &pb.GetMyGroupListReq{
		Page: &pb.Page{Page: 1, Size: 100000},
		Opt:  pb.GetMyGroupListReq_ONLY_ID,
	}, getGroupListResp)
	if err != nil {
		logx.Errorf("getGroupList error: %s", err.Error())
		return
	}
	var convIds []string
	for _, v := range getFriendListResp.Ids {
		convIds = append(convIds, pb.SingleConvId(c.Config.UserConfig.UserId, v))
	}
	for _, v := range getGroupListResp.Ids {
		convIds = append(convIds, pb.GroupConvId(v))
	}
	convIds = utils.Set(convIds)
	if len(convIds) == 0 {
		return
	}
	resp := &pb.BatchGetConvSeqResp{}
	err = c.RequestX("/v1/msg/batchGetConvSeq", &pb.BatchGetConvSeqReq{
		ConvIdList: convIds,
		CommonReq:  nil,
	}, resp)
	if err != nil {
		logx.Errorf("batchGetConvSeq error: %s", err.Error())
		return
	}
}

func (c *Client) write(method string, ws *websocket.Conn, dataBuff []byte) error {
	if method != "/v1/conn/white/setCxnParams" {
		if len(c.aesKey) > 0 && len(c.aesIv) > 0 {
			dataBuff = xaes.Encrypt(c.aesIv, c.aesKey, dataBuff)
		}
	}
	err := ws.Write(c.ctx, websocket.MessageBinary, dataBuff)
	if err != nil {
		logx.Errorf("write error: %s", err.Error())
		return err
	}
	return nil
}
