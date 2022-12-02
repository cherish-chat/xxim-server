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

// SetSingleConvSettingConfig ...
func SetSingleConvSettingConfig[REQ *pb.SetSingleConvSettingReq, RESP *pb.SetSingleConvSettingResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.SetSingleConvSettingReq, *pb.SetSingleConvSettingResp] {
	return wrapper.Config[*pb.SetSingleConvSettingReq, *pb.SetSingleConvSettingResp]{
		Do: func(ctx context.Context, in *pb.SetSingleConvSettingReq, opts ...grpc.CallOption) (*pb.SetSingleConvSettingResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().SetSingleConvSetting(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "SetSingleConvSetting", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.SetSingleConvSettingReq {
			return &pb.SetSingleConvSettingReq{}
		},
	}
}

// GetSingleConvSettingConfig ...
func GetSingleConvSettingConfig[REQ *pb.GetSingleConvSettingReq, RESP *pb.GetSingleConvSettingResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetSingleConvSettingReq, *pb.GetSingleConvSettingResp] {
	return wrapper.Config[*pb.GetSingleConvSettingReq, *pb.GetSingleConvSettingResp]{
		Do: func(ctx context.Context, in *pb.GetSingleConvSettingReq, opts ...grpc.CallOption) (*pb.GetSingleConvSettingResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().GetSingleConvSetting(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetSingleConvSetting", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetSingleConvSettingReq {
			return &pb.GetSingleConvSettingReq{}
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

// GetMyFriendEventListConfig ...
func GetMyFriendEventListConfig[REQ *pb.GetMyFriendEventListReq, RESP *pb.GetMyFriendEventListResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetMyFriendEventListReq, *pb.GetMyFriendEventListResp] {
	return wrapper.Config[*pb.GetMyFriendEventListReq, *pb.GetMyFriendEventListResp]{
		Do: func(ctx context.Context, in *pb.GetMyFriendEventListReq, opts ...grpc.CallOption) (*pb.GetMyFriendEventListResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.RelationService().GetMyFriendEventList(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetMyFriendEventList", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetMyFriendEventListReq {
			return &pb.GetMyFriendEventListReq{}
		},
	}
}
