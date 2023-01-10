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
	for _, pod := range l.svcCtx.ConnPodsMgr.AllConnServices() {
		podValue := *pod
		fs = append(fs, func() {
			xtrace.StartFuncSpan(l.ctx, "SendMsgToConnection", func(ctx context.Context) {
				resp, err := podValue.SendMsg(l.ctx, in)
				if err != nil {
					l.Errorf("SendMsg error: %v", err)
					return
				}
				l.Infof("resp.SuccessConnParams.length: %v", len(resp.SuccessConnParams))
				l.Infof("resp.FailedConnParams.length: %v", len(resp.FailedConnParams))
				failedConnParams = append(failedConnParams, resp.FailedConnParams...)
				successConnParams = append(successConnParams, resp.SuccessConnParams...)
			}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
				"userIds.length": strconv.Itoa(len(in.GetUserConnReq.UserIds)),
				"event":          in.Event.String(),
			}))
		})
	}
	mr.FinishVoid(fs...)

	return &pb.SendMsgResp{
		CommonResp:        pb.NewSuccessResp(),
		SuccessConnParams: successConnParams,
		FailedConnParams:  failedConnParams,
	}, nil
}
