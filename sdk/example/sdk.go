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
				UserId:   "testb4",
				Password: utils.Md5("123456"),
				Token:    "",
				Ext:      nil,
			},
			RsaPublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuiJyMvMqTKKC5Z4qWU3v
R9ZWm1JnEEP66xYC0a62XsNE+Vi/3OtKChrhsLGzzEfpAmYLtIdODK/Wm5VQeFqA
w/2UgWtIxPrKfLllA3tTcKkbw/K/9WkO24FKmPPg00L7OaVbfvg/0TorLnMyQ65R
OnG8fvs+LqrIRIDgGZIPGCytV4IdV988v/7KHLNUvyAoINLVIISriUwwr5cjAORL
RLsPVW0jJp4xNleE55Vi+0PlmloPwGtEt9xMRIaTIQzpgBzuLLymxF5a5ifbHg/V
xqDumvu1sYCot9fhDqktYsVz990FgpHJv7xeY11ZFvfKYl4T0VLg5Mvzq8+BX5ut
SQIDAQAB
-----END PUBLIC KEY-----
`,
			AesIv: "哈哈哈aaaaaaaasdsadsada",
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
