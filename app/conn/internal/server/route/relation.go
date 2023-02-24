package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterRelation(svcCtx *svc.ServiceContext) {
	// relation
	{
		// AcceptAddFriendReq AcceptAddFriendResp
		{
			route := conngateway.Route[*pb.AcceptAddFriendReq, *pb.AcceptAddFriendResp]{
				NewRequest: func() *pb.AcceptAddFriendReq {
					return &pb.AcceptAddFriendReq{}
				},
				Do: svcCtx.RelationService().AcceptAddFriend,
			}
			conngateway.AddRoute("/v1/relation/acceptAddFriend", route)
		}
		// GetMyFriendEventListReq GetMyFriendEventListResp
		{
			route := conngateway.Route[*pb.GetMyFriendEventListReq, *pb.GetMyFriendEventListResp]{
				NewRequest: func() *pb.GetMyFriendEventListReq {
					return &pb.GetMyFriendEventListReq{}
				},
				Do: svcCtx.RelationService().GetMyFriendEventList,
			}
			conngateway.AddRoute("/v1/relation/getMyFriendEventList", route)
		}
		// RequestAddFriendReq RequestAddFriendResp
		{
			route := conngateway.Route[*pb.RequestAddFriendReq, *pb.RequestAddFriendResp]{
				NewRequest: func() *pb.RequestAddFriendReq {
					return &pb.RequestAddFriendReq{}
				},
				Do: svcCtx.RelationService().RequestAddFriend,
			}
			conngateway.AddRoute("/v1/relation/requestAddFriend", route)
		}
		// GetFriendListReq GetFriendListResp
		{
			route := conngateway.Route[*pb.GetFriendListReq, *pb.GetFriendListResp]{
				NewRequest: func() *pb.GetFriendListReq {
					return &pb.GetFriendListReq{}
				},
				Do: svcCtx.RelationService().GetFriendList,
			}
			conngateway.AddRoute("/v1/relation/getFriendList", route)
		}
		// DeleteFriendReq DeleteFriendResp
		{
			route := conngateway.Route[*pb.DeleteFriendReq, *pb.DeleteFriendResp]{
				NewRequest: func() *pb.DeleteFriendReq {
					return &pb.DeleteFriendReq{}
				},
				Do: svcCtx.RelationService().DeleteFriend,
			}
			conngateway.AddRoute("/v1/relation/deleteFriend", route)
		}
	}
}
