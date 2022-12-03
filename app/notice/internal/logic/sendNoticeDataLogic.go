package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNoticeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendNoticeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNoticeDataLogic {
	return &SendNoticeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendNoticeData 发送通知数据
func (l *SendNoticeDataLogic) SendNoticeData(in *pb.SendNoticeDataReq) (*pb.SendNoticeDataResp, error) {
	m := &noticemodel.Notice{
		NoticeId:    utils.If(in.NoticeData.NoticeId != "", in.NoticeData.NoticeId, utils.GenId()),
		ConvId:      in.NoticeData.ConvId,
		CreateTime:  utils.If(in.NoticeData.CreateTime != "", utils.AnyToInt64(in.NoticeData.CreateTime), time.Now().UnixMilli()),
		Title:       in.NoticeData.Title,
		ContentType: in.NoticeData.ContentType,
		Content:     in.NoticeData.Content,
		Options: noticemodel.NoticeOption{
			StorageForClient: in.NoticeData.Options.StorageForClient,
			UpdateConvMsg:    in.NoticeData.Options.UpdateConvMsg,
			OnlinePushOnce:   in.NoticeData.Options.OnlinePushOnce,
		},
		IsBroadcast: in.GetIsBroadcast(),
		Ext:         in.NoticeData.Ext,
		UserId:      in.GetUserId(),
	}
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := m.Upsert(tx)
		if err != nil {
			l.Errorf("upsert notice error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.SendNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "PushNoticeData", func(ctx context.Context) {
		for i := 0; i < 12; i++ {
			resp, err := NewPushNoticeDataLogic(ctx, l.svcCtx).PushNoticeData(&pb.PushNoticeDataReq{
				CommonReq: in.CommonReq,
				NoticeId:  m.NoticeId,
			})
			if err == nil {
				break
			}
			if resp.CommonResp == nil {
				break
			}
			if !resp.CommonResp.Failed() {
				break
			}
			logx.WithContext(ctx).Errorf("push notice data error: %v, resp: %v", err, utils.AnyToString(resp))
			time.Sleep(2 * time.Second)
		}
	}, propagation.MapCarrier{})
	return &pb.SendNoticeDataResp{}, nil
}
