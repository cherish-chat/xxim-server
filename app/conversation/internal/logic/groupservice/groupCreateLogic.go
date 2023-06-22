package groupservicelogic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupCreate 创建群组
func (l *GroupCreateLogic) GroupCreate(in *pb.GroupCreateReq) (*pb.GroupCreateResp, error) {
	now := time.Now()
	//是否允许创建群组
	if !utils.EnumInSlice(in.Header.Platform, l.svcCtx.Config.Group.Create.AllowPlatform) {
		return &pb.GroupCreateResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "group_create_not_allow_platform")),
		}, nil
	}
	var (
		userMap    = make(map[string]*usermodel.User)
		settingMap = make(map[string]map[string]*usermodel.UserSetting)
	)
	// MemberIds中是否包含自己 如果包含 删除
	{
		tmp := make([]string, 0)
		for _, v := range in.MemberList {
			if v != in.Header.UserId {
				tmp = append(tmp, v)
			}
		}
		in.MemberList = tmp
	}
	// 获取用户信息
	{
		err := l.UserMap_(in, &userMap, &settingMap, append(in.MemberList, in.Header.UserId))
		if err != nil {
			l.Errorf("get user info error: %v", err)
			return &pb.GroupCreateResp{}, err
		}
	}
	// 验证角色
	{
		user := userMap[in.Header.UserId]
		accountRole := user.GetAccountMap().Get(pb.AccountTypeRole)
		if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Group.Create.AllowRole) {
			return &pb.GroupCreateResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "group_create_not_allow_role")),
			}, nil
		}
	}
	//验证加群数量
	{
		{
			user := userMap[in.Header.UserId]
			if user.GetCountMap().JoinGroupCount >= int64(l.svcCtx.Config.Group.JoinedMaxCount) {
				return &pb.GroupCreateResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "group_create_joined_max_count")),
				}, nil
			}
		}
		tmp := make([]string, 0)
		for _, userId := range in.MemberList {
			user := userMap[userId]
			if user.GetCountMap().JoinGroupCount >= int64(l.svcCtx.Config.Group.JoinedMaxCount) {
				continue
			}
			tmp = append(tmp, userId)
		}
		in.MemberList = tmp
	}
	//验证设置
	{
		tmp := make([]string, 0)
		for _, userId := range in.MemberList {
			setting := settingMap[userId][usermodel.UserSettingKeyAllowBeInvitedGroup]
			if setting.V == "" {
				setting.V = utils.AnyString(l.svcCtx.Config.Group.Invite.UserDefaultAllow)
			}
			if setting.V == "false" || setting.V == "0" {
				//不允许被邀请
				continue
			}
			tmp = append(tmp, userId)
		}
		in.MemberList = tmp
	}
	group := &groupmodel.Group{
		GroupId:        groupmodel.GroupModel.GenerateGroupId(),
		GroupName:      "",
		GroupAvatar:    "",
		GroupInfo:      utils.Map.SS2SA(in.InfoMap),
		OwnerUserId:    in.Header.UserId,
		ManagerUserIds: make([]string, 0),
		CreateTime:     primitive.NewDateTimeFromTime(now),
		UpdateTime:     primitive.NewDateTimeFromTime(now),
		DismissTime:    0,
		MemberCount:    len(in.MemberList) + 1,
	}
	//group info
	{
		if in.Name == nil {
			if l.svcCtx.Config.Group.RequiredName {
				return &pb.GroupCreateResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "group_create_required_name")),
				}, nil
			} else {
				switch l.svcCtx.Config.Group.DefaultNameRule {
				case "byMember":
					name := ""
					for _, user := range userMap {
						name += user.Nickname + "、"
					}
					name = strings.TrimRight(name, "、")
					group.GroupName = name
					//长度不得大于32
					group.GroupName = utils.String.Utf8Split(group.GroupName, 32)
				case "fixed":
					name := l.svcCtx.Config.Group.FixedName
					group.GroupName = name
				}
			}
		}
		if in.Avatar == nil {
			if l.svcCtx.Config.Group.RequiredAvatar {
				return &pb.GroupCreateResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "group_create_required_avatar")),
				}, nil
			} else {
				switch l.svcCtx.Config.Group.DefaultAvatarRule {
				case "byName":
					avatar := fmt.Sprintf("/image/generateAvatar?text=%s&w=200&h=200&bg=%s&fg=%s", url.QueryEscape(group.GroupName), utils.Random.SliceString(l.svcCtx.Config.Group.ByNameAvatarBgColors), utils.Random.SliceString(l.svcCtx.Config.Group.ByNameAvatarFgColors))
					group.GroupAvatar = avatar
				case "fixed":
					avatar := l.svcCtx.Config.Group.FixedAvatar
					group.GroupAvatar = avatar
				}
			}
		}
	}
	conversationMembers := make([]*conversationmodel.ConversationMember, 0)
	for _, memberId := range in.MemberList {
		conversationMembers = append(conversationMembers, &conversationmodel.ConversationMember{
			ConversationId:   group.GroupId,
			ConversationType: pb.ConversationType_Group,
			MemberUserId:     memberId,
			JoinTime:         primitive.NewDateTimeFromTime(now),
			JoinSource: bson.M{
				"event": "groupCreate",
				"from":  in.Header.UserId,
			},
			Settings: make([]*conversationmodel.ConversationSetting, 0),
		})
	}
	if len(in.MemberList) > 0 {
		_, err := l.svcCtx.ConversationMemberCollection.InsertMany(l.ctx, conversationMembers)
		if err != nil {
			l.Errorf("insert conversation member error: %v", err)
			return &pb.GroupCreateResp{}, err
		}
	}
	_, err := l.svcCtx.GroupCollection.InsertOne(l.ctx, group)
	if err != nil {
		l.Errorf("insert group error: %v", err)
		return &pb.GroupCreateResp{}, err
	}
	//更新用户加群数量
	go func() {
		for _, userId := range append(in.MemberList, in.Header.UserId) {
			_, err := l.svcCtx.InfoService.UpdateUserCountMap(context.Background(), &pb.UpdateUserCountMapReq{
				Header:     &pb.RequestHeader{UserId: userId},
				CountType:  pb.UpdateUserCountMapReq_joinGroupCount,
				Algorithm:  pb.UpdateUserCountMapReq_add,
				Count:      1,
				Statistics: false,
			})
			if err != nil {
				l.Errorf("update user count map error: %v", err)
			}
		}
		_, err = l.svcCtx.InfoService.UpdateUserCountMap(context.Background(), &pb.UpdateUserCountMapReq{
			Header:     &pb.RequestHeader{UserId: in.Header.UserId},
			CountType:  pb.UpdateUserCountMapReq_createGroupCount,
			Algorithm:  pb.UpdateUserCountMapReq_add,
			Count:      1,
			Statistics: false,
		})
		if err != nil {
			l.Errorf("update user count map error: %v", err)
		}
	}()
	//发消息
	go func() {
		owner := userMap[in.Header.UserId]
		_, err := l.svcCtx.MessageService.MessageSend(context.Background(), &pb.MessageSendReq{
			Header: &pb.RequestHeader{UserId: in.Header.UserId},
			Message: &pb.Message{
				ConversationId:   group.GroupIdString(),
				ConversationType: pb.ConversationType_Group,
				Sender: &pb.Message_Sender{
					Id:     owner.UserId,
					Name:   owner.Nickname,
					Avatar: owner.Avatar,
					Extra:  "",
				},
				Content: utils.Json.MarshalToBytes(&pb.MessageContentText{
					Items: []*pb.MessageContentText_Item{{
						Type:  pb.MessageContentText_Item_TEXT,
						Text:  l.svcCtx.Config.Group.Invite.DefaultWelcomeMessage,
						Image: nil,
						At:    nil,
					}},
				}),
				ContentType: pb.MessageContentType_Text,
				SendTime:    time.Now().UnixMilli(),
				Option: &pb.Message_Option{
					StorageForServer: true,
					StorageForClient: true,
					NeedDecrypt:      false,
					CountUnread:      true,
				},
				ExtraMap: map[string]string{
					"platformSource": "server",
				},
			},
			DisableQueue: false,
		})
		if err != nil {
			l.Errorf("send message error: %v", err)
		}
	}()
	//发通知
	go func() {
		notice := &noticemodel.BroadcastNotice{
			NoticeId:         utils.Snowflake.String(),
			ConversationId:   group.GroupIdString(),
			ConversationType: pb.ConversationType_Group,
			Content:          utils.Json.MarshalToString(&pb.NoticeContentNewGroup{}),
			ContentType:      pb.NoticeContentType_NewGroup,
			UpdateTime:       primitive.NewDateTimeFromTime(time.Now()),
		}
		utils.Retry.Do(func() error {
			_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &pb.NoticeSendReq{
				Header:    in.Header,
				Notice:    notice.ToPb(),
				Broadcast: true,
			})
			if err != nil {
				l.Errorf("send notice error: %v", err)
			}
			return err
		})
	}()
	return &pb.GroupCreateResp{GroupId: utils.AnyString(group.GroupId)}, nil
}

