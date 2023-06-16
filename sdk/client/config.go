package client

import (
	"errors"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/sdk/store"
	"runtime"
	"strconv"
	"time"
)

type Config struct {
	// Endpoints are the endpoints of the servers.
	Endpoints []string
	// Encoding is the content type of the request. Please select one of the following: protobuf, json; default is protobuf.
	Encoding *pb.EncodingProto
	// AppId is the id of the application. if you don't have one, please ignore it. See https://console.imcloudx.com/#/app for details.
	AppId string
	// InstallId is the id of the installation. If not filled in, it will be written locally.
	InstallId string
	// Platform is the platform of the installation. optional: 0(ios), 1(android), 2(web), 3(windows), 4(mac), 5(linux), 6(ipad), 7(androidPad). If not filled in, the native platform value will be used.
	Platform *pb.Platform
	// DeviceModel is the model of the device. If not filled in, the native device model value will be used.
	DeviceModel string
	// OsVersion is the version of the operating system. If not filled in, the native operating system version value will be used.
	OsVersion string
	// Language is the language of the device. Default is zh-CN.
	Language *pb.I18NLanguage
	// RequestTimeout is the timeout of the request. Default is 10s.
	RequestTimeout time.Duration
	// UserToken is the token of the user. Required.
	UserToken string
	// CustomHeader is the custom header of the request.
	CustomHeader string
	// KeepAliveSecond is the keep alive second of the websocket. Default is 30s.
	KeepAliveSecond time.Duration
}

var (
	ErrNoEndpoint = errors.New("no endpoint available")
)

func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return ErrNoEndpoint
	}
	if c.Encoding == nil {
		encoding := pb.EncodingProto(1)
		c.Encoding = &encoding
	}
	if c.InstallId == "" {
		//读取本地配置
		installId := store.Database.Config.FindByK("install_id")
		if installId == "" {
			installId = utils.Snowflake.String()
			store.Database.Config.Save("install_id", installId)
		}
		c.InstallId = installId
	}
	if c.Platform == nil {
		//获取本机系统
		goos := runtime.GOOS
		switch goos {
		case "darwin":
			platform := pb.Platform(4)
			c.Platform = &platform
		case "linux":
			platform := pb.Platform(5)
			c.Platform = &platform
		case "windows":
			platform := pb.Platform(3)
			c.Platform = &platform
		default:
			platform := pb.Platform(2)
			c.Platform = &platform
		}
	}
	if c.DeviceModel == "" {
		// 读取本机设备信息
		goos := runtime.GOOS
		goarch := runtime.GOARCH
		numCpu := runtime.NumCPU()
		c.DeviceModel = goos + "-" + goarch + "-" + strconv.Itoa(numCpu)
	}
	if c.OsVersion == "" {
		// 读取本机系统版本
		c.OsVersion = runtime.Version()
	}
	if c.Language == nil {
		// 读取本机语言
		language := pb.I18NLanguage_Chinese_Simplified
		c.Language = &language
	}
	if c.UserToken == "" {
		return errors.New("user token is required")
	}
	if c.RequestTimeout == 0 {
		c.RequestTimeout = 10 * time.Second
	}
	if c.KeepAliveSecond == 0 {
		c.KeepAliveSecond = 30 * time.Second
	}
	return nil
}
