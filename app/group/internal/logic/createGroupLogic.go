package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"strconv"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	now time.Time
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		now:    time.Now(),
	}
}

// CreateGroup 创建群聊
func (l *CreateGroupLogic) CreateGroup(in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	// 获取群id
	groupIdInt, err := l.svcCtx.Redis().HincrbyCtx(l.ctx, rediskey.IncrId(), rediskey.IncrGroup(), 1)
	if err != nil {
		l.Errorf("CreateGroup HincrbyCtx error: %v", err)
		return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	group := &groupmodel.Group{
		Id:          strconv.Itoa(groupIdInt),
		Name:        l.svcCtx.SystemConfigMgr.GetOrDefaultCtx(l.ctx, "default_group_name", "未命名群聊"),
		Avatar:      utils.AnyRandomInSlice(l.svcCtx.SystemConfigMgr.GetSliceCtx(l.ctx, "default_group_avatars"), ""),
		Owner:       in.CommonReq.UserId,
		Managers:    make([]string, 0),
		CreateTime:  time.Now().UnixMilli(),
		DismissTime: 0,
		Description: l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "default_group_description"),
		Setting: groupmodel.GroupSetting{
			AllMute:                  false,
			SpeakLimit:               nil,
			MaxMember:                int32(utils.AnyToInt64(l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "default_group_max_member"))),
			MemberCanStartTempChat:   true,
			MemberCanInviteFriend:    true,
			NewMemberHistoryMsgCount: int32(utils.AnyToInt64(l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "default_group_new_member_history_msg_count"))),
			AnonymousChat:            true,
			JoinGroupOption: groupmodel.JoinGroupOption{
				Type:     pb.GroupSetting_JoinGroupOpt_NEED_VERIFY,
				Question: utils.AnyToString(l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "default_group_join_group_question")),
				Answer:   "",
			},
		},
	}
	if in.Name != nil {
		group.Name = *in.Name
	}
	if in.Avatar != nil {
		group.Avatar = *in.Avatar
	}
	// 是否携带群成员
	// 插入到群成员表
	var inviteFriendToGroupResp *pb.InviteFriendToGroupResp
	xtrace.StartFuncSpan(l.ctx, "InviteFriendToGroupWithoutVerify", func(ctx context.Context) {
		inviteFriendToGroupResp, err = NewInviteFriendToGroupLogic(ctx, l.svcCtx).InviteFriendToGroupWithoutVerify(&pb.InviteFriendToGroupReq{
			CommonReq: in.CommonReq,
			GroupId:   group.Id,
			FriendIds: append(in.Members, in.CommonReq.UserId),
		})
	})
	if err != nil {
		l.Errorf("CreateGroup InviteFriendToGroup error: %v", err)
		return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if inviteFriendToGroupResp.CommonResp != nil && inviteFriendToGroupResp.CommonResp.Failed() {
		l.Errorf("CreateGroup InviteFriendToGroup failed: %v", inviteFriendToGroupResp.CommonResp)
		return &pb.CreateGroupResp{CommonResp: inviteFriendToGroupResp.CommonResp}, nil
	}
	// 插入群表
	err = xorm.InsertOne(l.svcCtx.Mysql(), group)
	if err != nil {
		// retry
		l.Errorf("CreateGroup InsertOne error: %v", err)
		return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.CreateGroupResp{
		GroupId: utils.AnyPtr(group.Id),
	}, nil
}
