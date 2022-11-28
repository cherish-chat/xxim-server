package relationhandler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/wrapper"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/grpc"
	"time"
)

// RequestAddFriendConfig ...
func RequestAddFriendConfig[REQ *pb.RequestAddFriendReq, RESP *pb.RequestAddFriendResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.RequestAddFriendReq, *pb.RequestAddFriendResp] {
	return wrapper.Config[*pb.RequestAddFriendReq, *pb.RequestAddFriendResp]{
		Do: func(ctx context.Context, in *pb.RequestAddFriendReq, opts ...grpc.CallOption) (*pb.RequestAddFriendResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().RequestAddFriend(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "RequestAddFriend", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.RequestAddFriendReq {
			return &pb.RequestAddFriendReq{}
		},
	}
}

// AcceptAddFriendConfig ...
func AcceptAddFriendConfig[REQ *pb.AcceptAddFriendReq, RESP *pb.AcceptAddFriendResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.AcceptAddFriendReq, *pb.AcceptAddFriendResp] {
	return wrapper.Config[*pb.AcceptAddFriendReq, *pb.AcceptAddFriendResp]{
		Do: func(ctx context.Context, in *pb.AcceptAddFriendReq, opts ...grpc.CallOption) (*pb.AcceptAddFriendResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().AcceptAddFriend(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "AcceptAddFriend", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.AcceptAddFriendReq {
			return &pb.AcceptAddFriendReq{}
		},
	}
}

// RejectAddFriendConfig ...
func RejectAddFriendConfig[REQ *pb.RejectAddFriendReq, RESP *pb.RejectAddFriendResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.RejectAddFriendReq, *pb.RejectAddFriendResp] {
	return wrapper.Config[*pb.RejectAddFriendReq, *pb.RejectAddFriendResp]{
		Do: func(ctx context.Context, in *pb.RejectAddFriendReq, opts ...grpc.CallOption) (*pb.RejectAddFriendResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().RejectAddFriend(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "RejectAddFriend", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.RejectAddFriendReq {
			return &pb.RejectAddFriendReq{}
		},
	}
}

// BlockUserConfig ...
func BlockUserConfig[REQ *pb.BlockUserReq, RESP *pb.BlockUserResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.BlockUserReq, *pb.BlockUserResp] {
	return wrapper.Config[*pb.BlockUserReq, *pb.BlockUserResp]{
		Do: func(ctx context.Context, in *pb.BlockUserReq, opts ...grpc.CallOption) (*pb.BlockUserResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().BlockUser(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "BlockUser", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.BlockUserReq {
			return &pb.BlockUserReq{}
		},
	}
}

// DeleteBlockUserConfig ...
func DeleteBlockUserConfig[REQ *pb.DeleteBlockUserReq, RESP *pb.DeleteBlockUserResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.DeleteBlockUserReq, *pb.DeleteBlockUserResp] {
	return wrapper.Config[*pb.DeleteBlockUserReq, *pb.DeleteBlockUserResp]{
		Do: func(ctx context.Context, in *pb.DeleteBlockUserReq, opts ...grpc.CallOption) (*pb.DeleteBlockUserResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().DeleteBlockUser(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "DeleteBlockUser", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.DeleteBlockUserReq {
			return &pb.DeleteBlockUserReq{}
		},
	}
}

// DeleteFriendConfig ...
func DeleteFriendConfig[REQ *pb.DeleteFriendReq, RESP *pb.DeleteFriendResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.DeleteFriendReq, *pb.DeleteFriendResp] {
	return wrapper.Config[*pb.DeleteFriendReq, *pb.DeleteFriendResp]{
		Do: func(ctx context.Context, in *pb.DeleteFriendReq, opts ...grpc.CallOption) (*pb.DeleteFriendResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().DeleteFriend(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "DeleteFriend", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.DeleteFriendReq {
			return &pb.DeleteFriendReq{}
		},
	}
}

// SetSingleChatSettingConfig ...
func SetSingleChatSettingConfig[REQ *pb.SetSingleChatSettingReq, RESP *pb.SetSingleChatSettingResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.SetSingleChatSettingReq, *pb.SetSingleChatSettingResp] {
	return wrapper.Config[*pb.SetSingleChatSettingReq, *pb.SetSingleChatSettingResp]{
		Do: func(ctx context.Context, in *pb.SetSingleChatSettingReq, opts ...grpc.CallOption) (*pb.SetSingleChatSettingResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().SetSingleChatSetting(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "SetSingleChatSetting", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.SetSingleChatSettingReq {
			return &pb.SetSingleChatSettingReq{}
		},
	}
}

// GetSingleChatSettingConfig ...
func GetSingleChatSettingConfig[REQ *pb.GetSingleChatSettingReq, RESP *pb.GetSingleChatSettingResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetSingleChatSettingReq, *pb.GetSingleChatSettingResp] {
	return wrapper.Config[*pb.GetSingleChatSettingReq, *pb.GetSingleChatSettingResp]{
		Do: func(ctx context.Context, in *pb.GetSingleChatSettingReq, opts ...grpc.CallOption) (*pb.GetSingleChatSettingResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().GetSingleChatSetting(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetSingleChatSetting", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetSingleChatSettingReq {
			return &pb.GetSingleChatSettingReq{}
		},
	}
}

// SetSingleMsgNotifyOptConfig ...
func SetSingleMsgNotifyOptConfig[REQ *pb.SetSingleMsgNotifyOptReq, RESP *pb.SetSingleMsgNotifyOptResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.SetSingleMsgNotifyOptReq, *pb.SetSingleMsgNotifyOptResp] {
	return wrapper.Config[*pb.SetSingleMsgNotifyOptReq, *pb.SetSingleMsgNotifyOptResp]{
		Do: func(ctx context.Context, in *pb.SetSingleMsgNotifyOptReq, opts ...grpc.CallOption) (*pb.SetSingleMsgNotifyOptResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().SetSingleMsgNotifyOpt(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "SetSingleMsgNotifyOpt", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.SetSingleMsgNotifyOptReq {
			return &pb.SetSingleMsgNotifyOptReq{}
		},
	}
}

// GetSingleMsgNotifyOptConfig ...
func GetSingleMsgNotifyOptConfig[REQ *pb.GetSingleMsgNotifyOptReq, RESP *pb.GetSingleMsgNotifyOptResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetSingleMsgNotifyOptReq, *pb.GetSingleMsgNotifyOptResp] {
	return wrapper.Config[*pb.GetSingleMsgNotifyOptReq, *pb.GetSingleMsgNotifyOptResp]{
		Do: func(ctx context.Context, in *pb.GetSingleMsgNotifyOptReq, opts ...grpc.CallOption) (*pb.GetSingleMsgNotifyOptResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().GetSingleMsgNotifyOpt(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetSingleMsgNotifyOpt", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetSingleMsgNotifyOptReq {
			return &pb.GetSingleMsgNotifyOptReq{}
		},
	}
}

// GetFriendListConfig ...
func GetFriendListConfig[REQ *pb.GetFriendListReq, RESP *pb.GetFriendListResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetFriendListReq, *pb.GetFriendListResp] {
	return wrapper.Config[*pb.GetFriendListReq, *pb.GetFriendListResp]{
		Do: func(ctx context.Context, in *pb.GetFriendListReq, opts ...grpc.CallOption) (*pb.GetFriendListResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().GetFriendList(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetFriendList", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetFriendListReq {
			return &pb.GetFriendListReq{}
		},
	}
}
