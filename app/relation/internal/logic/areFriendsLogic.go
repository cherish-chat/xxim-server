package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreFriendsLogic {
	return &AreFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreFriendsLogic) AreFriends(in *pb.AreFriendsReq) (*pb.AreFriendsResp, error) {
	friend, err := relationmodel.AreMyFriend(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&relationmodel.Friend{}), in.A, in.BList)
	if err != nil {
		l.Errorf("AreMyFriend failed, err: %v", err)
		return &pb.AreFriendsResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	friend[in.Requester.Id] = true
	return &pb.AreFriendsResp{FriendList: friend}, nil
}
