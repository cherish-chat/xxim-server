package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNoticeConvIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNoticeConvIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNoticeConvIdsLogic {
	return &GetUserNoticeConvIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserNoticeConvIds 获取用户所有的通知号
func (l *GetUserNoticeConvIdsLogic) GetUserNoticeConvIds(in *pb.GetUserNoticeConvIdsReq) (*pb.GetUserNoticeConvIdsResp, error) {
	var convIds = noticemodel.DefaultConvIds
	// 获取用户的所有好友id
	{
		getFriendList, err := l.svcCtx.RelationService().GetFriendList(l.ctx, &pb.GetFriendListReq{
			CommonReq: &pb.CommonReq{
				UserId: in.UserId,
			},
			Page: &pb.Page{
				Page: 1,
				Size: 999999,
			},
			Opt: pb.GetFriendListReq_OnlyId,
		})
		if err != nil {
			l.Errorf("get friend list failed, err: %v", err)
			return &pb.GetUserNoticeConvIdsResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		friends := getFriendList.Ids
		for _, friend := range friends {
			convIds = append(convIds, noticemodel.ConvIdUser(friend))
		}
	}
	// 获取用户的所有群组id
	{
		getGroupList, err := l.svcCtx.GroupService().GetMyGroupList(l.ctx, &pb.GetMyGroupListReq{
			CommonReq: &pb.CommonReq{
				UserId: in.UserId,
			},
			Page: &pb.Page{
				Page: 1,
				Size: 999999,
			},
			Opt: pb.GetMyGroupListReq_ONLY_ID,
		})
		if err != nil {
			l.Errorf("get group list failed, err: %v", err)
			return &pb.GetUserNoticeConvIdsResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		groups := getGroupList.Ids
		for _, group := range groups {
			convIds = append(convIds, noticemodel.ConvIdGroup(group))
		}
	}
	return &pb.GetUserNoticeConvIdsResp{
		ConvIds: utils.Set(convIds),
	}, nil
}
