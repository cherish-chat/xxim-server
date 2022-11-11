package immodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type UserConnectRecord struct {
	UserId         string            `json:"userId" bson:"userId"`
	DeviceId       string            `json:"deviceId" bson:"deviceId"`
	Platform       string            `json:"platform" bson:"platform"`
	Ips            string            `json:"ips" bson:"ips"`
	IpRegion       ip2region.Obj     `json:"ipRegion" bson:"ipRegion"`
	NetworkUsed    string            `json:"networkUsed" bson:"networkUsed"`
	Headers        map[string]string `json:"headers" bson:"headers"`
	PodIp          string            `json:"podIp" bson:"podIp"`
	ConnectTime    int64             `json:"connectTime" bson:"connectTime"`
	DisconnectTime int64             `json:"disconnectTime" bson:"disconnectTime"`
}

func (m *UserConnectRecord) CollectionName() string {
	return "user_connect_record"
}

func (m *UserConnectRecord) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key: []string{"userId", "deviceId"},
	}, {
		Key: []string{"userId"},
	}, {
		Key: []string{"connectTime"},
	}, {
		Key: []string{"disconnectTime"},
	}})
	return nil
}
