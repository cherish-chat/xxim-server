package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterGroup(svcCtx *svc.ServiceContext) {
	// group
	{
		// CreateGroupReq CreateGroupResp
		{
			route := conngateway.Route[*pb.CreateGroupReq, *pb.CreateGroupResp]{
				NewRequest: func() *pb.CreateGroupReq {
					return &pb.CreateGroupReq{}
				},
				Do: svcCtx.GroupService().CreateGroup,
			}
			conngateway.AddRoute("/v1/group/createGroup", route)
		}
		// GetMyGroupListReq GetMyGroupListResp
		{
			route := conngateway.Route[*pb.GetMyGroupListReq, *pb.GetMyGroupListResp]{
				NewRequest: func() *pb.GetMyGroupListReq {
					return &pb.GetMyGroupListReq{}
				},
				Do: svcCtx.GroupService().GetMyGroupList,
			}
			conngateway.AddRoute("/v1/group/getMyGroupList", route)
		}
		// SetGroupMemberInfoReq SetGroupMemberInfoResp
		{
			route := conngateway.Route[*pb.SetGroupMemberInfoReq, *pb.SetGroupMemberInfoResp]{
				NewRequest: func() *pb.SetGroupMemberInfoReq {
					return &pb.SetGroupMemberInfoReq{}
				},
				Do: svcCtx.GroupService().SetGroupMemberInfo,
			}
			conngateway.AddRoute("/v1/group/setGroupMemberInfo", route)
		}
		// GetGroupMemberInfoReq GetGroupMemberInfoResp
		{
			route := conngateway.Route[*pb.GetGroupMemberInfoReq, *pb.GetGroupMemberInfoResp]{
				NewRequest: func() *pb.GetGroupMemberInfoReq {
					return &pb.GetGroupMemberInfoReq{}
				},
				Do: svcCtx.GroupService().GetGroupMemberInfo,
			}
			conngateway.AddRoute("/v1/group/getGroupMemberInfo", route)
		}
		// ApplyToBeGroupMemberReq ApplyToBeGroupMemberResp
		{
			route := conngateway.Route[*pb.ApplyToBeGroupMemberReq, *pb.ApplyToBeGroupMemberResp]{
				NewRequest: func() *pb.ApplyToBeGroupMemberReq {
					return &pb.ApplyToBeGroupMemberReq{}
				},
				Do: svcCtx.GroupService().ApplyToBeGroupMember,
			}
			conngateway.AddRoute("/v1/group/applyToBeGroupMember", route)
		}
		// HandleGroupApplyReq HandleGroupApplyResp
		{
			route := conngateway.Route[*pb.HandleGroupApplyReq, *pb.HandleGroupApplyResp]{
				NewRequest: func() *pb.HandleGroupApplyReq {
					return &pb.HandleGroupApplyReq{}
				},
				Do: svcCtx.GroupService().HandleGroupApply,
			}
			conngateway.AddRoute("/v1/group/handleGroupApply", route)
		}
		// GetGroupApplyListReq GetGroupApplyListResp
		{
			route := conngateway.Route[*pb.GetGroupApplyListReq, *pb.GetGroupApplyListResp]{
				NewRequest: func() *pb.GetGroupApplyListReq {
					return &pb.GetGroupApplyListReq{}
				},
				Do: svcCtx.GroupService().GetGroupApplyList,
			}
			conngateway.AddRoute("/v1/group/getGroupApplyList", route)
		}
		// KickGroupMemberReq KickGroupMemberResp
		{
			route := conngateway.Route[*pb.KickGroupMemberReq, *pb.KickGroupMemberResp]{
				NewRequest: func() *pb.KickGroupMemberReq {
					return &pb.KickGroupMemberReq{}
				},
				Do: svcCtx.GroupService().KickGroupMember,
			}
			conngateway.AddRoute("/v1/group/kickGroupMember", route)
		}
		// SearchGroupsByKeywordReq SearchGroupsByKeywordResp
		{
			route := conngateway.Route[*pb.SearchGroupsByKeywordReq, *pb.SearchGroupsByKeywordResp]{
				NewRequest: func() *pb.SearchGroupsByKeywordReq {
					return &pb.SearchGroupsByKeywordReq{}
				},
				Do: svcCtx.GroupService().SearchGroupsByKeyword,
			}
			conngateway.AddRoute("/v1/group/searchGroupsByKeyword", route)
		}
		// GetGroupHomeReq GetGroupHomeResp
		{
			route := conngateway.Route[*pb.GetGroupHomeReq, *pb.GetGroupHomeResp]{
				NewRequest: func() *pb.GetGroupHomeReq {
					return &pb.GetGroupHomeReq{}
				},
				Do: svcCtx.GroupService().GetGroupHome,
			}
			conngateway.AddRoute("/v1/group/getGroupHome", route)
		}
		// GetGroupMemberListReq GetGroupMemberListResp
		{
			route := conngateway.Route[*pb.GetGroupMemberListReq, *pb.GetGroupMemberListResp]{
				NewRequest: func() *pb.GetGroupMemberListReq {
					return &pb.GetGroupMemberListReq{}
				},
				Do: svcCtx.GroupService().GetGroupMemberList,
			}
			conngateway.AddRoute("/v1/group/getGroupMemberList", route)
		}
	}
}
