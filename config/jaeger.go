package config

import "github.com/zeromicro/go-zero/core/trace"

type JaegerConfig struct {
	Endpoint string  `json:",optional"`
	Sampler  float64 `json:",default=1.0"`
	Disabled bool    `json:",optional"`
}

func (c Config) GetJaeger(serviceName string) trace.Config {
	if c.Jaeger.Disabled {
		return trace.Config{
			Name:         "",
			Endpoint:     "",
			Sampler:      0,
			Batcher:      "",
			OtlpHeaders:  nil,
			OtlpHttpPath: "",
			Disabled:     true,
		}
	}
	return trace.Config{
		Name:     serviceName,
		Endpoint: c.Jaeger.Endpoint,
		Sampler:  c.Jaeger.Sampler,
		Batcher:  "jaeger",
		Disabled: false,
	}
}
