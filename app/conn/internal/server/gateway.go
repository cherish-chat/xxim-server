package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic"
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func (s *ConnServer) registerRelation() {
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
}

func (s *ConnServer) registerUser() {
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
	// 白名单
	{
		// 登录
		{
			route := conngateway.Route[*pb.LoginReq, *pb.LoginResp]{
				NewRequest: func() *pb.LoginReq {
					return &pb.LoginReq{}
				},
				Do: s.svcCtx.UserService().Login,
			}
			conngateway.AddRoute("/v1/user/white/login", route)
		}
		// 确认注册
		{
			route := conngateway.Route[*pb.ConfirmRegisterReq, *pb.ConfirmRegisterResp]{
				NewRequest: func() *pb.ConfirmRegisterReq {
					return &pb.ConfirmRegisterReq{}
				},
				Do: s.svcCtx.UserService().ConfirmRegister,
			}
			conngateway.AddRoute("/v1/user/white/confirmRegister", route)
		}
	}
}

func (s *ConnServer) registerGroup() {
	// group
	{
		// CreateGroupReq CreateGroupResp
		{
			route := conngateway.Route[*pb.CreateGroupReq, *pb.CreateGroupResp]{
				NewRequest: func() *pb.CreateGroupReq {
					return &pb.CreateGroupReq{}
				},
				Do: s.svcCtx.GroupService().CreateGroup,
			}
			conngateway.AddRoute("/v1/group/createGroup", route)
		}
		// GetMyGroupListReq GetMyGroupListResp
		{
			route := conngateway.Route[*pb.GetMyGroupListReq, *pb.GetMyGroupListResp]{
				NewRequest: func() *pb.GetMyGroupListReq {
					return &pb.GetMyGroupListReq{}
				},
				Do: s.svcCtx.GroupService().GetMyGroupList,
			}
			conngateway.AddRoute("/v1/group/getMyGroupList", route)
		}
		// SetGroupMemberInfoReq SetGroupMemberInfoResp
		{
			route := conngateway.Route[*pb.SetGroupMemberInfoReq, *pb.SetGroupMemberInfoResp]{
				NewRequest: func() *pb.SetGroupMemberInfoReq {
					return &pb.SetGroupMemberInfoReq{}
				},
				Do: s.svcCtx.GroupService().SetGroupMemberInfo,
			}
			conngateway.AddRoute("/v1/group/setGroupMemberInfo", route)
		}
		// GetGroupMemberInfoReq GetGroupMemberInfoResp
		{
			route := conngateway.Route[*pb.GetGroupMemberInfoReq, *pb.GetGroupMemberInfoResp]{
				NewRequest: func() *pb.GetGroupMemberInfoReq {
					return &pb.GetGroupMemberInfoReq{}
				},
				Do: s.svcCtx.GroupService().GetGroupMemberInfo,
			}
			conngateway.AddRoute("/v1/group/getGroupMemberInfo", route)
		}
		// ApplyToBeGroupMemberReq ApplyToBeGroupMemberResp
		{
			route := conngateway.Route[*pb.ApplyToBeGroupMemberReq, *pb.ApplyToBeGroupMemberResp]{
				NewRequest: func() *pb.ApplyToBeGroupMemberReq {
					return &pb.ApplyToBeGroupMemberReq{}
				},
				Do: s.svcCtx.GroupService().ApplyToBeGroupMember,
			}
			conngateway.AddRoute("/v1/group/applyToBeGroupMember", route)
		}
		// HandleGroupApplyReq HandleGroupApplyResp
		{
			route := conngateway.Route[*pb.HandleGroupApplyReq, *pb.HandleGroupApplyResp]{
				NewRequest: func() *pb.HandleGroupApplyReq {
					return &pb.HandleGroupApplyReq{}
				},
				Do: s.svcCtx.GroupService().HandleGroupApply,
			}
			conngateway.AddRoute("/v1/group/handleGroupApply", route)
		}
		// KickGroupMemberReq KickGroupMemberResp
		{
			route := conngateway.Route[*pb.KickGroupMemberReq, *pb.KickGroupMemberResp]{
				NewRequest: func() *pb.KickGroupMemberReq {
					return &pb.KickGroupMemberReq{}
				},
				Do: s.svcCtx.GroupService().KickGroupMember,
			}
			conngateway.AddRoute("/v1/group/kickGroupMember", route)
		}
	}
}

