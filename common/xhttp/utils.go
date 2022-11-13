package xhttp

import (
	"net/http"
	"strings"
)

func GetRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		split := strings.Split(ip, ":")
		ip = ""
		for i, s := range split {
			if i != len(split)-1 {
				ip += s + ":"
			}
		}
		ip = strings.TrimSuffix(ip, ":")
	}
	if ip == "[::1]" {
		ip = "127.0.0.1"
	}
	return ip
}
