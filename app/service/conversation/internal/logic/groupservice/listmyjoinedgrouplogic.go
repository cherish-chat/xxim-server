package groupservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/groupmodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyJoinedGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyJoinedGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyJoinedGroupLogic {
	return &ListMyJoinedGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyJoinedGroupLogic) ListMyJoinedGroup(in *peerpb.ListMyJoinedGroupReq) (*peerpb.ListMyJoinedGroupResp, error) {
	{
		if in.Limit == 0 {
			in.Limit = 2000
		}
	}
	filter := bson.M{
		"memberUserId": in.Header.UserId,
	}
	conversationMembers := make([]*groupmodel.GroupMember, 0)
	queryI := l.svcCtx.GroupMemberCollection.Find(context.Background(), filter)
	queryI = queryI.Limit(int64(in.Limit))
	queryI = queryI.Skip(int64(in.Cursor))
	queryI = queryI.Sort("joinTime")
	err := queryI.All(&conversationMembers)
	if err != nil {
		l.Errorf("find conversation member error: %v", err)
		return &peerpb.ListMyJoinedGroupResp{}, err
	}
	var resp = &peerpb.ListMyJoinedGroupResp{}
	for _, conversationMember := range conversationMembers {
		resp.GroupList = append(resp.GroupList, &peerpb.ListMyJoinedGroupResp_Group{
			GroupId: conversationMember.GroupId,
		})
	}
	return resp, nil
}
