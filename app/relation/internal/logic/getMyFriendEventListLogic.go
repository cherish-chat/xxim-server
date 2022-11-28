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
		CommonReq: &pb.CommonReq{Id: in.CommonReq.Id},
		Keys:      []pb.UserSettingKey{pb.UserSettingKey_FriendEventList_ClearTime},
	})
	if err != nil {
		l.Errorf("get user setting error: %v", err)
		return &pb.GetMyFriendEventListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	clearTime := utils.AnyToInt64(setting.Settings[int32(pb.UserSettingKey_FriendEventList_ClearTime)].Value)
	var list []*relationmodel.RequestAddFriend
	err = l.svcCtx.Mysql().Model(&relationmodel.RequestAddFriend{}).
		Where("(fromUserId = ? OR toUserId = ?) AND createTime > ? AND createTime < ?", in.CommonReq.Id, in.CommonReq.Id, clearTime, pb.PageIndex(in.PageIndex)).
		Order("createTime desc").
		Find(&list).Error
	if err != nil {
		l.Errorf("get friend event list error: %v", err)
		return &pb.GetMyFriendEventListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var respList []*pb.FriendEvent
	var pageIndex int64 = math.MaxInt64
	for _, v := range list {
		respList = append(respList, &pb.FriendEvent{
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Status:     v.Status,
			CreateTime: utils.AnyToString(v.CreateTime),
			UpdateTime: utils.AnyToString(v.UpdateTime),
			Extra:      v.Extra,
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
