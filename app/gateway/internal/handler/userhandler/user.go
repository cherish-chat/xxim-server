package userhandler

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

func LoginConfig[REQ *pb.LoginReq, RESP *pb.LoginResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.LoginReq, *pb.LoginResp] {
	return wrapper.Config[*pb.LoginReq, *pb.LoginResp]{
		Do: func(ctx context.Context, in *pb.LoginReq, opts ...grpc.CallOption) (*pb.LoginResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().Login(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "Login", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.LoginReq {
			return &pb.LoginReq{}
		},
	}
}

func ConfirmRegisterConfig[REQ *pb.ConfirmRegisterReq, RESP *pb.ConfirmRegisterResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.ConfirmRegisterReq, *pb.ConfirmRegisterResp] {
	return wrapper.Config[*pb.ConfirmRegisterReq, *pb.ConfirmRegisterResp]{
		Do: func(ctx context.Context, in *pb.ConfirmRegisterReq, opts ...grpc.CallOption) (*pb.ConfirmRegisterResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().ConfirmRegister(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "ConfirmRegister", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.ConfirmRegisterReq {
			return &pb.ConfirmRegisterReq{}
		},
	}
}

func SearchUsersByKeywordConfig[REQ *pb.SearchUsersByKeywordReq, RESP *pb.SearchUsersByKeywordResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.SearchUsersByKeywordReq, *pb.SearchUsersByKeywordResp] {
	return wrapper.Config[*pb.SearchUsersByKeywordReq, *pb.SearchUsersByKeywordResp]{
		Do: func(ctx context.Context, in *pb.SearchUsersByKeywordReq, opts ...grpc.CallOption) (*pb.SearchUsersByKeywordResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().SearchUsersByKeyword(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "SearchUsersByKeyword", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.SearchUsersByKeywordReq {
			return &pb.SearchUsersByKeywordReq{}
		},
	}
}

func GetUserHomeConfig[REQ *pb.GetUserHomeReq, RESP *pb.GetUserHomeResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetUserHomeReq, *pb.GetUserHomeResp] {
	return wrapper.Config[*pb.GetUserHomeReq, *pb.GetUserHomeResp]{
		Do: func(ctx context.Context, in *pb.GetUserHomeReq, opts ...grpc.CallOption) (*pb.GetUserHomeResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().GetUserHome(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "GetUserHome", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetUserHomeReq {
			return &pb.GetUserHomeReq{}
		},
	}
}

func GetUserSettingsConfig[REQ *pb.GetUserSettingsReq, RESP *pb.GetUserSettingsResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetUserSettingsReq, *pb.GetUserSettingsResp] {
	return wrapper.Config[*pb.GetUserSettingsReq, *pb.GetUserSettingsResp]{
		Do: func(ctx context.Context, in *pb.GetUserSettingsReq, opts ...grpc.CallOption) (*pb.GetUserSettingsResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().GetUserSettings(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "GetUserSettings", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetUserSettingsReq {
			return &pb.GetUserSettingsReq{}
		},
	}
}

func SetUserSettingsConfig[REQ *pb.SetUserSettingsReq, RESP *pb.SetUserSettingsResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.SetUserSettingsReq, *pb.SetUserSettingsResp] {
	return wrapper.Config[*pb.SetUserSettingsReq, *pb.SetUserSettingsResp]{
		Do: func(ctx context.Context, in *pb.SetUserSettingsReq, opts ...grpc.CallOption) (*pb.SetUserSettingsResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.UserService().SetUserSettings(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "SetUserSettings", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.SetUserSettingsReq {
			return &pb.SetUserSettingsReq{}
		},
	}
}
