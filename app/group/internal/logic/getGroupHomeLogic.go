package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupHomeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupHomeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupHomeLogic {
	return &GetGroupHomeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupHome 获取群聊首页
func (l *GetGroupHomeLogic) GetGroupHome(in *pb.GetGroupHomeReq) (*pb.GetGroupHomeResp, error) {
	var groupByIds *pb.MapGroupByIdsResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "", func(ctx context.Context) {
		groupByIds, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
			CommonReq: in.CommonReq,
			Ids:       []string{in.GroupId},
		})
	})
	if err != nil {
		l.Errorf("MapGroupByIds error: %s", err.Error())
		return &pb.GetGroupHomeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	bytes, ok := groupByIds.GroupMap[in.GroupId]
	if !ok {
		l.Errorf("group not found: %s", in.GroupId)
		return &pb.GetGroupHomeResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "群聊不存在"))}, nil
	}
	group := groupmodel.GroupFromBytes(bytes)
	logx.Infof("group: %+v", group)
	getGroupMemberListResp, err := NewGetGroupMemberListLogic(l.ctx, l.svcCtx).GetGroupMemberList(&pb.GetGroupMemberListReq{
		CommonReq: in.CommonReq,
		GroupId:   in.GroupId,
		Page: &pb.Page{
			Page: 1,
			Size: 100,
		},
		Filter: &pb.GetGroupMemberListReq_GetGroupMemberListFilter{
			NoDisturb:  nil,
			OnlyOwner:  utils.AnyPtr(true),
			OnlyAdmin:  utils.AnyPtr(true),
			OnlyMember: nil,
		},
		Opt: &pb.GetGroupMemberListReq_GetGroupMemberListOpt{},
	})
	if err != nil {
		l.Errorf("GetGroupMemberList error: %v", err)
		return nil, err
	}
	var admins []*pb.UserBaseInfo
	for _, info := range getGroupMemberListResp.GetGroupMemberList() {
		if info.Role == pb.GroupRole_OWNER {
			admins = append(admins, info.UserBaseInfo)
		}
	}
	for _, info := range getGroupMemberListResp.GetGroupMemberList() {
		if info.Role == pb.GroupRole_OWNER {
			continue
		}
		admins = append(admins, info.UserBaseInfo)
	}
	return &pb.GetGroupHomeResp{
		CommonResp:         pb.NewSuccessResp(),
		GroupId:            group.Id,
		Name:               group.Name,
		Avatar:             group.Avatar,
		CreatedAt:          utils.TimeFormat(group.CreateTime),
		MemberCount:        int32(group.MemberCount),
		Introduction:       group.Description,
		Owner:              group.Owner,
		DismissTime:        group.DismissTime,
		AllMute:            group.AllMute,
		MemberCanAddFriend: group.MemberCanAddFriend,
		MemberStatistics:   nil,
		Admins:             admins,
		CanAddMember:       group.CanAddMember,
	}, nil
}
