package ip2region

import (
	"github.com/cherish-chat/xxim-server/common/pb"
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
	// 判断本地是否存在 ./ip2region.xdb
	var err error
	if !utils.FileExists("./ip2region.xdb") {
		err = utils.DownloadFile(downloadUrl, "./ip2region.xdb")
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

type Obj struct {
	Country  string `json:"country" bson:"country"`
	District string `json:"district" bson:"district"`
	Province string `json:"province" bson:"province"`
	City     string `json:"city" bson:"city"`
	ISP      string `json:"isp" bson:"isp"`
}

func (o Obj) String() string {
	return o.Country + "|" + o.District + "|" + o.Province + "|" + o.City + "|" + o.ISP
}

func (o Obj) Pb() *pb.IpRegion {
	return &pb.IpRegion{
		Country:  o.Country,
		Province: o.Province,
		City:     o.City,
		Isp:      o.ISP,
	}
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