func (s *ConnServer) registerGateway() {
	// 自带的
	{
		// 设置连接参数
		{
			route := conngateway.Route[*pb.SetCxnParamsReq, *pb.SetCxnParamsResp]{
				NewRequest: func() *pb.SetCxnParamsReq {
					return &pb.SetCxnParamsReq{}
				},
				Do:       logic.NewSetConnParamsLogic(s.svcCtx).SetConnParams,
				Callback: logic.NewSetConnParamsLogic(s.svcCtx).Callback,
			}
			conngateway.AddRoute("/v1/conn/white/setConnParams", route)
		}
		// 设置userId和token
		{
			route := conngateway.Route[*pb.SetUserParamsReq, *pb.SetUserParamsResp]{
				NewRequest: func() *pb.SetUserParamsReq {
					return &pb.SetUserParamsReq{}
				},
				Do:       logic.NewSetUserParamsLogic(s.svcCtx).SetUserParams,
				Callback: logic.NewSetUserParamsLogic(s.svcCtx).Callback,
			}
			conngateway.AddRoute("/v1/conn/white/setUserParams", route)
		}
	}
	s.registerMsg()
	s.registerNotice()
	s.registerRelation()
	s.registerUser()
	s.registerGroup()
	conngateway.PrintRoutes()
}

func (s *ConnServer) registerMsg() {
	// SendMsgListReq SendMsgListResp
	{
		route := conngateway.Route[*pb.SendMsgListReq, *pb.SendMsgListResp]{
			NewRequest: func() *pb.SendMsgListReq {
				return &pb.SendMsgListReq{}
			},
			Do: s.svcCtx.MsgService().SendMsgListAsync,
		}
		conngateway.AddRoute("/v1/msg/sendMsgList", route)
	}
	// BatchGetConvSeqReq BatchGetConvSeqResp
	{
		route := conngateway.Route[*pb.BatchGetConvSeqReq, *pb.BatchGetConvSeqResp]{
			NewRequest: func() *pb.BatchGetConvSeqReq {
				return &pb.BatchGetConvSeqReq{}
			},
			Do: s.svcCtx.MsgService().BatchGetConvSeq,
		}
		conngateway.AddRoute("/v1/msg/batchGetConvSeq", route)
	}
	// BatchGetMsgListByConvIdReq GetMsgListResp
	{
		route := conngateway.Route[*pb.BatchGetMsgListByConvIdReq, *pb.GetMsgListResp]{
			NewRequest: func() *pb.BatchGetMsgListByConvIdReq {
				return &pb.BatchGetMsgListByConvIdReq{}
			},
			Do: s.svcCtx.MsgService().BatchGetMsgListByConvId,
		}
		conngateway.AddRoute("/v1/msg/batchGetMsgListByConvId", route)
	}
	// GetMsgByIdReq GetMsgByIdResp
	{
		route := conngateway.Route[*pb.GetMsgByIdReq, *pb.GetMsgByIdResp]{
			NewRequest: func() *pb.GetMsgByIdReq {
				return &pb.GetMsgByIdReq{}
			},
			Do: s.svcCtx.MsgService().GetMsgById,
		}
		conngateway.AddRoute("/v1/msg/getMsgById", route)
	}
}

func (s *ConnServer) registerNotice() {
	// AckNoticeDataReq AckNoticeDataResp
	{
		route := conngateway.Route[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp]{
			NewRequest: func() *pb.AckNoticeDataReq {
				return &pb.AckNoticeDataReq{}
			},
			Do: s.svcCtx.NoticeService().AckNoticeData,
		}
		conngateway.AddRoute("/v1/notice/ackNoticeData", route)
	}
}
