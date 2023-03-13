package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"math"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFriendEventListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyFriendEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFriendEventListLogic {
	return &GetMyFriendEventListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyFriendEventListLogic) GetMyFriendEventList(in *pb.GetMyFriendEventListReq) (*pb.GetMyFriendEventListResp, error) {
	// 用户清空好友事件列表 最后一条事件的createTime
	setting, err := l.svcCtx.UserService().GetUserSettings(l.ctx, &pb.GetUserSettingsReq{
		CommonReq: &pb.CommonReq{UserId: in.CommonReq.UserId},
		Keys:      []pb.UserSettingKey{pb.UserSettingKey_FriendEventList_ClearTime},
	})
	if err != nil {
		l.Errorf("get user setting error: %v", err)
		return &pb.GetMyFriendEventListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	clearTime := utils.AnyToInt64(setting.Settings[int32(pb.UserSettingKey_FriendEventList_ClearTime)].Value)
	var list []*relationmodel.RequestAddFriend
	err = l.svcCtx.Mysql().Model(&relationmodel.RequestAddFriend{}).
		Where("(fromUserId = ? OR toUserId = ?) AND createTime > ? AND createTime < ?", in.CommonReq.UserId, in.CommonReq.UserId, clearTime, pb.PageIndex(in.PageIndex)).
		Order("createTime desc").
		Limit(20).
		Find(&list).Error
	if err != nil {
		l.Errorf("get friend event list error: %v", err)
		return &pb.GetMyFriendEventListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var respList []*pb.FriendEvent
	var pageIndex int64 = math.MaxInt64
	var otherUserIdList []string
	for _, friend := range list {
		if friend.FromUserId == in.CommonReq.UserId {
			otherUserIdList = append(otherUserIdList, friend.ToUserId)
		} else {
			otherUserIdList = append(otherUserIdList, friend.FromUserId)
		}
	}
	var userMap = make(map[string]*pb.UserBaseInfo)
	{
		otherUserIdList = utils.Set(otherUserIdList)
		userListResp, err := l.svcCtx.UserService().BatchGetUserBaseInfo(l.ctx, &pb.BatchGetUserBaseInfoReq{Ids: otherUserIdList})
		if err != nil {
			l.Errorf("get user info error: %v", err)
			return &pb.GetMyFriendEventListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		for _, user := range userListResp.UserBaseInfos {
			userMap[user.Id] = user
		}
	}
	for _, v := range list {
		var extra *pb.RequestAddFriendExtra
		if len(v.Extra) > 0 {
			extra = v.Extra[0]
		}
		otherId := v.FromUserId
		if v.FromUserId == in.CommonReq.UserId {
			otherId = v.ToUserId
		}
		info, ok := userMap[otherId]
		if !ok {
			info = &pb.UserBaseInfo{Id: otherId, Nickname: "用户已注销", Avatar: ""}
		}
		respList = append(respList, &pb.FriendEvent{
			FromUserId:    v.FromUserId,
			ToUserId:      v.ToUserId,
			OtherUserInfo: info,
			Status:        v.Status,
			CreateTime:    utils.AnyToString(v.CreateTime),
			UpdateTime:    utils.AnyToString(v.UpdateTime),
			Extra:         extra,
			RequestId:     v.Id,
		})
		if v.CreateTime < pageIndex {
			pageIndex = v.CreateTime
		}
	}
	return &pb.GetMyFriendEventListResp{
		FriendNotifyList: respList,
		PageIndex:        utils.AnyToString(pageIndex),
	}, nil
}
