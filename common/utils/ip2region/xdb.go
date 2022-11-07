package ip2region

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	. "github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var (
	cBuff     []byte
	_searcher *Searcher
)

func Init(downloadUrl string) {
	err := utils.DownloadFile(downloadUrl, "./ip2region.xdb")
	if err != nil {
		logx.Errorf("download ip2region.xdb failed: %s", err.Error())
		panic(err)
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

type Obj struct {
	Country  string
	District string
	Province string
	City     string
	ISP      string
}

func (o Obj) String() string {
	return o.Country + "|" + o.District + "|" + o.Province + "|" + o.City + "|" + o.ISP
}

func Ip2Region(ip string) Obj {
	if ip == "" {
		return Obj{}
	}
	split := strings.Split(ip, ",")
	for _, s := range split {
		str, _ := _searcher.SearchByStr(s)
		if str != "" {
			regionSplit := strings.Split(str, "|")
			if len(regionSplit) == 5 {
				return Obj{regionSplit[0], regionSplit[1], regionSplit[2], regionSplit[3], regionSplit[4]}
			}
		}
	}
	return Obj{}
}
