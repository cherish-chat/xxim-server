// Code generated by goctl. DO NOT EDIT.
// Source: user.peer.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/logic/accountservice"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"
)

type AccountServiceServer struct {
	svcCtx *svc.ServiceContext
	peerpb.UnimplementedAccountServiceServer
}

func NewAccountServiceServer(svcCtx *svc.ServiceContext) *AccountServiceServer {
	return &AccountServiceServer{
		svcCtx: svcCtx,
	}
}

// UserRegister 用户注册
func (s *AccountServiceServer) UserRegister(ctx context.Context, in *peerpb.UserRegisterReq) (*peerpb.UserRegisterResp, error) {
	l := accountservicelogic.NewUserRegisterLogic(ctx, s.svcCtx)
	return l.UserRegister(in)
}

// UserDestroy 用户注销
func (s *AccountServiceServer) UserDestroy(ctx context.Context, in *peerpb.UserDestroyReq) (*peerpb.UserDestroyResp, error) {
	l := accountservicelogic.NewUserDestroyLogic(ctx, s.svcCtx)
	return l.UserDestroy(in)
}

// UserToken 用户登录
func (s *AccountServiceServer) UserToken(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	l := accountservicelogic.NewUserTokenLogic(ctx, s.svcCtx)
	return l.UserToken(in)
}

// RefreshUserToken 刷新用户token
func (s *AccountServiceServer) RefreshUserToken(ctx context.Context, in *peerpb.RefreshUserTokenReq) (*peerpb.RefreshUserTokenResp, error) {
	l := accountservicelogic.NewRefreshUserTokenLogic(ctx, s.svcCtx)
	return l.RefreshUserToken(in)
}

// RevokeUserToken 注销用户token
func (s *AccountServiceServer) RevokeUserToken(ctx context.Context, in *peerpb.RevokeUserTokenReq) (*peerpb.RevokeUserTokenResp, error) {
	l := accountservicelogic.NewRevokeUserTokenLogic(ctx, s.svcCtx)
	return l.RevokeUserToken(in)
}

// UpdateUserAccountMap 更新用户账号信息
func (s *AccountServiceServer) UpdateUserAccountMap(ctx context.Context, in *peerpb.UpdateUserAccountMapReq) (*peerpb.UpdateUserAccountMapResp, error) {
	l := accountservicelogic.NewUpdateUserAccountMapLogic(ctx, s.svcCtx)
	return l.UpdateUserAccountMap(in)
}

// ResetUserAccountMap 重置用户账号信息
func (s *AccountServiceServer) ResetUserAccountMap(ctx context.Context, in *peerpb.ResetUserAccountMapReq) (*peerpb.ResetUserAccountMapResp, error) {
	l := accountservicelogic.NewResetUserAccountMapLogic(ctx, s.svcCtx)
	return l.ResetUserAccountMap(in)
}
