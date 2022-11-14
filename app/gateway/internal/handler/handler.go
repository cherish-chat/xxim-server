package handler

import (
	"crypto/tls"
	"errors"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xhttp"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
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
		if r.TLS != nil {
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
		proxy := httputil.NewSingleHostReverseProxy(target)
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
		connRpcIps = append(connRpcIps, connPod.PodIp)
	}
	if len(connRpcIps) == 0 {
		return "", errors.New("no conn service")
	}
	// randIp 随机选择一个 conn 服务实现负载均衡
	return utils.AnyRandomInSlice(connRpcIps, ""), nil
}
