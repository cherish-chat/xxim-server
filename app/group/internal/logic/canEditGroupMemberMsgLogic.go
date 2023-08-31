package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanEditGroupMemberMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCanEditGroupMemberMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanEditGroupMemberMsgLogic {
	return &CanEditGroupMemberMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CanEditGroupMemberMsg 是否可以编辑群成员信息
func (l *CanEditGroupMemberMsgLogic) CanEditGroupMemberMsg(in *pb.CanEditGroupMemberMsgReq) (*pb.CanEditGroupMemberMsgResp, error) {
	var group *groupmodel.Group
	{
		var resp *pb.MapGroupByIdsResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "mapGroupByIds", func(ctx context.Context) {
			resp, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
				CommonReq: in.CommonReq,
				Ids:       []string{in.GroupId},
			})
		})
		if err != nil {
			l.Errorf("getGroupMemberInfoLogic err: %v", err)
			return &pb.CanEditGroupMemberMsgResp{CommonResp: resp.CommonResp}, err
		}
		value, ok := resp.GroupMap[in.GroupId]
		if !ok {
			return &pb.CanEditGroupMemberMsgResp{CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "群聊不存在"),
			)}, nil
		}
		group = groupmodel.GroupFromBytes(value)
	}
	// 获取自己的群成员信息 看是否有权
	myRole := pb.GroupRole_MEMBER
	{
		var resp *pb.GetGroupMemberInfoResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "getGroupMemberInfoLogic", func(ctx context.Context) {
			resp, err = NewGetGroupMemberInfoLogic(ctx, l.svcCtx).GetGroupMemberInfo(&pb.GetGroupMemberInfoReq{
				CommonReq: in.CommonReq,
				GroupId:   in.GroupId,
				MemberId:  in.CommonReq.UserId,
			})
		})
		if err != nil {
			l.Errorf("getGroupMemberInfoLogic err: %v", err)
			return &pb.CanEditGroupMemberMsgResp{CommonResp: resp.CommonResp}, err
		}
		myRole = resp.GroupMemberInfo.Role
	}

	// 判断是否有权限
	// 如果是member 只能修改自己的信息
	if myRole == pb.GroupRole_MEMBER && in.MemberId != in.CommonReq.UserId {
		return &pb.CanEditGroupMemberMsgResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "没有权限"),
			),
		}, nil
	}
	// 如果是管理员 不能修改群主的信息
	if myRole == pb.GroupRole_MANAGER && group.Owner == in.MemberId {
		return &pb.CanEditGroupMemberMsgResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "没有权限"),
			),
		}, nil
	}
	return &pb.CanEditGroupMemberMsgResp{}, nil
}
