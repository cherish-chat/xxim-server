package friendservicelogic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountFriendLogic {
	return &CountFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CountFriend 统计好友数量
func (l *CountFriendLogic) CountFriend(in *pb.CountFriendReq) (*pb.CountFriendResp, error) {
	count, err := l.svcCtx.FriendCollection.Find(l.ctx, bson.M{
		"userId": in.Header.UserId,
	}).Count()
	if err != nil {
		l.Errorf("find friend error: %v", err)
		return &pb.CountFriendResp{}, err
	}
	return &pb.CountFriendResp{Count: count}, nil
}
