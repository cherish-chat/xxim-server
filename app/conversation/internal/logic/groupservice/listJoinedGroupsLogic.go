package groupservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"go.mongodb.org/mongo-driver/bson"
	"sort"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListJoinedGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListJoinedGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJoinedGroupsLogic {
	return &ListJoinedGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListJoinedGroups 列出加入的群组
func (l *ListJoinedGroupsLogic) ListJoinedGroups(in *pb.ListJoinedGroupsReq) (*pb.ListJoinedGroupsResp, error) {
	{
		if in.Limit == 0 {
			in.Limit = 20
		}
	}

	filter := bson.M{
		"memberUserId": in.Header.UserId,
	}
	// filter
	{
		if len(in.GetFilter().GetSettingList()) > 0 {
			for _, kv := range in.GetFilter().GetSettingList() {
				k := "settings." + kv.GetKey().String()
				switch kv.Operator {
				case pb.ListJoinedGroupsReq_Filter_SettingKV_Equal:
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
				case pb.ListJoinedGroupsReq_Filter_SettingKV_NotEqual:
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
		return &pb.ListJoinedGroupsResp{}, err
	}
	resp := &pb.ListJoinedGroupsResp{
		GroupList: make([]*pb.ListJoinedGroupsResp_Group, 0),
	}
	groupMap := make(map[string]*pb.ListJoinedGroupsResp_Group)
	for _, conversationMember := range conversationMembers {
		groupMap[conversationMember.ConversationId] = &pb.ListJoinedGroupsResp_Group{
			GroupId: conversationMember.ConversationId,
			SelfMemberInfo: &pb.ListJoinedGroupsResp_Group_SelfMemberInfo{
				JoinTime: int64(conversationMember.JoinTime),
			},
		}
	}
	for _, group := range groupMap {
		resp.GroupList = append(resp.GroupList, group)
	}
	// sort by join time asc
	sort.Slice(resp.GroupList, func(i, j int) bool {
		return resp.GroupList[i].SelfMemberInfo.JoinTime < resp.GroupList[j].SelfMemberInfo.JoinTime
	})
	return resp, nil
}
