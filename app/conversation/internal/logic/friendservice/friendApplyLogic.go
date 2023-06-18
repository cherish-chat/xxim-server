package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/mr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
func (l *FriendApplyLogic) FriendApply(in *pb.FriendApplyReq) (*pb.FriendApplyResp, error) {
	//预先验证
	{
		//是否允许申请加好友
		if !utils.EnumInSlice(in.Header.Platform, l.svcCtx.Config.Friend.AllowPlatform) {
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_not_allow_platform")),
			}, nil
		}
		//是不是相同id
		if in.Header.UserId == in.ToUserId {
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_same_id")),
			}, nil
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
		return &pb.FriendApplyResp{}, err
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
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_already")),
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
		err = NewFriendApplyHandleLogic(l.ctx, l.svcCtx).AddFriend(fromUserInfo.UserId, toUserInfo.UserId)
		return &pb.FriendApplyResp{}, err
	}

	//处理
	err = l.handle_(in, fromUserInfo, toUserInfo, toUserSettingMap)
	if err != nil {
		return &pb.FriendApplyResp{}, err
	}
	return &pb.FriendApplyResp{}, nil
}

func (l *FriendApplyLogic) fromUserInfo_(in *pb.FriendApplyReq, fromUserId string, fromUserInfo *usermodel.User) error {
	fromUserModelResp, err := l.svcCtx.InfoService.GetUserModelById(l.ctx, &pb.GetUserModelByIdReq{
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

func (l *FriendApplyLogic) toUserInfo_(in *pb.FriendApplyReq, toUserId string, toUserInfo *usermodel.User, toUserSettingMap *map[string]*usermodel.UserSetting) error {
	toUserModelResp, err := l.svcCtx.InfoService.GetUserModelById(l.ctx, &pb.GetUserModelByIdReq{
		Header: in.Header,
		UserId: toUserId,
		Opt: &pb.GetUserModelByIdReq_Opt{
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

func (l *FriendApplyLogic) verifyUserInfo_(in *pb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User) (*pb.FriendApplyResp, error) {
	//验证from用户状态
	accountStatus := fromUserInfo.GetAccountMap().Get(pb.AccountTypeStatus)
	switch accountStatus {
	case "", "0":
		//正常
	default:
		//账号状态异常
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_from_user_status_error")),
		}, nil
	}
	//验证to用户状态
	accountStatus = toUserInfo.GetAccountMap().Get(pb.AccountTypeStatus)
	switch accountStatus {
	case "", "0":
		//正常
	default:
		//账号状态异常
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_from_user_status_error")),
		}, nil
	}
	//to用户是否注销
	if toUserInfo.DestroyTime == 0 {
		// 正常
	} else {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_to_user_destroy")),
		}, nil
	}
	accountRole := fromUserInfo.GetAccountMap().Get(pb.AccountTypeRole)
	//该角色是否允许添加好友
	if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Friend.AllowRoleApply) {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_from_user_role_error")),
		}, nil
	}
	accountRole = toUserInfo.GetAccountMap().Get(pb.AccountTypeRole)
	//该角色是否允许被添加好友
	if !utils.AnyInSlice(accountRole, l.svcCtx.Config.Friend.AllowRoleBeApplied) {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_to_user_role_error")),
		}, nil
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyUserSetting_(in *pb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User, settingMap map[string]*usermodel.UserSetting) (*pb.FriendApplyResp, bool, error) {
	userSettingKeyFriendApplySetting := settingMap[usermodel.UserSettingKeyFriendApply]

	userSettingFriendApply := &usermodel.UserSettingFriendApply{}
	if userSettingKeyFriendApplySetting.V == "" {
		userSettingKeyFriendApplySetting.V = l.svcCtx.Config.Friend.DefaultApplySetting
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
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_answer_question_error")),
			}, false, err
		}
		answer := *in.Answer
		if answer != userSettingFriendApply.Answer {
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_answer_question_error")),
			}, false, errors.New("friend apply answer question error")
		}
		return nil, true, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestionAndConfirm:
		if in.Answer == nil || *in.Answer == "" {
			err = errors.New("friend apply answer question empty")
			l.Errorf("friend apply answer question empty")
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_answer_question_error")),
			}, false, err
		}
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestionAndConfirmWithRightAnswer:
		if in.Answer == nil {
			err = errors.New("friend apply answer question empty")
			l.Errorf("friend apply answer question empty")
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_answer_question_error")),
			}, false, err
		}
		answer := *in.Answer
		if answer != userSettingFriendApply.Answer {
			return &pb.FriendApplyResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_answer_question_error")),
			}, false, errors.New("friend apply answer question error")
		}
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeNoOne:
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_no_one")),
		}, false, errors.New("friend apply no one")
	default:
		return nil, false, nil
	}
}

func (l *FriendApplyLogic) handle_(in *pb.FriendApplyReq, fromUserInfo *usermodel.User, toUserInfo *usermodel.User, settingMap map[string]*usermodel.UserSetting) error {
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
	_, err := l.svcCtx.FriendApplyRecordCollection.InsertOne(l.ctx, friendApplyRecord)
	if err != nil {
		l.Errorf("insert friend apply record error: %v", err)
		return err
	}
	go l.sendNotice_(in, fromUserInfo, toUserInfo, friendApplyRecord)
	return nil
}

func (l *FriendApplyLogic) verifyAreFriend_(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (*pb.FriendApplyResp, error) {
	yes, err := NewFriendApplyHandleLogic(l.ctx, l.svcCtx).AreFriends(fromUser.UserId, toUser.UserId)
	if err != nil {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_error")),
		}, err
	}
	if yes {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_are_friend")),
		}, errors.New("friend apply are friend")
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyFriendLimit_(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (*pb.FriendApplyResp, error) {
	yes, msg, err := NewFriendApplyHandleLogic(l.ctx, l.svcCtx).IsFriendLimit(fromUser.UserId, toUser.UserId)
	if err != nil {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, msg),
		}, err
	}
	if yes {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, msg),
		}, errors.New("friend apply friend limit")
	}
	return nil, nil
}

func (l *FriendApplyLogic) verifyApplyExist_(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (bool, error) {
	count, err := l.svcCtx.FriendApplyRecordCollection.Find(l.ctx, bson.M{
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

func (l *FriendApplyLogic) sendNotice_(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User, friendApplyRecord *friendmodel.FriendApplyRecord) {
	notice := &noticemodel.BroadcastNotice{
		NoticeId:         utils.Snowflake.String(),
		ConversationId:   subscriptionmodel.ConversationIdFriendNotification,
		ConversationType: pb.ConversationType_Subscription,
		Content:          utils.Json.MarshalToString(&pb.NoticeContentNewFriendRequest{}),
		ContentType:      pb.NoticeContentType_NewFriendRequest,
		UpdateTime:       primitive.NewDateTimeFromTime(time.Now()),
	}
	utils.Retry.Do(func() error {
		_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &pb.NoticeSendReq{
			Header:    in.Header,
			Notice:    notice.ToPb(),
			UserIds:   []string{toUser.UserId},
			Broadcast: false,
		})
		if err != nil {
			l.Errorf("send notice error: %v", err)
		}
		return err
	})
}
