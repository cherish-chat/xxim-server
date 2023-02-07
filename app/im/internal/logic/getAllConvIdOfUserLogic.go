package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllConvIdOfUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllConvIdOfUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllConvIdOfUserLogic {
	return &GetAllConvIdOfUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllConvIdOfUserLogic) GetAllConvIdOfUser(in *pb.GetAllConvIdOfUserReq) (*pb.GetAllConvIdOfUserResp, error) {
	var friendIds []string
	var groupIds []string
	var noticIds []string
	var convIds []string
	// 默认的
	{
		convIds = append(convIds, pb.HiddenConvIdCommand(), pb.HiddenConvIdFriendMember(), pb.HiddenConvIdGroupMember())
		noticIds = append(noticIds, pb.HiddenConvIdCommand(), pb.HiddenConvIdFriendMember(), pb.HiddenConvIdGroupMember())
	}
	// 获取用户订阅的好友列表
	{
		getFriendList, err := l.svcCtx.RelationService().GetFriendList(l.ctx, &pb.GetFriendListReq{
			CommonReq: &pb.CommonReq{
				UserId: in.UserId,
			},
			Page: &pb.Page{
				Page: 1,
				Size: 0,
			},
			Opt: pb.GetFriendListReq_OnlyId,
		})
		if err != nil {
			l.Errorf("get friend list error: %v", err)
			return &pb.GetAllConvIdOfUserResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		friendIds = getFriendList.Ids
		for _, id := range friendIds {
			convIds = append(convIds, pb.SingleConvId(in.UserId, id))
			convIds = append(convIds, pb.HiddenConvIdFriend(id))
			convIds = append(convIds, pb.HiddenConvId(pb.SingleConvId(in.UserId, id)))
			noticIds = append(noticIds, pb.HiddenConvId(pb.SingleConvId(in.UserId, id)))
			noticIds = append(noticIds, pb.HiddenConvIdFriend(id))
		}
	}
	// 获取用户订阅的群组列表
	{
		getMyGroupList, err := l.svcCtx.GroupService().GetMyGroupList(l.ctx, &pb.GetMyGroupListReq{
			CommonReq: &pb.CommonReq{
				UserId: in.UserId,
			},
			Page: &pb.Page{Page: 1},
			Filter: &pb.GetMyGroupListReq_Filter{
				FilterFold:   true,
				FilterShield: true,
			},
			Opt: pb.GetMyGroupListReq_ONLY_ID,
		})
		if err != nil {
			l.Errorf("get group list error: %v", err)
			return &pb.GetAllConvIdOfUserResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		groupIds = getMyGroupList.Ids
		for _, id := range groupIds {
			convIds = append(convIds, pb.GroupConvId(id))
			convIds = append(convIds, pb.HiddenConvIdGroup(id))
			noticIds = append(noticIds, pb.HiddenConvIdGroup(id))
		}
	}
	return &pb.GetAllConvIdOfUserResp{
		ConvIds:   convIds,
		GroupIds:  groupIds,
		FriendIds: friendIds,
		NoticeIds: noticIds,
	}, nil
}
