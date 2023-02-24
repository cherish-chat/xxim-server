package main

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/config"
	"github.com/cherish-chat/xxim-server/sdk/conn"
	"github.com/cherish-chat/xxim-server/sdk/handler"
	"github.com/cherish-chat/xxim-server/sdk/svc"
	"log"
)

func main() {
	conf := config.Config{
		Client: conn.Config{
			Addr: "wss://api.cherish.chat:443/ws",
			DeviceConfig: conn.DeviceConfig{
				PackageId:   utils.GenId(),
				Platform:    "macos",
				DeviceId:    utils.GenId(),
				DeviceModel: "macos golang",
				OsVersion:   "10.15.7",
				AppVersion:  "v1.0.0",
				Language:    "zh",
				NetworkUsed: "wifi",
				Ext:         nil,
			},
			UserConfig: conn.UserConfig{
				UserId:   "test123456",
				Password: utils.Md5("123456"),
				Token:    "",
				Ext:      nil,
			},
		},
	}
	svcCtx := svc.NewServiceContext(conf)
	svcCtx.SetEventHandler(handler.NewEventHandler(svcCtx))

	err := svcCtx.Client().Connect()
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	select {}
}
