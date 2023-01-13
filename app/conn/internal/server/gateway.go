package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/common/pb"
	"strconv"
)

func (s *ConnServer) registerGateway() {
	{
		// SendMsgListReq SendMsgListResp
		{
			route := conngateway.Route[*pb.SendMsgListReq, *pb.SendMsgListResp]{
				NewRequest: func() *pb.SendMsgListReq {
					return &pb.SendMsgListReq{}
				},
				Do: s.svcCtx.MsgService().SendMsgListAsync,
			}
			conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SendMsgList.Number())), route)
		}
		// BatchGetConvSeqReq BatchGetConvSeqResp
		{
			route := conngateway.Route[*pb.BatchGetConvSeqReq, *pb.BatchGetConvSeqResp]{
				NewRequest: func() *pb.BatchGetConvSeqReq {
					return &pb.BatchGetConvSeqReq{}
				},
				Do: s.svcCtx.MsgService().BatchGetConvSeq,
			}
			conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SyncConvSeq.Number())), route)
		}
		// BatchGetMsgListByConvIdReq GetMsgListResp
		{
			route := conngateway.Route[*pb.BatchGetMsgListByConvIdReq, *pb.GetMsgListResp]{
				NewRequest: func() *pb.BatchGetMsgListByConvIdReq {
					return &pb.BatchGetMsgListByConvIdReq{}
				},
				Do: s.svcCtx.MsgService().BatchGetMsgListByConvId,
			}
			conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SyncMsgList.Number())), route)
		}
		// GetMsgByIdReq GetMsgByIdResp
		{
			route := conngateway.Route[*pb.GetMsgByIdReq, *pb.GetMsgByIdResp]{
				NewRequest: func() *pb.GetMsgByIdReq {
					return &pb.GetMsgByIdReq{}
				},
				Do: s.svcCtx.MsgService().GetMsgById,
			}
			conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_GetMsgById.Number())), route)
		}
		// AckNoticeDataReq AckNoticeDataResp
		{
			route := conngateway.Route[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp]{
				NewRequest: func() *pb.AckNoticeDataReq {
					return &pb.AckNoticeDataReq{}
				},
				Do: s.svcCtx.NoticeService().AckNoticeData,
			}
			conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_AckNotice.Number())), route)
		}
	}
	// relation
	{
		// AcceptAddFriendReq AcceptAddFriendResp
		{
			route := conngateway.Route[*pb.AcceptAddFriendReq, *pb.AcceptAddFriendResp]{
				NewRequest: func() *pb.AcceptAddFriendReq {
					return &pb.AcceptAddFriendReq{}
				},
				Do: s.svcCtx.RelationService().AcceptAddFriend,
			}
			conngateway.AddRoute("/v1/relation/acceptAddFriend", route)
		}
		// GetMyFriendEventListReq GetMyFriendEventListResp
		{
			route := conngateway.Route[*pb.GetMyFriendEventListReq, *pb.GetMyFriendEventListResp]{
				NewRequest: func() *pb.GetMyFriendEventListReq {
					return &pb.GetMyFriendEventListReq{}
				},
				Do: s.svcCtx.RelationService().GetMyFriendEventList,
			}
			conngateway.AddRoute("/v1/relation/getMyFriendEventList", route)
		}
		// RequestAddFriendReq RequestAddFriendResp
		{
			route := conngateway.Route[*pb.RequestAddFriendReq, *pb.RequestAddFriendResp]{
				NewRequest: func() *pb.RequestAddFriendReq {
					return &pb.RequestAddFriendReq{}
				},
				Do: s.svcCtx.RelationService().RequestAddFriend,
			}
			conngateway.AddRoute("/v1/relation/requestAddFriend", route)
		}
		// GetFriendListReq GetFriendListResp
		{
			route := conngateway.Route[*pb.GetFriendListReq, *pb.GetFriendListResp]{
				NewRequest: func() *pb.GetFriendListReq {
					return &pb.GetFriendListReq{}
				},
				Do: s.svcCtx.RelationService().GetFriendList,
			}
			conngateway.AddRoute("/v1/relation/getFriendList", route)
		}
	}
	// user
	{
		// SearchUsersByKeywordReq SearchUsersByKeywordResp
		{
			route := conngateway.Route[*pb.SearchUsersByKeywordReq, *pb.SearchUsersByKeywordResp]{
				NewRequest: func() *pb.SearchUsersByKeywordReq {
					return &pb.SearchUsersByKeywordReq{}
				},
				Do: s.svcCtx.UserService().SearchUsersByKeyword,
			}
			conngateway.AddRoute("/v1/user/searchUsersByKeyword", route)
		}
		// GetUserSettingsReq GetUserSettingsResp
		{
			route := conngateway.Route[*pb.GetUserSettingsReq, *pb.GetUserSettingsResp]{
				NewRequest: func() *pb.GetUserSettingsReq {
					return &pb.GetUserSettingsReq{}
				},
				Do: s.svcCtx.UserService().GetUserSettings,
			}
			conngateway.AddRoute("/v1/user/getUserSettings", route)
		}
		// SetUserSettingsReq SetUserSettingsResp
		{
			route := conngateway.Route[*pb.SetUserSettingsReq, *pb.SetUserSettingsResp]{
				NewRequest: func() *pb.SetUserSettingsReq {
					return &pb.SetUserSettingsReq{}
				},
				Do: s.svcCtx.UserService().SetUserSettings,
			}
			conngateway.AddRoute("/v1/user/setUserSettings", route)
		}
		// UpdateUserInfoReq UpdateUserInfoResp
		{
			route := conngateway.Route[*pb.UpdateUserInfoReq, *pb.UpdateUserInfoResp]{
				NewRequest: func() *pb.UpdateUserInfoReq {
					return &pb.UpdateUserInfoReq{}
				},
				Do: s.svcCtx.UserService().UpdateUserInfo,
			}
			conngateway.AddRoute("/v1/user/updateUserInfo", route)
		}
	}
}
