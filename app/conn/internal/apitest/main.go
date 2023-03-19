package main

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	r "github.com/imroc/req/v3"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	loginReq := &pb.LoginReq{
		Id:       "test123456",
		Password: utils.Md5("123456"),
	}
	data, _ := proto.Marshal(loginReq)
	commonReq := &pb.CommonReq{
		UserId:      "",
		Token:       "",
		DeviceModel: "",
		DeviceId:    utils.GenId(),
		OsVersion:   "",
		Platform:    "android",
		PackageId:   utils.GenId(),
		AppVersion:  "",
		Language:    "",
		Data:        data,
		Ip:          "",
		UserAgent:   "",
	}
	data, _ = proto.Marshal(commonReq)
	response, err := r.R().SetBody(data).Post("https://api.cherish.chat/v1/user/white/login")
	if err != nil {
		log.Fatalf("request error: %s", err.Error())
	}
	if response.StatusCode != 200 {
		log.Fatalf("request error: %d", response.StatusCode)
	}
	
}

func Request(
	path string,
	req conngateway.IReq,
	resp conngateway.IResp,
) error {
	return nil
}
