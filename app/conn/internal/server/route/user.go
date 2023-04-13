package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterUser(svcCtx *svc.ServiceContext) {
	// user
	{
		// SearchUsersByKeywordReq SearchUsersByKeywordResp
		{
			route := conngateway.Route[*pb.SearchUsersByKeywordReq, *pb.SearchUsersByKeywordResp]{
				NewRequest: func() *pb.SearchUsersByKeywordReq {
					return &pb.SearchUsersByKeywordReq{}
				},
				Do: svcCtx.UserService().SearchUsersByKeyword,
			}
			conngateway.AddRoute("/v1/user/searchUsersByKeyword", route)
		}
		// GetUserSettingsReq GetUserSettingsResp
		{
			route := conngateway.Route[*pb.GetUserSettingsReq, *pb.GetUserSettingsResp]{
				NewRequest: func() *pb.GetUserSettingsReq {
					return &pb.GetUserSettingsReq{}
				},
				Do: svcCtx.UserService().GetUserSettings,
			}
			conngateway.AddRoute("/v1/user/getUserSettings", route)
		}
		// SetUserSettingsReq SetUserSettingsResp
		{
			route := conngateway.Route[*pb.SetUserSettingsReq, *pb.SetUserSettingsResp]{
				NewRequest: func() *pb.SetUserSettingsReq {
					return &pb.SetUserSettingsReq{}
				},
				Do: svcCtx.UserService().SetUserSettings,
			}
			conngateway.AddRoute("/v1/user/setUserSettings", route)
		}
		// UpdateUserInfoReq UpdateUserInfoResp
		{
			route := conngateway.Route[*pb.UpdateUserInfoReq, *pb.UpdateUserInfoResp]{
				NewRequest: func() *pb.UpdateUserInfoReq {
					return &pb.UpdateUserInfoReq{}
				},
				Do: svcCtx.UserService().UpdateUserInfo,
			}
			conngateway.AddRoute("/v1/user/updateUserInfo", route)
		}
		// UpdateUserPasswordReq UpdateUserPasswordResp
		{
			route := conngateway.Route[*pb.UpdateUserPasswordReq, *pb.UpdateUserPasswordResp]{
				NewRequest: func() *pb.UpdateUserPasswordReq {
					return &pb.UpdateUserPasswordReq{}
				},
				Do: svcCtx.UserService().UpdateUserPassword,
			}
			conngateway.AddRoute("/v1/user/updateUserPassword", route)
		}
		// GetUserHomeReq GetUserHomeResp
		{
			route := conngateway.Route[*pb.GetUserHomeReq, *pb.GetUserHomeResp]{
				NewRequest: func() *pb.GetUserHomeReq {
					return &pb.GetUserHomeReq{}
				},
				Do: svcCtx.UserService().GetUserHome,
			}
			conngateway.AddRoute("/v1/user/getUserHome", route)
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
				Do: svcCtx.UserService().Login,
			}
			conngateway.AddRoute("/v1/user/white/login", route)
		}
		// 确认注册
		{
			route := conngateway.Route[*pb.ConfirmRegisterReq, *pb.ConfirmRegisterResp]{
				NewRequest: func() *pb.ConfirmRegisterReq {
					return &pb.ConfirmRegisterReq{}
				},
				Do: svcCtx.UserService().ConfirmRegister,
			}
			conngateway.AddRoute("/v1/user/white/confirmRegister", route)
		}
		// 注册
		{
			route := conngateway.Route[*pb.RegisterReq, *pb.RegisterResp]{
				NewRequest: func() *pb.RegisterReq {
					return &pb.RegisterReq{}
				},
				Do: svcCtx.UserService().Register,
			}
			conngateway.AddRoute("/v1/user/white/register", route)
		}
		// SendSmsReq SendSmsResp
		{
			route := conngateway.Route[*pb.SendSmsReq, *pb.SendSmsResp]{
				NewRequest: func() *pb.SendSmsReq {
					return &pb.SendSmsReq{}
				},
				Do: svcCtx.UserService().SendSms,
			}
			conngateway.AddRoute("/v1/user/white/sendSms", route)
		}
		// VerifySmsReq VerifySmsResp
		{
			route := conngateway.Route[*pb.VerifySmsReq, *pb.VerifySmsResp]{
				NewRequest: func() *pb.VerifySmsReq {
					return &pb.VerifySmsReq{}
				},
				Do: svcCtx.UserService().VerifySms,
			}
			conngateway.AddRoute("/v1/user/white/verifySms", route)
		}
		// GetCaptchaCodeReq GetCaptchaCodeResp
		{
			route := conngateway.Route[*pb.GetCaptchaCodeReq, *pb.GetCaptchaCodeResp]{
				NewRequest: func() *pb.GetCaptchaCodeReq {
					return &pb.GetCaptchaCodeReq{}
				},
				Do: svcCtx.UserService().GetCaptchaCode,
			}
			conngateway.AddRoute("/v1/user/white/getCaptchaCode", route)
		}
		// VerifyCaptchaCodeReq VerifyCaptchaCodeResp
		{
			route := conngateway.Route[*pb.VerifyCaptchaCodeReq, *pb.VerifyCaptchaCodeResp]{
				NewRequest: func() *pb.VerifyCaptchaCodeReq {
					return &pb.VerifyCaptchaCodeReq{}
				},
				Do: svcCtx.UserService().VerifyCaptchaCode,
			}
			conngateway.AddRoute("/v1/user/white/verifyCaptchaCode", route)
		}
		// ResetPasswordReq ResetPasswordResp
		{
			route := conngateway.Route[*pb.ResetPasswordReq, *pb.ResetPasswordResp]{
				NewRequest: func() *pb.ResetPasswordReq {
					return &pb.ResetPasswordReq{}
				},
				Do: svcCtx.UserService().ResetPassword,
			}
			conngateway.AddRoute("/v1/user/white/resetPassword", route)
		}
		// ReportUserReq ReportUserResp
		{
			route := conngateway.Route[*pb.ReportUserReq, *pb.ReportUserResp]{
				NewRequest: func() *pb.ReportUserReq {
					return &pb.ReportUserReq{}
				},
				Do: svcCtx.UserService().ReportUser,
			}
			conngateway.AddRoute("/v1/user/reportUser", route)
		}
		// DestroyAccountReq DestroyAccountResp
		{
			route := conngateway.Route[*pb.DestroyAccountReq, *pb.DestroyAccountResp]{
				NewRequest: func() *pb.DestroyAccountReq {
					return &pb.DestroyAccountReq{}
				},
				Do: svcCtx.UserService().DestroyAccount,
			}
			conngateway.AddRoute("/v1/user/destroyAccount", route)
		}
		// RecoverAccountReq RecoverAccountResp
		{
			route := conngateway.Route[*pb.RecoverAccountReq, *pb.RecoverAccountResp]{
				NewRequest: func() *pb.RecoverAccountReq {
					return &pb.RecoverAccountReq{}
				},
				Do: svcCtx.UserService().RecoverAccount,
			}
			conngateway.AddRoute("/v1/user/white/recoverAccount", route)
		}
	}
}
