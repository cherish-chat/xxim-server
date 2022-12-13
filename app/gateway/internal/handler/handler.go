package handler

import (
	"crypto/tls"
	"errors"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/wrapper"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xhttp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}
}

var wsProxyLogger = log.New(os.Stdout, "【proxy-error】", log.LstdFlags)

func WsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, err := getConnPodHost(svcCtx)
		if err != nil {
			// 502
			logx.WithContext(r.Context()).Errorf("get conn pod host error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		ur := "http://" + host + "/ws"
		if svcCtx.Config.EnableSSL {
			ur = "https://" + host + "/ws"
		}
		target, _ := url.Parse(ur)
		var transport http.RoundTripper = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 3 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			DisableCompression:    true,
		}
		proxy := xhttp.NewSingleHostReverseProxy(target)
		proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("X-Frame-Options")
			return nil
		}
		proxy.ErrorLog = wsProxyLogger
		r.Header.Set("X-Real-IP", xhttp.GetRequestIP(r))
		proxy.ServeHTTP(w, r)
	}
}

// 自定义负载均衡
func getConnPodHost(svcCtx *svc.ServiceContext) (string, error) {
	// 代理到 conn 服务
	var connRpcIps []string
	for _, connPod := range svcCtx.ConnPodsMgr.AllConnServices() {
		connRpcIps = append(connRpcIps, connPod.PodIpPort)
	}
	if len(connRpcIps) == 0 {
		return "", errors.New("no conn service")
	}
	// randIp 随机选择一个 conn 服务实现负载均衡
	return utils.AnyRandomInSlice(connRpcIps, ""), nil
}

func AuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requester := &pb.CommonReq{}
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		err := proto.Unmarshal(body, requester)
		if err != nil {
			wrapper.RequestValidateErr(w, err.Error())
			return
		}
		requester.Ip = xhttp.GetRequestIP(r)
		requester.UserAgent = r.UserAgent()
		resp := logic.NewAuthLogic(r, svcCtx).Auth(requester)
		if resp.Code != pb.CommonResp_Success {
			wrapper.AuthError(w, resp)
			return
		}
		wrapper.Success(w, nil, resp)
	}
}

func ShieldHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			ConvId string `path:"convId"`
		}{}
		err := httpx.Parse(r, req)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("parse request error: %v", err)
			wrapper.RequestValidateErr(w, err.Error())
			return
		}
		if req.ConvId == "" {
			logx.WithContext(r.Context()).Errorf("convId is empty")
			wrapper.RequestValidateErr(w, "convId is required")
			return
		}
		svgResp, err := logic.NewShieldLogic(r, svcCtx).Shield(req.ConvId)
		if err != nil {
			wrapper.InternalErr(w, err)
			return
		}
		// 返回 svg
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = w.Write([]byte(svgResp))
	}
}

func Cors(svcCtx *svc.ServiceContext) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, User-Agent, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, X-Forwarded-Host, X-Forwarded-Port, X-Forwarded-Server")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			next(w, r)
		}
	}
}
