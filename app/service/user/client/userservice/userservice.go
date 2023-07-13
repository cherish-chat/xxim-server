// Code generated by goctl. DO NOT EDIT.
// Source: user.peer.proto

package userservice

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateRobotReq           = peerpb.CreateRobotReq
	CreateRobotResp          = peerpb.CreateRobotResp
	GetSelfUserInfoReq       = peerpb.GetSelfUserInfoReq
	GetSelfUserInfoResp      = peerpb.GetSelfUserInfoResp
	GetUserInfoReq           = peerpb.GetUserInfoReq
	GetUserInfoResp          = peerpb.GetUserInfoResp
	GetUserModelByIdReq      = peerpb.GetUserModelByIdReq
	GetUserModelByIdReq_Opt  = peerpb.GetUserModelByIdReq_Opt
	GetUserModelByIdResp     = peerpb.GetUserModelByIdResp
	GetUserModelByIdsReq     = peerpb.GetUserModelByIdsReq
	GetUserModelByIdsReq_Opt = peerpb.GetUserModelByIdsReq_Opt
	GetUserModelByIdsResp    = peerpb.GetUserModelByIdsResp
	RefreshUserTokenReq      = peerpb.RefreshUserTokenReq
	RefreshUserTokenResp     = peerpb.RefreshUserTokenResp
	ResetUserAccountMapReq   = peerpb.ResetUserAccountMapReq
	ResetUserAccountMapResp  = peerpb.ResetUserAccountMapResp
	RevokeUserTokenReq       = peerpb.RevokeUserTokenReq
	RevokeUserTokenResp      = peerpb.RevokeUserTokenResp
	UpdateUserAccountMapReq  = peerpb.UpdateUserAccountMapReq
	UpdateUserAccountMapResp = peerpb.UpdateUserAccountMapResp
	UpdateUserCountMapReq    = peerpb.UpdateUserCountMapReq
	UpdateUserCountMapResp   = peerpb.UpdateUserCountMapResp
	UpdateUserExtraMapReq    = peerpb.UpdateUserExtraMapReq
	UpdateUserExtraMapResp   = peerpb.UpdateUserExtraMapResp
	UpdateUserProfileMapReq  = peerpb.UpdateUserProfileMapReq
	UpdateUserProfileMapResp = peerpb.UpdateUserProfileMapResp
	UserAfterKeepAliveReq    = peerpb.UserAfterKeepAliveReq
	UserAfterKeepAliveResp   = peerpb.UserAfterKeepAliveResp
	UserAfterOfflineReq      = peerpb.UserAfterOfflineReq
	UserAfterOfflineResp     = peerpb.UserAfterOfflineResp
	UserAfterOnlineReq       = peerpb.UserAfterOnlineReq
	UserAfterOnlineResp      = peerpb.UserAfterOnlineResp
	UserBeforeConnectReq     = peerpb.UserBeforeConnectReq
	UserBeforeConnectResp    = peerpb.UserBeforeConnectResp
	UserBeforeRequestReq     = peerpb.UserBeforeRequestReq
	UserBeforeRequestResp    = peerpb.UserBeforeRequestResp
	UserDestroyReq           = peerpb.UserDestroyReq
	UserDestroyResp          = peerpb.UserDestroyResp
	UserRegisterReq          = peerpb.UserRegisterReq
	UserRegisterResp         = peerpb.UserRegisterResp
	UserTokenReq             = peerpb.UserTokenReq
	UserTokenResp            = peerpb.UserTokenResp

	UserService interface {
		// UpdateUserProfileMap 更新用户个人信息
		UpdateUserProfileMap(ctx context.Context, in *UpdateUserProfileMapReq, opts ...grpc.CallOption) (*UpdateUserProfileMapResp, error)
		// UpdateUserExtraMap 更新用户扩展信息
		UpdateUserExtraMap(ctx context.Context, in *UpdateUserExtraMapReq, opts ...grpc.CallOption) (*UpdateUserExtraMapResp, error)
		// UpdateUserCountMap 更新用户计数信息
		UpdateUserCountMap(ctx context.Context, in *UpdateUserCountMapReq, opts ...grpc.CallOption) (*UpdateUserCountMapResp, error)
		// GetSelfUserInfo 获取自己的用户信息
		GetSelfUserInfo(ctx context.Context, in *GetSelfUserInfoReq, opts ...grpc.CallOption) (*GetSelfUserInfoResp, error)
		// GetUserInfo 获取用户信息
		GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
		// GetUserModelById 获取用户模型
		GetUserModelById(ctx context.Context, in *GetUserModelByIdReq, opts ...grpc.CallOption) (*GetUserModelByIdResp, error)
		// GetUserModelByIds 批量获取用户模型
		GetUserModelByIds(ctx context.Context, in *GetUserModelByIdsReq, opts ...grpc.CallOption) (*GetUserModelByIdsResp, error)
	}

	defaultUserService struct {
		cli zrpc.Client
	}
)

func NewUserService(cli zrpc.Client) UserService {
	return &defaultUserService{
		cli: cli,
	}
}

// UpdateUserProfileMap 更新用户个人信息
func (m *defaultUserService) UpdateUserProfileMap(ctx context.Context, in *UpdateUserProfileMapReq, opts ...grpc.CallOption) (*UpdateUserProfileMapResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.UpdateUserProfileMap(ctx, in, opts...)
}

// UpdateUserExtraMap 更新用户扩展信息
func (m *defaultUserService) UpdateUserExtraMap(ctx context.Context, in *UpdateUserExtraMapReq, opts ...grpc.CallOption) (*UpdateUserExtraMapResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.UpdateUserExtraMap(ctx, in, opts...)
}

// UpdateUserCountMap 更新用户计数信息
func (m *defaultUserService) UpdateUserCountMap(ctx context.Context, in *UpdateUserCountMapReq, opts ...grpc.CallOption) (*UpdateUserCountMapResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.UpdateUserCountMap(ctx, in, opts...)
}

// GetSelfUserInfo 获取自己的用户信息
func (m *defaultUserService) GetSelfUserInfo(ctx context.Context, in *GetSelfUserInfoReq, opts ...grpc.CallOption) (*GetSelfUserInfoResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.GetSelfUserInfo(ctx, in, opts...)
}

// GetUserInfo 获取用户信息
func (m *defaultUserService) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

// GetUserModelById 获取用户模型
func (m *defaultUserService) GetUserModelById(ctx context.Context, in *GetUserModelByIdReq, opts ...grpc.CallOption) (*GetUserModelByIdResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.GetUserModelById(ctx, in, opts...)
}

// GetUserModelByIds 批量获取用户模型
func (m *defaultUserService) GetUserModelByIds(ctx context.Context, in *GetUserModelByIdsReq, opts ...grpc.CallOption) (*GetUserModelByIdsResp, error) {
	client := peerpb.NewUserServiceClient(m.cli.Conn())
	return client.GetUserModelByIds(ctx, in, opts...)
}
