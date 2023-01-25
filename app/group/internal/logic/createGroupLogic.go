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
	"go.opentelemetry.io/otel/propagation"
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
	in.Members = utils.Set(in.Members)
	// 判断 members 是否包含自己
	if utils.InSlice(in.Members, in.CommonReq.UserId) {
		// 报错
		return &pb.CreateGroupResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "不能邀请自己"),
			),
		}, nil
	}
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
				Type:     0,
				Question: utils.AnyToString(l.svcCtx.SystemConfigMgr.GetCtx(l.ctx, "default_group_join_group_question")),
				Answer:   "",
			},
		},
		MemberCount: 1 + len(in.Members),
	}
	if in.Name != nil {
		group.Name = *in.Name
	}
	if in.Avatar != nil {
		group.Avatar = *in.Avatar
	}
	// 删除缓存
	{
		err = groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), group.Id)
		if err != nil {
			l.Errorf("CreateGroup CleanGroupCache error: %v", err)
			return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), append(in.Members, in.CommonReq.UserId)...)
		if err != nil {
			l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
			return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		// 群成员
		members := make([]*groupmodel.GroupMember, 0)
		for _, member := range in.Members {
			members = append(members, &groupmodel.GroupMember{
				GroupId:    group.Id,
				UserId:     member,
				CreateTime: l.now.UnixMilli(),
				Role:       groupmodel.RoleType_MEMBER,
			})
		}
		// 群主
		members = append(members, &groupmodel.GroupMember{
			GroupId:    group.Id,
			UserId:     group.Owner,
			CreateTime: l.now.UnixMilli(),
			Role:       groupmodel.RoleType_OWNER,
		})
		err := xorm.InsertMany(tx, &groupmodel.GroupMember{}, members)
		if err != nil {
			l.Errorf("InviteFriendToGroup InsertMany error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		err := xorm.InsertOne(l.svcCtx.Mysql(), group)
		if err != nil {
			l.Errorf("CreateGroup InsertOne error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvIdGroup(group.Id),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvMsg:    false,
			},
			ContentType: pb.NoticeContentType_CreateGroup,
			Content: utils.AnyToBytes(pb.NoticeContent_CreateGroup{
				GroupId: group.Id,
			}),
			Title: "",
			Ext:   nil,
		}
		err = notice.Insert(l.ctx, tx)
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.CreateGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	{
		// 预热缓存
		// 刷新订阅
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			// 删除缓存
			{
				err = groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), group.Id)
				if err != nil {
					l.Errorf("CreateGroup CleanGroupCache error: %v", err)
					return err
				}
				err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), append(in.Members, in.CommonReq.UserId)...)
				if err != nil {
					l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
					return err
				}
			}
			// 预热缓存
			go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarmUp", func(ctx context.Context) {
				_, err := groupmodel.ListGroupByIdsFromMysql(ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), []string{group.Id})
				if err != nil {
					l.Errorf("CreateGroup ListGroupByIdsFromMysql error: %v", err)
				}
				for _, userId := range append(in.Members, in.CommonReq.UserId) {
					_, err = groupmodel.ListGroupsByUserIdFromMysql(ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), userId)
					if err != nil {
						l.Errorf("CreateGroup ListGroupsByUserIdFromMysql error: %v", err)
					}
				}
			}, propagation.MapCarrier{
				"group_id": group.Id,
			})
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: append(in.Members, in.CommonReq.UserId)})
			if err != nil {
				l.Errorf("FlushUsersSubConv failed, err: %v", err)
				return err
			}
			_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				UserId:    "",
				ConvId:    pb.HiddenConvIdGroup(group.Id),
				DeviceId:  nil,
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
