package groupservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/groupmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAfterKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAfterKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAfterKeepAliveLogic {
	return &GroupAfterKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupAfterKeepAliveLogic) GroupAfterKeepAlive(in *peerpb.GroupAfterKeepAliveReq) (*peerpb.GroupAfterKeepAliveResp, error) {
	listMyJoinedGroupResp, err := NewListMyJoinedGroupLogic(context.Background(), l.svcCtx).ListMyJoinedGroup(&peerpb.ListMyJoinedGroupReq{
		Header: in.Header,
		Cursor: 0,
		Limit:  2000,
		Filter: nil,
		Option: nil,
	})
	if err != nil {
		l.Errorf("list my joined group error: %v", err)
		return &peerpb.GroupAfterKeepAliveResp{}, err
	}
	if len(listMyJoinedGroupResp.GroupList) == 0 {
		return &peerpb.GroupAfterKeepAliveResp{}, nil
	}

	var groupIds = make([]string, 0)
	for _, group := range listMyJoinedGroupResp.GroupList {
		groupIds = append(groupIds, group.GroupId)
	}

	// 2. 批量更新用户的群组订阅时间
	// filter: {memberUserId: "xxx", groupId: {$in: ["xxx", "xxx"]}}
	// set: {subscribeTime: "xxx"}
	bulk := l.svcCtx.GroupSubscribeCacheCollection.Bulk()
	for _, groupId := range groupIds {
		bulk.Upsert(bson.M{
			"memberUserId": in.Header.UserId,
			"groupId":      groupId,
		}, &groupmodel.GroupSubscribeCache{
			GroupId:       groupId,
			MemberUserId:  in.Header.UserId,
			SubscribeTime: primitive.NewDateTimeFromTime(time.Now()),
		})
	}
	_, err = bulk.Run(context.Background())
	if err != nil {
		l.Errorf("bulk.Run err: %v", err)
		return &peerpb.GroupAfterKeepAliveResp{}, err
	}
	return &peerpb.GroupAfterKeepAliveResp{}, nil
}
