package groupservicelogic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

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
func (l *GroupCreateLogic) GroupCreate(in *peerpb.GroupCreateReq) (*peerpb.GroupCreateResp, error) {
	now := time.Now()
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
			return &peerpb.GroupCreateResp{}, err
		}
	}
	// 验证角色
	{
		user := userMap[in.Header.UserId]
		accountRole := user.GetAccountMap().Get(peerpb.AccountType_Role.String())
		if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Group.AllowRoleCreate) {
			return &peerpb.GroupCreateResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.GroupCreateNotAllowRole),
			}, nil
		}
	}
	//验证加群数量
	{
		{
			user := userMap[in.Header.UserId]
			if user.GetCountMap().JoinGroupCount >= int64(l.svcCtx.Config.Group.JoinedMaxCount) {
				return &peerpb.GroupCreateResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.GroupCreateJoinedMaxCount),
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
				setting.V = "true"
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
				return &peerpb.GroupCreateResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.GroupCreateRequiredName),
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
				return &peerpb.GroupCreateResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.GroupCreateRequiredAvatar),
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
	conversationMembers := make([]*groupmodel.GroupMember, 0)
	for _, memberId := range append(in.MemberList, in.Header.UserId) {
		conversationMembers = append(conversationMembers, &groupmodel.GroupMember{
			GroupId:      group.GroupIdString(),
			MemberUserId: memberId,
			JoinTime:     primitive.NewDateTimeFromTime(now),
			JoinSource: bson.M{
				"event": "groupCreate",
				"from":  in.Header.UserId,
			},
			Settings: bson.M{},
		})
	}
	if len(append(in.MemberList, in.Header.UserId)) > 0 {
		_, err := l.svcCtx.GroupMemberCollection.InsertMany(context.Background(), conversationMembers)
		if err != nil {
			l.Errorf("insert conversation member error: %v", err)
			return &peerpb.GroupCreateResp{}, err
		}
	}
	_, err := l.svcCtx.GroupCollection.InsertOne(context.Background(), group)
	if err != nil {
		l.Errorf("insert group error: %v", err)
		return &peerpb.GroupCreateResp{}, err
	}
	//更新用户加群数量
	go func() {
		for _, userId := range append(in.MemberList, in.Header.UserId) {
			_, err := l.svcCtx.UserService.UpdateUserCountMap(context.Background(), &peerpb.UpdateUserCountMapReq{
				Header:     &peerpb.RequestHeader{UserId: userId},
				CountType:  peerpb.UpdateUserCountMapReq_joinGroupCount,
				Algorithm:  peerpb.UpdateUserCountMapReq_add,
				Count:      1,
				Statistics: false,
			})
			if err != nil {
				l.Errorf("update user count map error: %v", err)
			}
		}
		_, err = l.svcCtx.UserService.UpdateUserCountMap(context.Background(), &peerpb.UpdateUserCountMapReq{
			Header:     &peerpb.RequestHeader{UserId: in.Header.UserId},
			CountType:  peerpb.UpdateUserCountMapReq_createGroupCount,
			Algorithm:  peerpb.UpdateUserCountMapReq_add,
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
		_, err := l.svcCtx.MessageService.MessageSend(context.Background(), &peerpb.MessageSendReq{
			Header: &peerpb.RequestHeader{UserId: in.Header.UserId},
			Message: &peerpb.Message{
				ConversationId:   group.GroupIdString(),
				ConversationType: peerpb.ConversationType_Group,
				Sender: &peerpb.Message_Sender{
					Id:     owner.UserId,
					Name:   owner.Nickname,
					Avatar: owner.Avatar,
					Extra:  "",
				},
				Content: utils.Proto.Marshal(&peerpb.MessageContentText{
					Items: []*peerpb.MessageContentText_Item{{
						Type:  peerpb.MessageContentText_Item_TEXT,
						Text:  l.svcCtx.Config.Group.DefaultWelcomeMessage,
						Image: nil,
						At:    nil,
					}},
				}),
				ContentType: peerpb.MessageContentType_Text,
				SendTime:    uint32(time.Now().UnixMilli()),
				Option: &peerpb.Message_Option{
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
		notices := make([]*peerpb.Message, len(in.MemberList))
		for _, userId := range in.MemberList {
			notice := &peerpb.Message{
				MessageId:        utils.Snowflake.String(),
				ConversationId:   peerpb.GetSingleChatConversationId(channelmodel.ConversationIdGroupHelper, userId),
				ConversationType: peerpb.ConversationType_Single,
				Content: utils.Proto.Marshal(&peerpb.NoticeContentJoinNewGroup{
					GroupId: group.GroupIdString(),
				}),
				ContentType: peerpb.MessageContentType_JoinNewGroup,
				Option: &peerpb.Message_Option{
					StorageForServer: true,
					StorageForClient: true,
					CountUnread:      false,
				},
				Sender: &peerpb.Message_Sender{
					Id:         channelmodel.ConversationIdGroupHelper,
					SenderType: peerpb.SenderType_ChannelSender,
					Name:       channelmodel.GroupHelperChannel.Nickname,
					Avatar:     channelmodel.GroupHelperChannel.Avatar,
					Extra:      "",
				},
			}
			notices = append(notices, notice)
		}
		utils.Retry.Do(func() error {
			_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &peerpb.NoticeSendReq{
				Header:  in.Header,
				Notices: notices,
			})
			if err != nil {
				l.Errorf("send notice error: %v", err)
			}
			return err
		})
	}()
	return &peerpb.GroupCreateResp{GroupId: utils.AnyString(group.GroupId)}, nil
}

func (l *GroupCreateLogic) toUserInfo_(in *peerpb.GroupCreateReq, userId string, toUserInfo *usermodel.User) error {
	toUserModelResp, err := l.svcCtx.UserService.GetUserModelById(context.Background(), &peerpb.GetUserModelByIdReq{
		Header: in.Header,
		UserId: userId,
		Opt: &peerpb.GetUserModelByIdReq_Opt{
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

func (l *GroupCreateLogic) UserMap_(in *peerpb.GroupCreateReq, userMap *map[string]*usermodel.User, settingMap *map[string]map[string]*usermodel.UserSetting, userIds []string) error {
	toUserModelsResp, err := l.svcCtx.UserService.GetUserModelByIds(context.Background(), &peerpb.GetUserModelByIdsReq{
		Header:  in.Header,
		UserIds: userIds,
		Opt: &peerpb.GetUserModelByIdsReq_Opt{
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
