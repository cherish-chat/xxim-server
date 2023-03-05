package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/propagation"
	"strconv"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(in *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	var fs []func()
	var failedConnParams []*pb.ConnParam
	var successConnParams []*pb.ConnParam
	xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendMsgToConnection", func(ctx context.Context) {
		for _, pod := range l.svcCtx.ConnPodsMgr.AllConnServices() {
			podValue := *pod
			fs = append(fs, func() {
				resp, err := podValue.SendMsg(ctx, in)
				if err != nil {
					l.Errorf("SendMsg error: %v", err)
					return
				}
				l.Debugf("resp.SuccessConnParams.length: %v", len(resp.SuccessConnParams))
				l.Debugf("resp.FailedConnParams.length: %v", len(resp.FailedConnParams))
				failedConnParams = append(failedConnParams, resp.FailedConnParams...)
				successConnParams = append(successConnParams, resp.SuccessConnParams...)
			})
		}
		mr.FinishVoid(fs...)
	}, propagation.MapCarrier{
		"userIds.length": strconv.Itoa(len(in.GetUserConnReq.UserIds)),
		"event":          in.Event.String(),
	})

	return &pb.SendMsgResp{
		CommonResp:        pb.NewSuccessResp(),
		SuccessConnParams: successConnParams,
		FailedConnParams:  failedConnParams,
	}, nil
}
