package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/mr"

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

	//是否允许申请加好友
	if !utils.EnumInSlice(in.Header.Platform, l.svcCtx.Config.Friend.AllowPlatform) {
		return &pb.FriendApplyResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "friend_apply_not_allow_platform")),
		}, nil
	}
	//验证用户信息
	friendApplyResp, err := l.verifyUserInfo_(in, fromUserInfo, toUserInfo)
	if err != nil {
		return friendApplyResp, err
	}
	//验证用户设置
	friendApplyResp, skipApply, err := l.verifyUserSetting_(in, fromUserInfo, toUserInfo, toUserSettingMap)
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

	//跳过申请
	if skipApply {
		err = NewFriendApplyHandleLogic(l.ctx, l.svcCtx).AddFriend(in, fromUserInfo, toUserInfo)
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
		//TODO: 账号状态异常
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
		//TODO: 账号状态异常
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
		// TODO: 发送验证消息
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestion:
		// TODO: 验证答案
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
		// TODO: 填写答案，并发送验证消息
		return nil, false, nil
	case usermodel.UserSettingFriendApplyTypeAnswerQuestionAndConfirmWithRightAnswer:
		// TODO：验证答案，并发送验证消息
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
	//TODO: handle your logic here and delete this line
	return nil
}

func (l *FriendApplyLogic) verifyAreFriend_(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (*pb.FriendApplyResp, error) {
	yes, err := NewFriendApplyHandleLogic(l.ctx, l.svcCtx).AreFriends(in, fromUser, toUser)
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
	yes, msg, err := NewFriendApplyHandleLogic(l.ctx, l.svcCtx).IsFriendLimit(in, fromUser, toUser)
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
