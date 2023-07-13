package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/mr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendApplyLogic {
	return &FriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendApply 添加好友
func (l *FriendApplyLogic) FriendApply(in *peerpb.FriendApplyReq) (*peerpb.FriendApplyResp, error) {
	//预先验证
	{
		//是不是相同id
		if in.Header.UserId == in.ToUserId {
			return &peerpb.FriendApplyResp{}, nil
		}
	}

	// 查询两个用户信息和用户设置
	var (
		fromUserId       = in.Header.UserId
		toUserId         = in.ToUserId
		fromUserInfo     = &usermodel.User{}
		toUserInfo       = &usermodel.User{}
		toUserSettingMap = make(map[string]*usermodel.UserSetting)
	)

	var getPreSourceFunctions []func() error
	getPreSourceFunctions = append(getPreSourceFunctions, func() error {
		return l.fromUserInfo_(in, fromUserId, fromUserInfo)
	}, func() error {
		return l.toUserInfo_(in, toUserId, toUserInfo, &toUserSettingMap)
	})
	err := mr.Finish(getPreSourceFunctions...)
	if err != nil {
		return &peerpb.FriendApplyResp{}, err
	}

	//验证用户信息
	friendApplyResp, err := l.verifyUserInfo_(in, fromUserInfo, toUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	//验证是否已经是好友
	friendApplyResp, err = l.verifyAreFriend_(in, fromUserInfo, toUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	//验证两人的好友数量上限
	friendApplyResp, err = l.verifyFriendLimit_(in, fromUserInfo, toUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	//验证from用户是否已经申请过to用户
	exist, err := l.verifyApplyExist_(in, fromUserInfo, toUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	if exist {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAlready),
		}, nil
	}
	//验证to用户是否已经申请过from用户
	skipApply, err := l.verifyApplyExist_(in, toUserInfo, fromUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	if !skipApply {
		//验证用户设置
		friendApplyResp, skipApply, err = l.verifyUserSetting_(in, fromUserInfo, toUserInfo, toUserSettingMap)
		if err != nil {
			return friendApplyResp, err
		}
	}
	//跳过申请
	if skipApply {
		err = NewFriendApplyHandleLogic(context.Background(), l.svcCtx).AddFriend(fromUserInfo.UserId, toUserInfo)
		return &peerpb.FriendApplyResp{}, err
	}

	//处理
	err = l.handle_(in, fromUserInfo, toUserInfo, toUserSettingMap)
	if err != nil {
		return &peerpb.FriendApplyResp{}, err
	}
	return &peerpb.FriendApplyResp{}, nil
}

func (l *FriendApplyLogic) fromUserInfo_(in *peerpb.FriendApplyReq, fromUserId string, fromUserInfo *usermodel.User) error {
	fromUserModelResp, err := l.svcCtx.UserService.GetUserModelById(context.Background(), &peerpb.GetUserModelByIdReq{
		Header: in.Header,
		UserId: fromUserId,
		Opt:    nil,
	})
	if err != nil {
		l.Errorf("get user model error: %v", err)
		return err
	}
	user := &usermodel.User{}
	err = utils.Json.Unmarshal(fromUserModelResp.UserModelJson, user)
	if err != nil {
		l.Errorf("unmarshal user model error: %v", err)
		return err
	}
	*fromUserInfo = *user
	return nil
}

func (l *FriendApplyLogic) toUserInfo_(in *peerpb.FriendApplyReq, toUserId string, toUserInfo *usermodel.User, toUserSettingMap *map[string]*usermodel.UserSetting) error {
	toUserModelResp, err := l.svcCtx.UserService.GetUserModelById(context.Background(), &peerpb.GetUserModelByIdReq{
		Header: in.Header,
		UserId: toUserId,
		Opt: &peerpb.GetUserModelByIdReq_Opt{
			WithUserSettings: true,
			UserSettingKeys: []string{
				usermodel.UserSettingKeyFriendApply,
			},
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

	userSettingMap := make(map[string]*usermodel.UserSetting)
	err = utils.Json.Unmarshal(toUserModelResp.UserSettingsJson, &userSettingMap)
	if err != nil {
		l.Errorf("unmarshal user setting error: %v", err)
		return err
	}
	*toUserSettingMap = userSettingMap

	return nil
}

func (l *FriendApplyLogic) verifyUserInfo_(in *peerpb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User) (*peerpb.FriendApplyResp, error) {
	//验证from用户状态
	accountStatus := fromUserInfo.GetAccountMap().Get(peerpb.AccountType_Status.String())
	switch accountStatus {
	case "", "0":
		//正常
	default:
		//账号状态异常
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyFromUserStatusError),
		}, nil
	}
	//验证to用户状态
	accountStatus = toUserInfo.GetAccountMap().Get(peerpb.AccountType_Status.String())
	switch accountStatus {
	case "", "0":
		//正常
	default:
		//账号状态异常
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyToUserStatusError),
		}, nil
	}
	//to用户是否注销
	if toUserInfo.DestroyTime == 0 {
		// 正常
	} else {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyToUserDestroy),
		}, nil
	}
	accountRole := fromUserInfo.GetAccountMap().Get(peerpb.AccountType_Role.String())
	//该角色是否允许添加好友
	if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Friend.AllowRoleApply) {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyFromUserRoleError),
		}, nil
	}
	accountRole = toUserInfo.GetAccountMap().Get(peerpb.AccountType_Role.String())
	//该角色是否允许被添加好友
	if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Friend.AllowRoleBeApplied) {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyToUserRoleError),
		}, nil
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyUserSetting_(in *peerpb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User, settingMap map[string]*usermodel.UserSetting) (*peerpb.FriendApplyResp, bool, error) {
	userSettingKeyFriendApplySetting := settingMap[usermodel.UserSettingKeyFriendApply]

	userSettingFriendApply := &usermodel.UserSettingFriendApply{}
	if userSettingKeyFriendApplySetting.V == "" {
		userSettingKeyFriendApplySetting.V = `{}`
	}
	err := utils.Json.Unmarshal([]byte(userSettingKeyFriendApplySetting.V), userSettingFriendApply)
	if err != nil {
		l.Errorf("unmarshal user setting friend apply error: %v", err)
		return nil, false, err
	}
	switch userSettingFriendApply.ApplyType {
	case usermodel.UserSettingFriendApplyTypeAny:
		return nil, true, nil
	case usermodel.UserSettingFriendApplyTypeVerifyMessage:
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestion:
		if in.Answer == nil {
			err = errors.New("friend apply answer question empty")
			l.Errorf("friend apply answer question empty")
			return &peerpb.FriendApplyResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAnswerQuestionError),
			}, false, err
		}
		answer := *in.Answer
		if answer != userSettingFriendApply.Answer {
			return &peerpb.FriendApplyResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAnswerQuestionError),
			}, false, errors.New("friend apply answer question error")
		}
		return nil, true, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestionAndConfirm:
		if in.Answer == nil || *in.Answer == "" {
			err = errors.New("friend apply answer question empty")
			l.Errorf("friend apply answer question empty")
			return &peerpb.FriendApplyResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAnswerQuestionError),
			}, false, err
		}
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestionAndConfirmWithRightAnswer:
		if in.Answer == nil {
			err = errors.New("friend apply answer question empty")
			l.Errorf("friend apply answer question empty")
			return &peerpb.FriendApplyResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAnswerQuestionError),
			}, false, err
		}
		answer := *in.Answer
		if answer != userSettingFriendApply.Answer {
			return &peerpb.FriendApplyResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAnswerQuestionError),
			}, false, errors.New("friend apply answer question error")
		}
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeNoOne:
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyNoOne),
		}, false, errors.New("friend apply no one")
	default:
		return nil, false, nil
	}
}

