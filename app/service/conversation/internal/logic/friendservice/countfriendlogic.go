package friendservicelogic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

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

func (l *CountFriendLogic) CountFriend(in *peerpb.CountFriendReq) (*peerpb.CountFriendResp, error) {
	count, err := l.svcCtx.FriendCollection.Find(context.Background(), bson.M{
		"userId": in.Header.UserId,
	}).Count()
	if err != nil {
		l.Errorf("find friend error: %v", err)
		return &peerpb.CountFriendResp{}, err
	}
	return &peerpb.CountFriendResp{Count: uint32(count)}, nil
}
