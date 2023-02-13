package logic

import (
	"context"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteFriendToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	now time.Time
}

func NewInviteFriendToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteFriendToGroupLogic {
	return &InviteFriendToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		now:    time.Now(),
	}
}

// InviteFriendToGroup 邀请好友加入群聊
func (l *InviteFriendToGroupLogic) InviteFriendToGroup(in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	if len(in.FriendIds) == 0 {
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "请选择好友"))}, nil
	}
	// 验证是否是我的好友
	areFriendsResp, err := l.svcCtx.RelationService().AreFriends(l.ctx, &pb.AreFriendsReq{
		CommonReq: in.CommonReq,
		A:         in.CommonReq.UserId,
		BList:     in.FriendIds,
	})
	if err != nil {
		l.Errorf("InviteFriendToGroup AreFriends error: %v", err)
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, id := range in.FriendIds {
		if is, ok := areFriendsResp.FriendList[id]; !ok || !is {
			l.Errorf("InviteFriendToGroup AreFriends error: %v", err)
			return &pb.InviteFriendToGroupResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "只能邀请好友加入群聊"))}, err
		}
	}
	return l.inviteFriendToGroup(in)
}

func (l *InviteFriendToGroupLogic) inviteFriendToGroup(in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	//TODO implement the business logic of InviteFriendToGroup
	return &pb.InviteFriendToGroupResp{}, nil
}
