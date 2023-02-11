package mgmtmodel

import "github.com/cherish-chat/xxim-server/common/pb"

type OperationLog struct {
	Id string `gorm:"column:id;primarykey;comment:'主键'"`
	// 请求时间
	ReqTime    int64  `gorm:"column:reqTime;not null;comment:'请求时间';index;"`
	ReqTimeStr string `gorm:"column:reqTimeStr;not null;comment:'请求时间'"`
	// 响应时间
	RespTime    int64  `gorm:"column:respTime;not null;comment:'响应时间'"`
	RespTimeStr string `gorm:"column:respTimeStr;not null;comment:'响应时间'"`
	// 操作类型
	OperationType string `gorm:"column:operationType;not null;default:'';comment:'操作类型';index;"` // add:新增 update:修改 delete:删除 other:其他
	// 操作标题
	OperationTitle string `gorm:"column:operationTitle;not null;default:'';comment:'操作标题';index;"`
	// 请求路径
	ReqPath string `gorm:"column:reqPath;not null;default:'';comment:'请求路径'index;"`
	// 请求参数
	ReqParams string `gorm:"column:reqParams;comment:'请求参数';type:text;"`
	// 请求结果是否成功
	ResultSuccess bool `gorm:"column:resultSuccess;not null;default:0;comment:'请求结果是否成功'index;"`
	// 请求结果
	ReqResultMsg string `gorm:"column:reqResultMsg;not null;default:'';comment:'请求结果'"`
	// 响应
	Resp string `gorm:"column:resp;comment:'响应';type:text;"`
	// 请求ip
	ReqIp string `gorm:"column:reqIp;not null;default:'';comment:'请求ip'index;"`
	// ip来源
	IpSource string `gorm:"column:ipSource;not null;default:'';comment:'ip来源'"`
	// 请求耗时
	ReqCost int64 `gorm:"column:reqCost;not null;comment:'请求耗时'"`
	// 操作人
	Operator string `gorm:"column:operator;not null;default:'';comment:'操作人';index;"`
}

func (m *OperationLog) TableName() string {
	return MGMT_TABLE_PREFIX + "operationlog"
}

func (m *OperationLog) ToPB() *pb.MSOperationLog {
	return &pb.MSOperationLog{
		Id:             m.Id,
		ReqTime:        m.ReqTime,
		ReqTimeStr:     m.ReqTimeStr,
		RespTime:       m.RespTime,
		RespTimeStr:    m.RespTimeStr,
		OperationType:  m.OperationType,
		OperationTitle: m.OperationTitle,
		ReqPath:        m.ReqPath,
		ReqParams:      m.ReqParams,
		ResultSuccess:  m.ResultSuccess,
		ReqResultMsg:   m.ReqResultMsg,
		Resp:           m.Resp,
		ReqIp:          m.ReqIp,
		IpSource:       m.IpSource,
		ReqCost:        m.ReqCost,
		Operator:       m.Operator,
	}
}
