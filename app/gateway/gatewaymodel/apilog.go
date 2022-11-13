package gatewaymodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
)

type ApiLog struct {
	Service         string             `json:"service" bson:"service"`                 // 业务
	Req             string             `json:"req" bson:"req"`                         // 请求
	Resp            string             `json:"resp" bson:"resp"`                       // 响应
	Err             string             `json:"err" bson:"err"`                         // 错误
	RespCode        pb.CommonResp_Code `json:"respCode" bson:"respCode"`               // 响应码
	Requester       *pb.Requester      `json:"requester" bson:"requester"`             // 请求者
	IpRegion        *pb.IpRegion       `json:"ipRegion" bson:"ipRegion"`               // ip 地区
	RequestTime     int64              `json:"requestTime" bson:"requestTime"`         // 请求时间
	ResponseTime    int64              `json:"responseTime" bson:"responseTime"`       // 响应时间
	RequestTimeStr  string             `json:"requestTimeStr" bson:"requestTimeStr"`   // 请求时间
	ResponseTimeStr string             `json:"responseTimeStr" bson:"responseTimeStr"` // 响应时间
	TraceId         string             `json:"traceId" bson:"traceId"`                 // traceId
}

func (m *ApiLog) CollectionName() string {
	return "api_log"
}

func (m *ApiLog) Indexes(c *qmgo.Collection) error {
	return c.CreateIndexes(context.Background(), []opts.IndexModel{{
		Key: []string{"service"},
	}, {
		Key: []string{"err"},
	}, {
		Key: []string{"http_code"},
	}, {
		Key: []string{"resp_code"},
	}, {
		Key: []string{"requester.id"},
	}, {
		Key: []string{"requester.deviceId"},
	}, {
		Key: []string{"requester.ip"},
	}, {
		Key: []string{"ip_region.province"},
	}, {
		Key: []string{"ip_region.city"},
	}, {
		Key: []string{"requestTime"},
	}, {
		Key: []string{"responseTime"},
	}})
}
