package conn

import (
	"net/http"
)

type Header struct {
	Key   string
	Value string
}
type DeviceConfig struct {
	PackageId   string
	Platform    string
	DeviceId    string
	DeviceModel string
	OsVersion   string
	AppVersion  string
	Language    string
	NetworkUsed string
	Ext         []byte
}
type UserConfig struct {
	UserId   string
	Password string
	Token    string
	Ext      []byte
}
type Config struct {
	Addr         string   // ws://xxx.xxx.xxx.xxx:xxxx or wss://xxx.xxx.xxx.xxx:xxxx
	Headers      []Header // http header
	DeviceConfig DeviceConfig
	UserConfig   UserConfig
}

func (c Config) GetHeader() http.Header {
	var header http.Header
	for _, h := range c.Headers {
		header.Add(h.Key, h.Value)
	}
	return header
}
