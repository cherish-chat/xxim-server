package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/zeromicro/go-zero/core/netx"
	"os"
)

const (
	envPodIp = "POD_IP"
)

func Md5(s string) string {
	return Md5Bytes([]byte(s))
}

func Md5Bytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

func Md516(s string) string {
	// 将中间的第9位到第24位提取出来
	return Md5(s)[8:24]
}

func GetPodIp() string {
	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	return ip
}
