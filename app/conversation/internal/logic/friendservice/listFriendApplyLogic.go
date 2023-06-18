package friendservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendApplyLogic {
	return &ListFriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListFriendApply 列出好友申请
func (l *ListFriendApplyLogic) ListFriendApply(in *pb.ListFriendApplyReq) (*pb.ListFriendApplyResp, error) {
	filter := bson.M{}
	if in.Cursor > 0 {
		// < cursor
		filter["applyTime"] = bson.M{"$lt": primitive.DateTime(in.Cursor)}
	}
	if in.GetFilter().Status != nil {
		status := *in.GetFilter().Status
		filter["status"] = status
	}
	if in.GetOption().GetIncludeApplyByMe() {
		// fromId = in.Header.UserId or toId = in.Header.UserId
		filter["$or"] = bson.A{
			bson.M{"fromId": in.Header.UserId},
			bson.M{"toId": in.Header.UserId},
		}
	} else {
		// toId = in.Header.UserId
		filter["toId"] = in.Header.UserId
	}
	var result []*friendmodel.FriendApplyRecord
	err := l.svcCtx.FriendApplyRecordCollection.Find(l.ctx, filter).Sort("-applyTime").Limit(in.Limit).All(&result)
	if err != nil {
		l.Errorf("find friend apply record error: %v", err)
		return &pb.ListFriendApplyResp{}, err
	}
	if len(result) == 0 {
		return &pb.ListFriendApplyResp{
			Cursor:          0,
			FriendApplyList: nil,
		}, nil
	}
	var resp = &pb.ListFriendApplyResp{
		Cursor:          0,
		FriendApplyList: make([]*pb.ListFriendApplyResp_FriendApply, 0),
	}
	minCursor := int64(math.MaxInt64)
	for _, record := range result {
		resp.FriendApplyList = append(resp.FriendApplyList, &pb.ListFriendApplyResp_FriendApply{
			ApplyId:    record.ApplyId,
			FromUserId: record.FromId,
			ToUserId:   record.ToId,
			Message:    record.Message,
			Answer:     record.Answer,
		})
		if int64(record.ApplyTime) < minCursor {
			minCursor = int64(record.ApplyTime)
		}
	}
	return resp, nil
}
