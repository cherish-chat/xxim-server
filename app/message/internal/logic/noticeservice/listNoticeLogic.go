package noticeservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeLogic {
	return &ListNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListNotice 获取通知列表
func (l *ListNoticeLogic) ListNotice(in *pb.ListNoticeReq) (*pb.ListNoticeResp, error) {
	var broadcastResults []*noticemodel.BroadcastNotice
	err := l.svcCtx.BroadcastNoticeCollection.Find(l.ctx, bson.M{
		"conversationId":   in.ConversationId,
		"conversationType": in.ConversationType,
		"updateTime": bson.M{
			"$gt": primitive.DateTime(in.UpdateTimeGt),
		},
	}).Sort("updateTime").Limit(in.Limit).All(&broadcastResults)
	if err != nil {
		l.Errorf("find broadcast notice error: %v", err)
		return nil, err
	}
	l.Debugf("broadcastResults: %v", broadcastResults)
	var subscriptionResults []*noticemodel.SubscriptionNotice
	var contentMap = make(map[string]*noticemodel.SubscriptionNoticeContent)
	var contentIds []string
	if in.ConversationType == pb.ConversationType_Subscription {
		err = l.svcCtx.SubscriptionNoticeCollection.Find(l.ctx, bson.M{
			"userId":         in.Header.UserId,
			"subscriptionId": in.ConversationId,
			"updateTime": bson.M{
				"$gt": primitive.DateTime(in.UpdateTimeGt),
			},
		}).Sort("updateTime").Limit(in.Limit).All(&subscriptionResults)
		if err != nil {
			l.Errorf("find subscription notice error: %v", err)
			return nil, err
		}
		for _, result := range subscriptionResults {
			contentIds = append(contentIds, result.ContentId)
		}
	}
	l.Debugf("subscriptionResults: %v", subscriptionResults)
	if len(contentIds) > 0 {
		var contents []*noticemodel.SubscriptionNoticeContent
		err = l.svcCtx.SubscriptionNoticeContentCollection.Find(l.ctx, bson.M{
			"contentId": bson.M{
				"$in": contentIds,
			},
		}).All(&contents)
		if err != nil {
			l.Errorf("find subscription notice content error: %v", err)
			return nil, err
		}
		for _, content := range contents {
			contentMap[content.ContentId] = content
		}
	}

	var resp = &pb.ListNoticeResp{
		Notices: make([]*pb.Notice, 0),
	}
	for _, result := range broadcastResults {
		resp.Notices = append(resp.Notices, result.ToPb())
	}
	for _, result := range subscriptionResults {
		content, _ := contentMap[result.ContentId]
		if content == nil {
			content = &noticemodel.SubscriptionNoticeContent{}
		}
		resp.Notices = append(resp.Notices, result.ToPb(content.Content))
	}
	return resp, nil
}
