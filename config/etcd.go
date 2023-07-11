package config

import "github.com/zeromicro/go-zero/core/discov"

type EtcdConfig struct {
	Hosts              []string
	ID                 int64  `json:",optional"`
	User               string `json:",optional"`
	Pass               string `json:",optional"`
	CertFile           string `json:",optional"`
	CertKeyFile        string `json:",optional=CertFile"`
	CACertFile         string `json:",optional=CertFile"`
	InsecureSkipVerify bool   `json:",optional"`
}

func (c Config) GetEtcd(serviceName string) discov.EtcdConf {
	return discov.EtcdConf{
		Hosts:              c.Etcd.Hosts,
		Key:                serviceName,
		ID:                 c.Etcd.ID,
		User:               c.Etcd.User,
		Pass:               c.Etcd.Pass,
		CertFile:           c.Etcd.CertFile,
		CertKeyFile:        c.Etcd.CertKeyFile,
		CACertFile:         c.Etcd.CACertFile,
		InsecureSkipVerify: c.Etcd.InsecureSkipVerify,
	}
}
