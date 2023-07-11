package config

import "github.com/zeromicro/go-zero/core/logx"

type LogConfig struct {
	Mode  string `json:",default=console,options=[console,file,volume]"`
	Level string `json:",default=info,options=[debug,info,error,severe]"`
}

func (c Config) GetLog(serviceName string) logx.LogConf {
	return logx.LogConf{
		ServiceName:         serviceName,
		Mode:                c.Log.Mode,
		Encoding:            "json",
		Path:                "logs",
		Level:               c.Log.Level,
		MaxContentLength:    1024,
		Compress:            false,
		Stat:                false,
		KeepDays:            7,
		StackCooldownMillis: 100,
		Rotation:            "daily",
	}
}
