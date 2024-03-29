// Code generated by goctl. DO NOT EDIT.
// Source: third.proto

package emailservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CaptchaVerifyReq    = pb.CaptchaVerifyReq
	CaptchaVerifyResp   = pb.CaptchaVerifyResp
	EmailCodeSendReq    = pb.EmailCodeSendReq
	EmailCodeSendResp   = pb.EmailCodeSendResp
	EmailCodeVerifyReq  = pb.EmailCodeVerifyReq
	EmailCodeVerifyResp = pb.EmailCodeVerifyResp
	GetCaptchaReq       = pb.GetCaptchaReq
	GetCaptchaResp      = pb.GetCaptchaResp
	SmsCodeSendReq      = pb.SmsCodeSendReq
	SmsCodeSendResp     = pb.SmsCodeSendResp
	SmsCodeVerifyReq    = pb.SmsCodeVerifyReq
	SmsCodeVerifyResp   = pb.SmsCodeVerifyResp

	EmailService interface {
		// EmailCodeSend 发送邮件
		EmailCodeSend(ctx context.Context, in *EmailCodeSendReq, opts ...grpc.CallOption) (*EmailCodeSendResp, error)
		// EmailCodeVerify 验证邮件
		EmailCodeVerify(ctx context.Context, in *EmailCodeVerifyReq, opts ...grpc.CallOption) (*EmailCodeVerifyResp, error)
	}

	defaultEmailService struct {
		cli zrpc.Client
	}
)

func NewEmailService(cli zrpc.Client) EmailService {
	return &defaultEmailService{
		cli: cli,
	}
}

// EmailCodeSend 发送邮件
func (m *defaultEmailService) EmailCodeSend(ctx context.Context, in *EmailCodeSendReq, opts ...grpc.CallOption) (*EmailCodeSendResp, error) {
	client := pb.NewEmailServiceClient(m.cli.Conn())
	return client.EmailCodeSend(ctx, in, opts...)
}

// EmailCodeVerify 验证邮件
func (m *defaultEmailService) EmailCodeVerify(ctx context.Context, in *EmailCodeVerifyReq, opts ...grpc.CallOption) (*EmailCodeVerifyResp, error) {
	client := pb.NewEmailServiceClient(m.cli.Conn())
	return client.EmailCodeVerify(ctx, in, opts...)
}
