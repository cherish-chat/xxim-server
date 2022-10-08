package ip2region

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xenv"
	. "github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var (
	cBuff     []byte
	_searcher *Searcher
)

type IpRegion struct {
	Country  string
	District string
	Province string
	City     string
	Isp      string
}

func init() {
	var err error
	// 判断是否存在本地文件
	if !utils.FileExists("./ip2region.xdb") {
		err = utils.DownloadFile(xenv.Ip2RegionXdbUrl, "./ip2region.xdb")
		if err != nil {
			logx.Errorf("download ip2region.xdb failed: %s", err.Error())
			panic(err)
		}
	}
	cBuff, err = LoadContentFromFile("./ip2region.xdb")
	if err != nil {
		logx.Errorf("load ip2region.xdb failed: %s", err.Error())
		panic(err)
	}
	_searcher, err = NewWithBuffer(cBuff)
	if err != nil {
		logx.Errorf("new ip2region searcher failed: %s", err.Error())
		panic(err)
	}
}

func Ip2Region(ip string) string {
	if ip == "" {
		return ""
	}
	if len(strings.Split(ip, ",")) > 1 {
		ip = strings.Split(ip, ",")[0]
	}
	str, _ := _searcher.SearchByStr(ip)
	return str
}

func RegionSplit(ip string) IpRegion {
	region := Ip2Region(ip)
	split := strings.Split(region, "|")
	if len(split) == 5 {
		return IpRegion{split[0], split[2], split[3], split[1], split[4]}
	}
	return IpRegion{"", "", "", "", ""}
}
