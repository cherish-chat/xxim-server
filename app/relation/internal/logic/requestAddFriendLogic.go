package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestAddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRequestAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestAddFriendLogic {
	return &RequestAddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RequestAddFriend 请求添加好友
func (l *RequestAddFriendLogic) RequestAddFriend(in *pb.RequestAddFriendReq) (*pb.RequestAddFriendResp, error) {
	// 看看我和对方是否已经是好友
	// 如果是好友，直接返回
	{
		var areFriendsResp *pb.AreFriendsResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "AreFriends", func(ctx context.Context) {
			areFriendsResp, err = NewAreFriendsLogic(ctx, l.svcCtx).AreFriends(&pb.AreFriendsReq{
				CommonReq: in.CommonReq,
				A:         in.CommonReq.UserId,
				BList:     []string{in.To},
			})
		})
		if err != nil {
			l.Errorf("AreFriends failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if is, ok := areFriendsResp.FriendList[in.To]; is && ok {
			// 已经是好友了
			return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "你们已经是好友了"))}, nil
		}
	}
	// 他是否把我拉黑
	{
		var areBlackListResp *pb.AreBlackListResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "AreBlackList", func(ctx context.Context) {
			areBlackListResp, err = NewAreBlackListLogic(ctx, l.svcCtx).AreBlackList(&pb.AreBlackListReq{
				CommonReq: in.CommonReq,
				A:         in.To,
				BList:     []string{in.CommonReq.UserId},
			})
		})
		if err != nil {
			l.Errorf("AreBlackList failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if is, ok := areBlackListResp.BlackList[in.CommonReq.UserId]; is && ok {
			// 已经被拉黑
			return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "对方已把你拉黑"))}, nil
		}
	}
	// 我的好友总数是否已达上限
	{
		var getFriendCountResp *pb.GetFriendCountResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "GetFriendCount", func(ctx context.Context) {
			getFriendCountResp, err = NewGetFriendCountLogic(ctx, l.svcCtx).GetFriendCount(&pb.GetFriendCountReq{
				CommonReq: in.CommonReq,
			})
		})
		if err != nil {
			l.Errorf("GetFriendCount failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if int64(getFriendCountResp.Count) >= utils.AnyToInt64(l.svcCtx.SystemConfigMgr.Get("app.friend_max_count")) {
			return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "好友数量已达上限"))}, nil
		}
	}
	// 对方的加好友设置
	{
		getUserSettingsResp, err := l.svcCtx.UserService().GetUserSettings(l.ctx, &pb.GetUserSettingsReq{
			CommonReq: &pb.CommonReq{UserId: in.To},
			Keys:      []pb.UserSettingKey{pb.UserSettingKey_HowToAddFriend, pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Answer},
		})
		if err != nil {
			l.Errorf("GetUserSettings failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		// 对方不允许任何人添加好友
		if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend)].Value == pb.UserSettingValue_HowToAddFriend_DontAllowAnyone {
			return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "对方不允许任何人添加好友"))}, nil
		} else
		// 对方允许任何人添加好友
		if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend)].Value == pb.UserSettingValue_HowToAddFriend_AllowAnyone {
			return l.allowAddFriend(in)
		} else
		// 对方需要正确回答问题
		if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend)].Value == pb.UserSettingValue_HowToAddFriend_NeedAnswerQuestionCorrectly {
			// 对方的问题的答案 和 in.Message 一致
			if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Answer)].Value != in.Message {
				return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "回答问题错误"))}, nil
			} else {
				return l.allowAddFriend(in)
			}
		} else
		// 对方需要验证
		if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend)].Value == pb.UserSettingValue_HowToAddFriend_NeedConfirm {
			return l.requestAddFriend(in)
		} else
		// 对方需要正确回答问题且需要验证
		if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend)].Value == pb.UserSettingValue_HowToAddFriend_NeedAnswerQuestionCorrectlyAndConfirm {
			// 对方的问题的答案 和 in.Message 一致
			if getUserSettingsResp.Settings[int32(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Answer)].Value != in.Message {
				return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "回答问题错误"))}, nil
			} else {
				return l.requestAddFriend(in)
			}
		}
	}
	return &pb.RequestAddFriendResp{}, nil
}

