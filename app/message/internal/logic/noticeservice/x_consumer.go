package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumerLogic {
	return &ConsumerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConsumerLogic) NoticeBatchSend(topic string, msg []byte) error {
	req := &pb.NoticeBatchSendReq{}
	err := utils.Json.Unmarshal(msg, req)
	if err != nil {
		l.Errorf("unmarshal message error: %v", err)
		return nil
	}
	var broadcastNotices []*noticemodel.BroadcastNotice
	for _, notice := range req.Notices {
		if notice.Broadcast {
			broadcastNotices = append(broadcastNotices, noticemodel.BroadcastNoticeFromPb(notice.Notice))
		}
	}
	if len(broadcastNotices) > 0 {
		// upsert
		bulk := l.svcCtx.BroadcastNoticeCollection.Bulk()
		for _, notice := range broadcastNotices {
			if notice.Sort == 0 {
				//生成一个唯一的sort
				sort, e := l.svcCtx.Redis.Incr(xcache.RedisVal.IncrKeyNoticeSort)
				if e != nil {
					l.Errorf("incr notice sort error: %v", e)
					return e
				}
				notice.Sort = sort
			}
			bulk.Upsert(bson.M{
				"conversationId":   notice.ConversationId,
				"conversationType": notice.ConversationType,
				"contentType":      notice.ContentType,
			}, notice)
		}
		_, err = bulk.Run(l.ctx)
		// 判断是不是唯一索引冲突
		if err != nil && !xmgo.ErrIsDup(err) {
			return err
		}
	}

	var subscriptionNoticeContents []*noticemodel.SubscriptionNoticeContent
	var subscriptionNotices []*noticemodel.SubscriptionNotice
	for _, notice := range req.Notices {
		if !notice.Broadcast {
			//判断会话类型
			if notice.Notice.ConversationType != pb.ConversationType_Subscription {
				l.Errorf("conversation type error: %v", notice.Notice.ConversationType)
				continue
			}
			contentId := utils.Snowflake.String()
			subscriptionNoticeContents = append(subscriptionNoticeContents, &noticemodel.SubscriptionNoticeContent{
				ContentId: contentId,
				Content:   notice.Notice.Content,
			})
			for _, userId := range notice.UserIds {
				subscriptionNotices = append(subscriptionNotices, &noticemodel.SubscriptionNotice{
					UserId:         userId,
					SubscriptionId: notice.Notice.ConversationId,
					ContentId:      contentId,
					UpdateTime:     primitive.DateTime(notice.Notice.UpdateTime),
					ContentType:    notice.Notice.ContentType,
				})
			}
		}
	}

	if len(subscriptionNoticeContents) > 0 {
		_, err = l.svcCtx.SubscriptionNoticeContentCollection.InsertMany(l.ctx, subscriptionNoticeContents)
		// 判断是不是唯一索引冲突
		if err != nil {
			return err
		}
	}

	if len(subscriptionNotices) > 0 {
		// upsert
		bulk := l.svcCtx.SubscriptionNoticeCollection.Bulk()

		for _, notice := range subscriptionNotices {
			if notice.Sort == 0 {
				//生成一个唯一的sort
				sort, e := l.svcCtx.Redis.Incr(xcache.RedisVal.IncrKeyNoticeSort)
				if e != nil {
					l.Errorf("incr notice sort error: %v", e)
					return e
				}
				notice.Sort = sort
			}
			bulk = bulk.Upsert(bson.M{
				"userId":         notice.UserId,
				"subscriptionId": notice.SubscriptionId,
				"contentType":    notice.ContentType,
			}, notice)
		}

		_, err = bulk.Run(l.ctx)
		if err != nil {
			l.Errorf("upsert subscription notice error: %v", err)
			return err
		}
	}

	go l.pushBroadcastNotice(broadcastNotices)
	go l.pushSubscriptionNotice(subscriptionNoticeContents, subscriptionNotices)

	return nil
}

func (l *ConsumerLogic) pushBroadcastNotice(broadcastNotices []*noticemodel.BroadcastNotice) {
	l.Debugf("push broadcast notice: %v", broadcastNotices)
}

func (l *ConsumerLogic) pushSubscriptionNotice(subscriptionNoticeContents []*noticemodel.SubscriptionNoticeContent, subscriptionNotices []*noticemodel.SubscriptionNotice) {
	for _, notice := range subscriptionNotices {
		content := ""
		for _, contentItem := range subscriptionNoticeContents {
			if contentItem.ContentId == notice.ContentId {
				content = contentItem.Content
				break
			}
		}
		pbNotice := &pb.Notice{
			NoticeId:         utils.AnyString(notice.Sort),
			ConversationId:   notice.SubscriptionId,
			ConversationType: pb.ConversationType_Subscription,
			Content:          content,
			ContentType:      notice.ContentType,
			UpdateTime:       int64(notice.UpdateTime),
			Sort:             notice.Sort,
		}
		gatewayWriteDataToWsResp, err := l.svcCtx.GatewayService.GatewayWriteDataToWsWrapper(context.Background(), &pb.GatewayWriteDataToWsWrapperReq{
			Filter: &pb.GatewayGetConnectionFilter{
				UserIds: []string{notice.UserId},
			},
			Data: &pb.GatewayWriteDataContent{
				DataType: pb.GatewayWriteDataType_PushNotice,
				Response: nil,
				Message:  nil,
				Notice:   pbNotice,
			},
		})
		if err != nil {
			l.Errorf("push subscription notice error: %v", err)
		}
		_ = gatewayWriteDataToWsResp
	}
}
