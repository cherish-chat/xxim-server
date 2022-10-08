package xenv

import "os"

const (
	KeyIp2RegionXdbUrl  = "IP2REGION_XDB_URL"
	KeyEnableDebugToken = "ENABLE_DEBUG_TOKEN"
)

var (
	Ip2RegionXdbUrl  = "https://xxim-public-1312910328.cos.ap-guangzhou.myqcloud.com/ip2region.xdb"
	EnableDebugToken = false
)

func init() {
	if v := os.Getenv(KeyIp2RegionXdbUrl); v != "" {
		Ip2RegionXdbUrl = v
	}
	if v := os.Getenv(KeyEnableDebugToken); v != "" {
		EnableDebugToken = true
	}
}
