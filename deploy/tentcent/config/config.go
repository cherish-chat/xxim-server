package config

type TDMQTopic struct {
	Partition uint64
}
type StatefulSet struct {
	MaxPod int
}
type Deployment struct {
	Port int
}

type Config struct {
	Tencent struct {
		SecretId  string // 腾讯云 SecretId
		SecretKey string // 腾讯云 SecretKey
		Region    string // 腾讯云 地域 ap-guangzhou/ap-shanghai/ap-beijing/ap-hongkong ...
		TDMQ      struct {
			ClusterName string // 集群名称
			Namespace   string // 命名空间
			Topics      struct {
				Msg TDMQTopic
			}
		}
	}
	StatefulSets struct {
		Conn StatefulSet
	}
	Deployments struct {
		Msg Deployment
	}
}
