package handler

import (
	"crypto/tls"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xhttp"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}
}

func WsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 代理到 conn 服务
		var connRpcIps []string
		for _, connPod := range svcCtx.ConnPodsMgr.AllConnServices() {
			connRpcIps = append(connRpcIps, connPod.PodIp)
		}
		if len(connRpcIps) == 0 {
			// 502
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		// randIp 随机选择一个 conn 服务实现负载均衡
		host := utils.AnyRandomInSlice(connRpcIps, "")
		ur := "http://" + host + "/ws"
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
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("X-Frame-Options")
			return nil
		}
		r.Header.Set("X-Real-IP", xhttp.GetRequestIP(r))
		proxy.ServeHTTP(w, r)
	}
}
