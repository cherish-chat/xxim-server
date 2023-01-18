package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"
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
			MaxMember:                int32(utils.AnyToInt64(l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "app.default_group_max_member"))),
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
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := xorm.InsertOne(l.svcCtx.Mysql(), group)
		if err != nil {
			l.Errorf("CreateGroup InsertOne error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		// 发送一条订阅号消息 订阅号的convId = notice:group@groupId  noticeId = UpdateGroupInfo
		data := &pb.NoticeData{
			ConvId:         noticemodel.ConvIdGroup(group.Id),
			UnreadCount:    0,
			UnreadAbsolute: false,
			NoticeId:       "UpdateGroupInfo",
			ContentType:    0,
			Content:        []byte{},
			Options: &pb.NoticeData_Options{
				StorageForClient: false,
				UpdateConvMsg:    false,
				OnlinePushOnce:   false,
			},
			Ext: nil,
		}
		m := noticemodel.NoticeFromPB(data, true, "")
		err := m.Upsert(tx)
		if err != nil {
			l.Errorf("Upsert failed, err: %v", err)
		}
		return err
	})
	if err != nil {
		return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	{
		// 删除缓存
		// 预热缓存
		// 刷新订阅
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: append(in.Members, in.CommonReq.UserId)})
			if err != nil {
				l.Errorf("FlushUsersSubConv failed, err: %v", err)
				return err
			}
			_, err = l.svcCtx.NoticeService().SetUserSubscriptions(l.ctx, &pb.SetUserSubscriptionsReq{
				UserIds: append(in.Members, in.CommonReq.UserId),
			})
			if err != nil {
				l.Errorf("SetUserSubscriptions failed, err: %v", err)
				return err
			}
			_, err = l.svcCtx.NoticeService().SendNoticeData(l.ctx, &pb.SendNoticeDataReq{
				CommonReq: in.CommonReq,
				NoticeData: &pb.NoticeData{
					NoticeId: "UpdateGroupInfo",
					ConvId:   noticemodel.ConvIdGroup(group.Id),
				},
				UserId:      nil,
				IsBroadcast: utils.AnyPtr(true),
				Inserted:    utils.AnyPtr(true),
			})
			if err != nil {
				l.Errorf("SendNoticeData failed, err: %v", err)
			}
			return err
		})
		// 群主发送消息：欢迎加入群聊
		l.sendMsg(in, group)
	}

	return &pb.CreateGroupResp{
		GroupId: utils.AnyPtr(group.Id),
	}, nil
}

func (l *CreateGroupLogic) sendMsg(in *pb.CreateGroupReq, group *groupmodel.Group) {
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendMsg", func(ctx context.Context) {
		// 获取接受者info
		userByIds, err := l.svcCtx.UserService().MapUserByIds(ctx, &pb.MapUserByIdsReq{Ids: []string{in.CommonReq.UserId}})
		if err != nil {
			l.Errorf("MapUserByIds failed, err: %v", err)
		} else {
			selfInfo, ok := userByIds.Users[in.CommonReq.UserId]
			if ok {
				self := usermodel.UserFromBytes(selfInfo)
				_, err = msgservice.SendMsgSync(l.svcCtx.MsgService(), ctx, []*pb.MsgData{
					msgmodel.CreateTextMsgToGroup(
						&pb.UserBaseInfo{
							Id:       self.Id,
							Nickname: self.Nickname,
							Avatar:   self.Avatar,
							Xb:       self.Xb,
							Birthday: self.Birthday,
						},
						group.Id,
						l.svcCtx.T(in.CommonReq.Language, "欢迎加入群聊"),
						msgmodel.MsgOptions{
							OfflinePush:       true,
							StorageForServer:  true,
							StorageForClient:  true,
							UpdateUnreadCount: false,
							NeedDecrypt:       false,
							UpdateConvMsg:     true,
						},
						&msgmodel.MsgOfflinePush{
							Title:   group.Name,
							Content: "欢迎加入群聊",
							Payload: "",
						},
						nil,
					).ToMsgData(),
				})
				if err != nil {
					l.Errorf("SendMsgSync failed, err: %v", err)
					err = nil
				}
			}
		}
	}, nil)
}
