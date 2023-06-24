package conversationservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"go.mongodb.org/mongo-driver/bson"
	"sort"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListJoinedConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListJoinedConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJoinedConversationsLogic {
	return &ListJoinedConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListJoinedConversations 列出加入的会话
func (l *ListJoinedConversationsLogic) ListJoinedConversations(in *pb.ListJoinedConversationsReq) (*pb.ListJoinedConversationsResp, error) {
	{
		if in.Limit == 0 {
			in.Limit = 20
		}
	}

	filter := bson.M{
		"memberUserId":     in.Header.UserId,
		"conversationType": in.ConversationType,
	}
	// filter
	{
		if len(in.GetFilter().GetSettingList()) > 0 {
			for _, kv := range in.GetFilter().GetSettingList() {
				k := "settings." + kv.GetKey().String()
				switch kv.Operator {
				case pb.ListJoinedConversationsReq_Filter_SettingKV_Equal:
					if kv.OrExists {
						filter["$or"] = []bson.M{
							{
								k: kv.GetValue(),
							},
							{
								k: bson.M{
									"$exists": false,
								},
							},
						}
					} else {
						filter[k] = kv.GetValue()
					}
				case pb.ListJoinedConversationsReq_Filter_SettingKV_NotEqual:
					if kv.OrNotExists {
						filter["$or"] = []bson.M{
							{
								k: bson.M{
									"$ne": kv.GetValue(),
								},
							},
							{
								k: bson.M{
									"$exists": false,
								},
							},
						}
					} else {
						filter[k] = bson.M{
							"$ne": kv.GetValue(),
						}
					}
				}
			}
		}
	}
	conversationMembers := make([]*conversationmodel.ConversationMember, 0)
	queryI := l.svcCtx.ConversationMemberCollection.Find(l.ctx, filter)
	queryI = queryI.Limit(in.Limit)
	queryI = queryI.Skip(in.Cursor)
	queryI = queryI.Sort("joinTime")
	err := queryI.All(&conversationMembers)
	if err != nil {
		l.Errorf("find conversation member error: %v", err)
		return &pb.ListJoinedConversationsResp{}, err
	}
	resp := &pb.ListJoinedConversationsResp{
		ConversationList: make([]*pb.ListJoinedConversationsResp_Conversation, 0),
	}
	conversationMap := make(map[string]*pb.ListJoinedConversationsResp_Conversation)
	for _, conversationMember := range conversationMembers {
		conversationMap[conversationMember.ConversationId] = &pb.ListJoinedConversationsResp_Conversation{
			ConversationId: conversationMember.ConversationId,
			SelfMemberInfo: &pb.ListJoinedConversationsResp_Conversation_SelfMemberInfo{
				JoinTime: int64(conversationMember.JoinTime),
			},
		}
	}
	for _, conversation := range conversationMap {
		resp.ConversationList = append(resp.ConversationList, conversation)
	}
	// sort by join time asc
	sort.Slice(resp.ConversationList, func(i, j int) bool {
		return resp.ConversationList[i].SelfMemberInfo.JoinTime < resp.ConversationList[j].SelfMemberInfo.JoinTime
	})
	return resp, nil
}
