package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/gatewaymodel"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/propagation"
	"time"
)

type ApiLogLogic struct {
	_ctx    context.Context
	svcCtx  *svc.ServiceContext
	traceId string

	logx.Logger
}

func NewApiLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiLogLogic {
	return &ApiLogLogic{_ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx), traceId: xtrace.TraceIdFromContext(ctx)}
}

func (l *ApiLogLogic) ApiLog(requester *pb.CommonReq, service string, commonResp *pb.CommonResp, req string, resp string, requestTime time.Time, responseTime time.Time, err error) {
	defer func() {
		if r := recover(); r != nil {
			l.Error("ApiLog panic", r)
		}
	}()
	if commonResp == nil {
		commonResp = pb.NewSuccessResp()
	}
	var ipRegion = &pb.IpRegion{}
	if requester != nil {
		if requester.Ip != "" {
			ipRegion = ip2region.Ip2Region(requester.Ip).Pb()
		}
	}
	attr := make(map[string]string)
	if requester != nil {
		attr["userId"] = requester.Id
		attr["appVersion"] = requester.AppVersion
		attr["ip"] = requester.Ip
		attr["deviceId"] = requester.DeviceId
		attr["deviceModel"] = requester.DeviceModel
		attr["platform"] = requester.Platform
		attr["osVersion"] = requester.OsVersion
		attr["userAgent"] = requester.UserAgent
		attr["ipRegion.country"] = ipRegion.Country
		attr["ipRegion.province"] = ipRegion.Province
		attr["ipRegion.city"] = ipRegion.City
		attr["ipRegion.isp"] = ipRegion.Isp
	}
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	xtrace.RunWithTrace(l.traceId, "InsertApiLog", func(ctx context.Context) {
		apiLog := &gatewaymodel.ApiLog{
			Service:         service,
			Req:             req,
			Resp:            resp,
			Err:             errStr,
			RespCode:        commonResp.Code,
			CommonReq:       utils.AnyToString(requester),
			IpRegion:        utils.AnyToString(ipRegion),
			RequestTime:     requestTime.UnixMilli(),
			ResponseTime:    responseTime.UnixMilli(),
			RequestTimeStr:  requestTime.Format("2006-01-02 15:04:05.000"),
			ResponseTimeStr: responseTime.Format("2006-01-02 15:04:05.000"),
			TraceId:         l.traceId,
		}
		err = xorm.InsertOne(l.svcCtx.Mysql(), apiLog)
		if err != nil {
			l.Errorf("ApiLog err: %v", err)
		}
	}, propagation.MapCarrier(attr))
}
