package gatewaymodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
)

type ApiLog struct {
	Service         string             `json:"service" bson:"service" gorm:"column:service;index;type:char(64);"`                   // 业务
	Req             string             `json:"req" bson:"req" gorm:"column:req;type:text;"`                                         // 请求
	Resp            string             `json:"resp" bson:"resp" gorm:"column:resp;type:text;"`                                      // 响应
	Err             string             `json:"err" bson:"err" gorm:"column:err;index;type:varchar(255);"`                           // 错误
	RespCode        pb.CommonResp_Code `json:"respCode" bson:"respCode" gorm:"column:respCode;index;"`                              // 响应码
	Requester       string             `json:"requester" bson:"requester" gorm:"column:requester;type:text;"`                       // 请求者
	IpRegion        string             `json:"ipRegion" bson:"ipRegion" gorm:"column:ipRegion;type:text;"`                          // ip 地区
	RequestTime     int64              `json:"requestTime" bson:"requestTime" gorm:"column:requestTime;index;type:bigint(13);"`     // 请求时间
	ResponseTime    int64              `json:"responseTime" bson:"responseTime" gorm:"column:responseTime;index;type:bigint(13);"`  // 响应时间
	RequestTimeStr  string             `json:"requestTimeStr" bson:"requestTimeStr" gorm:"column:requestTimeStr;type:char(32);"`    // 请求时间
	ResponseTimeStr string             `json:"responseTimeStr" bson:"responseTimeStr" gorm:"column:responseTimeStr;type:char(32);"` // 响应时间
	TraceId         string             `json:"traceId" bson:"traceId" gorm:"column:traceId;index;type:char(64);"`                   // traceId
}

func (m *ApiLog) TableName() string {
	return "api_log"
}