func (l *RequestAddFriendLogic) allowAddFriend(in *pb.RequestAddFriendReq) (*pb.RequestAddFriendResp, error) {
	var acceptAddFriendResp *pb.AcceptAddFriendResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "AcceptAddFriend", func(ctx context.Context) {
		acceptAddFriendResp, err = NewAcceptAddFriendLogic(ctx, l.svcCtx).AcceptAddFriend(&pb.AcceptAddFriendReq{CommonReq: &pb.CommonReq{
			UserId:   in.To,
			Language: in.CommonReq.Language,
		}})
	})
	if err != nil {
		l.Errorf("AcceptAddFriend failed, err: %v", err)
		return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.RequestAddFriendResp{CommonResp: acceptAddFriendResp.CommonResp}, nil
}

func (l *RequestAddFriendLogic) requestAddFriend(in *pb.RequestAddFriendReq) (*pb.RequestAddFriendResp, error) {
	now := time.Now().UnixMilli()
	extra := make([]*pb.RequestAddFriendExtra, 0)
	if in.Message != "" {
		extra = append(extra, &pb.RequestAddFriendExtra{
			UserId:  in.CommonReq.UserId,
			Content: in.Message,
		})
	}
	model := &relationmodel.RequestAddFriend{
		Id:         utils.GenId(),
		FromUserId: in.CommonReq.UserId,
		ToUserId:   in.To,
		Status:     pb.RequestAddFriendStatus_Unhandled,
		CreateTime: now,
		UpdateTime: now,
		Extra:      extra,
	}
	// 判断是否有没处理的 添加好友请求
	{
		var count int64
		err := l.svcCtx.Mysql().Model(model).Where("fromUserId = ? and toUserId = ? and status = ?", model.FromUserId, model.ToUserId, pb.RequestAddFriendStatus_Unhandled).Count(&count).Error
		if err != nil {
			l.Errorf("Find failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if count > 0 {
			return &pb.RequestAddFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "请勿重复添加好友"))}, nil
		}
	}
	// 插入 添加好友请求
	{
		err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			err := xorm.InsertOne(tx, model)
			if err != nil {
				l.Errorf("InsertOne failed, err: %v", err)
			}
			return err
		}, func(tx *gorm.DB) error {
			data := &pb.NoticeData{
				ConvId:         noticemodel.ConvId_FriendNotice,
				UnreadCount:    0,
				UnreadAbsolute: false,
				NoticeId:       fmt.Sprintf("%s", in.To),
				CreateTime:     "",
				Title:          "",
				ContentType:    1,
				Content:        nil,
				Options: &pb.NoticeData_Options{
					StorageForClient: false,
					UpdateConvMsg:    false,
					OnlinePushOnce:   false,
				},
				Ext: nil,
			}
			m := noticemodel.NoticeFromPB(data, false, in.To)
			err := m.Upsert(tx)
			if err != nil {
				l.Errorf("Upsert failed, err: %v", err)
			}
			return err
		})
		if err != nil {
			l.Errorf("Transaction failed, err: %v", err)
			return &pb.RequestAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	l.svcCtx.NoticeService().SendNoticeData(l.ctx, &pb.SendNoticeDataReq{
		CommonReq:   in.CommonReq,
		NoticeData:  &pb.NoticeData{NoticeId: fmt.Sprintf("%s", in.To)},
		UserId:      utils.AnyPtr(in.To),
		IsBroadcast: nil,
		Inserted:    utils.AnyPtr(true),
	})
	return &pb.RequestAddFriendResp{}, nil
}
