package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
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
	return l.InviteFriendToGroupWithoutVerify(in)
}

func (l *InviteFriendToGroupLogic) InviteFriendToGroupWithoutVerify(in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	members := make([]*groupmodel.GroupMember, 0)
	for _, member := range in.FriendIds {
		members = append(members, &groupmodel.GroupMember{
			GroupId:    in.GroupId,
			UserId:     member,
			CreateTime: l.now.UnixMilli(),
		})
	}
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := xorm.InsertMany(tx, &groupmodel.GroupMember{}, members)
		if err != nil {
			l.Errorf("InviteFriendToGroup InsertMany error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "邀请失败"), l.svcCtx.T(in.CommonReq.Language, "群成员可能已经存在"))}, nil
	}
	if in.MinSeq != nil {
		// 设置群成员的最小消息序列号
		_, err = l.svcCtx.MsgService().BatchSetMinSeq(l.ctx, &pb.BatchSetMinSeqReq{
			CommonReq:  in.CommonReq,
			ConvId:     in.GroupId,
			UserIdList: in.FriendIds,
			MinSeq:     *in.MinSeq,
		})
		if err != nil {
			l.Errorf("InviteFriendToGroup BatchSetMinSeq error: %v", err)
			return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	// 删除缓存
	{
		err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), append(in.FriendIds, in.CommonReq.UserId)...)
		if err != nil {
			l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
			return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.InviteFriendToGroupResp{}, nil
}