func (l *FriendApplyLogic) handle_(in *peerpb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User, settingMap map[string]*usermodel.UserSetting) error {
	message := ""
	if in.Message != nil {
		message = *in.Message
	}
	answer := ""
	if in.Answer != nil {
		answer = *in.Answer
	}
	friendApplyRecord := &friendmodel.FriendApplyRecord{
		ApplyId:        utils.Snowflake.String(),
		FromId:         fromUserInfo.UserId,
		ToId:           toUserInfo.UserId,
		Message:        message,
		Answer:         answer,
		ApplyTime:      primitive.NewDateTimeFromTime(time.Now()),
		Status:         friendmodel.FriendApplyStatusApplying,
		FromDeleteTime: 0,
		ToDeleteTime:   0,
	}
	_, err := l.svcCtx.FriendApplyRecordCollection.InsertOne(context.Background(), friendApplyRecord)
	if err != nil {
		l.Errorf("insert friend apply record error: %v", err)
		return err
	}
	go l.sendNotice_(in, fromUserInfo, toUserInfo, friendApplyRecord)
	return nil
}

func (l *FriendApplyLogic) verifyAreFriend_(in *peerpb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (*peerpb.FriendApplyResp, error) {
	yes, err := NewFriendApplyHandleLogic(context.Background(), l.svcCtx).AreFriends(fromUser.UserId, toUser.UserId)
	if err != nil {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyError),
		}, err
	}
	if yes {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.FriendApplyAreFriend),
		}, errors.New("friend apply are friend")
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyFriendLimit_(in *peerpb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (*peerpb.FriendApplyResp, error) {
	yes, msg, err := NewFriendApplyHandleLogic(context.Background(), l.svcCtx).IsFriendLimit(fromUser.UserId, toUser.UserId)
	if err != nil {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, msg),
		}, err
	}
	if yes {
		return &peerpb.FriendApplyResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, msg),
		}, errors.New("friend apply friend limit")
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyApplyExist_(in *peerpb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (bool, error) {
	count, err := l.svcCtx.FriendApplyRecordCollection.Find(context.Background(), bson.M{
		"fromId": fromUser.UserId,
		"toId":   toUser.UserId,
		"status": friendmodel.FriendApplyStatusApplying,
	}).Count()
	if err != nil {
		l.Errorf("verify apply exist error: %v", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (l *FriendApplyLogic) sendNotice_(in *peerpb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User, friendApplyRecord *friendmodel.FriendApplyRecord) {
	utils.Retry.Do(func() error {
		_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &peerpb.NoticeSendReq{
			Header: in.Header,
			Notices: []*peerpb.Message{&peerpb.Message{
				MessageId:        utils.Snowflake.String(),
				ConversationId:   peerpb.GetSingleChatConversationId(channelmodel.ConversationIdFriendHelper, toUser.UserId),
				ConversationType: peerpb.ConversationType_Single,
				Content:          utils.Json.MarshalToBytes(&peerpb.NoticeContentNewFriendRequest{}),
				ContentType:      peerpb.MessageContentType_NewFriendRequest,
				Option: &peerpb.Message_Option{
					StorageForServer: true,
					StorageForClient: true,
					CountUnread:      true,
				},
				Sender: &peerpb.Message_Sender{
					Id:         channelmodel.ConversationIdFriendHelper,
					SenderType: peerpb.SenderType_ChannelSender,
					Name:       channelmodel.FriendHelperChannel.Nickname,
					Avatar:     channelmodel.FriendHelperChannel.Avatar,
					Extra:      "",
				},
			}},
		})
		if err != nil {
			l.Errorf("send notice error: %v", err)
		}
		return err
	})
}
