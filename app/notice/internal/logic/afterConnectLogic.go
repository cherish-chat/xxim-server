package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/propagation"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AfterConnect conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	return l.afterConnect(in)
}

func (l *AfterConnectLogic) afterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	// 查询用户所有订阅的会话 检测是否有未消费的消息 进行推送
	var convIds []string
	var err error
	convIds, err = l.getAllConv(in)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	// 并发查询未消费的消息
	var fs []func() error
	for _, convId := range convIds {
		convId := convId
		fs = append(fs, func() error {
			var err error
			xtrace.StartFuncSpan(l.ctx, "getUserNoticeData", func(ctx context.Context) {
				_, err = NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
					CommonReq: &pb.CommonReq{
						UserId: in.ConnParam.UserId,
					},
					UserId: in.ConnParam.UserId,
					ConvId: convId,
				})
			}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
				"conv_id": convId,
			}))
			if err != nil {
				l.Errorf("get user notice data error: %v", err)
			}
			return err
		})
	}
	err = mr.Finish(fs...)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	return pb.NewSuccessResp(), nil
}

func (l *AfterConnectLogic) getAllConv(in *pb.AfterConnectReq) ([]string, error) {
	var userId = in.ConnParam.UserId
	var convIds []string
	convIdOfUser, err := l.svcCtx.ImService().GetAllConvIdOfUser(l.ctx, &pb.GetAllConvIdOfUserReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("get all conv id of user error: %v", err)
		return convIds, err
	}
	convIds = convIdOfUser.ConvIds
	return convIds, nil
}
