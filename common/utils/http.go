package utils

import (
	"io"
	"net/http"
	"strings"
)

type xHttp struct {
}

var Http = &xHttp{}

func (x *xHttp) GetClientIP(r *http.Request) string {
	// 先取 X-Real-IP
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 取 X-Forwarded-For
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			// 取 RemoteAddr
			ip = r.RemoteAddr
		}
	}
	return strings.Split(ip, ",")[0]
}

func (x *xHttp) GetResponseBody(response *http.Response) []byte {
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return body
}
