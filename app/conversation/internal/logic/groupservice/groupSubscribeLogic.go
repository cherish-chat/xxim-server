package groupservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GroupSubscribeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSubscribeLogic {
	return &GroupSubscribeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupSubscribe 群组订阅
func (l *GroupSubscribeLogic) GroupSubscribe(in *pb.GroupSubscribeReq) (*pb.GroupSubscribeResp, error) {
	// 1. 获取用户加入的群组
	listJoinedGroupsResp, err := l.svcCtx.ConversationService.ListJoinedConversations(l.ctx, &pb.ListJoinedConversationsReq{
		Header:           in.Header,
		ConversationType: pb.ConversationType_Group,
		Cursor:           0,
		Limit:            int64(l.svcCtx.Config.Group.JoinedMaxCount),
		Filter: &pb.ListJoinedConversationsReq_Filter{
			SettingList: []*pb.ListJoinedConversationsReq_Filter_SettingKV{{
				//是否被屏蔽这个key != "true" or 不存在
				Key:         pb.ConversationSettingKey_IsBlocked, // 是否被屏蔽
				Value:       "true",
				Operator:    pb.ListJoinedConversationsReq_Filter_SettingKV_NotEqual,
				OrNotExists: true,
				OrExists:    false,
			}},
		},
		Option: &pb.ListJoinedConversationsReq_Option{
			IncludeSelfMemberInfo: true,
		},
	})
	if err != nil {
		l.Errorf("listJoinedGroupsResp err: %v", err)
		return &pb.GroupSubscribeResp{}, err
	}

	if len(listJoinedGroupsResp.GetConversationList()) == 0 {
		return &pb.GroupSubscribeResp{}, nil
	}

	var groupIds []string
	for _, group := range listJoinedGroupsResp.GetConversationList() {
		groupIds = append(groupIds, group.GetConversationId())
	}

	// 2. 批量更新用户的群组订阅时间
	// filter: {memberUserId: "xxx", groupId: {$in: ["xxx", "xxx"]}}
	// set: {subscribeTime: "xxx"}
	bulk := l.svcCtx.GroupSubscribeCollection.Bulk()
	for _, groupId := range groupIds {
		bulk.Upsert(bson.M{
			"memberUserId": in.Header.UserId,
			"groupId":      groupId,
		}, &groupmodel.GroupSubscribe{
			GroupId:       groupId,
			MemberUserId:  in.Header.UserId,
			SubscribeTime: primitive.NewDateTimeFromTime(time.Now()),
		})
	}
	_, err = bulk.Run(l.ctx)
	if err != nil {
		l.Errorf("bulk.Run err: %v", err)
		return &pb.GroupSubscribeResp{}, err
	}
	return &pb.GroupSubscribeResp{}, nil
}