func (l *GroupCreateLogic) toUserInfo_(in *pb.GroupCreateReq, userId string, toUserInfo *usermodel.User) error {
	toUserModelResp, err := l.svcCtx.InfoService.GetUserModelById(l.ctx, &pb.GetUserModelByIdReq{
		Header: in.Header,
		UserId: userId,
		Opt: &pb.GetUserModelByIdReq_Opt{
			WithUserSettings: false,
			UserSettingKeys:  []string{},
		},
	})
	if err != nil {
		l.Errorf("get user model error: %v", err)
		return err
	}

	user := &usermodel.User{}
	err = utils.Json.Unmarshal(toUserModelResp.UserModelJson, user)
	if err != nil {
		l.Errorf("unmarshal user model error: %v", err)
		return err
	}
	*toUserInfo = *user

	return nil
}

func (l *GroupCreateLogic) UserMap_(in *pb.GroupCreateReq, userMap *map[string]*usermodel.User, settingMap *map[string]map[string]*usermodel.UserSetting, userIds []string) error {
	toUserModelsResp, err := l.svcCtx.InfoService.GetUserModelByIds(l.ctx, &pb.GetUserModelByIdsReq{
		Header:  in.Header,
		UserIds: userIds,
		Opt: &pb.GetUserModelByIdsReq_Opt{
			WithUserSettings: true,
			UserSettingKeys: []string{
				//是否允许被邀请入群
				usermodel.UserSettingKeyAllowBeInvitedGroup,
			},
		},
	})
	if err != nil {
		l.Errorf("get user model error: %v", err)
		return err
	}
	var um = make(map[string]*usermodel.User)
	var sm = make(map[string]map[string]*usermodel.UserSetting)
	for userId, v := range toUserModelsResp.UserModelJsons {
		user := &usermodel.User{}
		err = utils.Json.Unmarshal(v, user)
		if err != nil {
			l.Errorf("unmarshal user model error: %v", err)
			return err
		}
		um[userId] = user
	}
	for userId, v := range toUserModelsResp.UserSettingsJsons {
		userSettingMap := make(map[string]*usermodel.UserSetting)
		err = utils.Json.Unmarshal(v, &userSettingMap)
		if err != nil {
			l.Errorf("unmarshal user setting error: %v", err)
			return err
		}
		sm[userId] = userSettingMap
	}
	//是否有用户不存在
	for _, v := range userIds {
		if _, ok := um[v]; !ok {
			return status.Errorf(codes.Internal, "user not exist: %s", v)
		}
	}
	//是否有setting不存在
	for _, v := range userIds {
		if sm, ok := sm[v]; !ok {
			return status.Errorf(codes.Internal, "user setting not exist: %s", v)
		} else {
			if _, ok := sm[usermodel.UserSettingKeyAllowBeInvitedGroup]; !ok {
				return status.Errorf(codes.Internal, "user setting not exist: %s", v)
			}
		}
	}
	*userMap = um
	*settingMap = sm
	return nil
}
