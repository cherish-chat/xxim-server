// Code generated by goctl. DO NOT EDIT!
// Source: appmgmt.proto

package appmgmtservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddAppMgmtEmojiGroupReq         = pb.AddAppMgmtEmojiGroupReq
	AddAppMgmtEmojiGroupResp        = pb.AddAppMgmtEmojiGroupResp
	AddAppMgmtEmojiReq              = pb.AddAppMgmtEmojiReq
	AddAppMgmtEmojiResp             = pb.AddAppMgmtEmojiResp
	AddAppMgmtLinkReq               = pb.AddAppMgmtLinkReq
	AddAppMgmtLinkResp              = pb.AddAppMgmtLinkResp
	AddAppMgmtNoticeReq             = pb.AddAppMgmtNoticeReq
	AddAppMgmtNoticeResp            = pb.AddAppMgmtNoticeResp
	AddAppMgmtRichArticleReq        = pb.AddAppMgmtRichArticleReq
	AddAppMgmtRichArticleResp       = pb.AddAppMgmtRichArticleResp
	AddAppMgmtShieldWordReq         = pb.AddAppMgmtShieldWordReq
	AddAppMgmtShieldWordResp        = pb.AddAppMgmtShieldWordResp
	AddAppMgmtVersionReq            = pb.AddAppMgmtVersionReq
	AddAppMgmtVersionResp           = pb.AddAppMgmtVersionResp
	AddAppMgmtVpnReq                = pb.AddAppMgmtVpnReq
	AddAppMgmtVpnResp               = pb.AddAppMgmtVpnResp
	AppGetAllConfigReq              = pb.AppGetAllConfigReq
	AppGetAllConfigResp             = pb.AppGetAllConfigResp
	AppGetRichArticleListReq        = pb.AppGetRichArticleListReq
	AppGetRichArticleListResp       = pb.AppGetRichArticleListResp
	AppMgmtConfig                   = pb.AppMgmtConfig
	AppMgmtEmoji                    = pb.AppMgmtEmoji
	AppMgmtEmojiGroup               = pb.AppMgmtEmojiGroup
	AppMgmtLink                     = pb.AppMgmtLink
	AppMgmtNotice                   = pb.AppMgmtNotice
	AppMgmtRichArticle              = pb.AppMgmtRichArticle
	AppMgmtShieldWord               = pb.AppMgmtShieldWord
	AppMgmtVersion                  = pb.AppMgmtVersion
	AppMgmtVpn                      = pb.AppMgmtVpn
	DeleteAppMgmtEmojiGroupReq      = pb.DeleteAppMgmtEmojiGroupReq
	DeleteAppMgmtEmojiGroupResp     = pb.DeleteAppMgmtEmojiGroupResp
	DeleteAppMgmtEmojiReq           = pb.DeleteAppMgmtEmojiReq
	DeleteAppMgmtEmojiResp          = pb.DeleteAppMgmtEmojiResp
	DeleteAppMgmtLinkReq            = pb.DeleteAppMgmtLinkReq
	DeleteAppMgmtLinkResp           = pb.DeleteAppMgmtLinkResp
	DeleteAppMgmtNoticeReq          = pb.DeleteAppMgmtNoticeReq
	DeleteAppMgmtNoticeResp         = pb.DeleteAppMgmtNoticeResp
	DeleteAppMgmtRichArticleReq     = pb.DeleteAppMgmtRichArticleReq
	DeleteAppMgmtRichArticleResp    = pb.DeleteAppMgmtRichArticleResp
	DeleteAppMgmtShieldWordReq      = pb.DeleteAppMgmtShieldWordReq
	DeleteAppMgmtShieldWordResp     = pb.DeleteAppMgmtShieldWordResp
	DeleteAppMgmtVersionReq         = pb.DeleteAppMgmtVersionReq
	DeleteAppMgmtVersionResp        = pb.DeleteAppMgmtVersionResp
	DeleteAppMgmtVpnReq             = pb.DeleteAppMgmtVpnReq
	DeleteAppMgmtVpnResp            = pb.DeleteAppMgmtVpnResp
	GetAllAppMgmtConfigReq          = pb.GetAllAppMgmtConfigReq
	GetAllAppMgmtConfigResp         = pb.GetAllAppMgmtConfigResp
	GetAllAppMgmtEmojiGroupReq      = pb.GetAllAppMgmtEmojiGroupReq
	GetAllAppMgmtEmojiGroupResp     = pb.GetAllAppMgmtEmojiGroupResp
	GetAllAppMgmtEmojiReq           = pb.GetAllAppMgmtEmojiReq
	GetAllAppMgmtEmojiResp          = pb.GetAllAppMgmtEmojiResp
	GetAllAppMgmtLinkReq            = pb.GetAllAppMgmtLinkReq
	GetAllAppMgmtLinkResp           = pb.GetAllAppMgmtLinkResp
	GetAllAppMgmtNoticeReq          = pb.GetAllAppMgmtNoticeReq
	GetAllAppMgmtNoticeResp         = pb.GetAllAppMgmtNoticeResp
	GetAllAppMgmtRichArticleReq     = pb.GetAllAppMgmtRichArticleReq
	GetAllAppMgmtRichArticleResp    = pb.GetAllAppMgmtRichArticleResp
	GetAllAppMgmtShieldWordReq      = pb.GetAllAppMgmtShieldWordReq
	GetAllAppMgmtShieldWordResp     = pb.GetAllAppMgmtShieldWordResp
	GetAllAppMgmtVersionReq         = pb.GetAllAppMgmtVersionReq
	GetAllAppMgmtVersionResp        = pb.GetAllAppMgmtVersionResp
	GetAllAppMgmtVpnReq             = pb.GetAllAppMgmtVpnReq
	GetAllAppMgmtVpnResp            = pb.GetAllAppMgmtVpnResp
	GetAppAddressBookReq            = pb.GetAppAddressBookReq
	GetAppAddressBookResp           = pb.GetAppAddressBookResp
	GetAppAddressBookUrlReq         = pb.GetAppAddressBookUrlReq
	GetAppAddressBookUrlResp        = pb.GetAppAddressBookUrlResp
	GetAppMgmtEmojiDetailReq        = pb.GetAppMgmtEmojiDetailReq
	GetAppMgmtEmojiDetailResp       = pb.GetAppMgmtEmojiDetailResp
	GetAppMgmtEmojiGroupDetailReq   = pb.GetAppMgmtEmojiGroupDetailReq
	GetAppMgmtEmojiGroupDetailResp  = pb.GetAppMgmtEmojiGroupDetailResp
	GetAppMgmtLinkDetailReq         = pb.GetAppMgmtLinkDetailReq
	GetAppMgmtLinkDetailResp        = pb.GetAppMgmtLinkDetailResp
	GetAppMgmtNoticeDetailReq       = pb.GetAppMgmtNoticeDetailReq
	GetAppMgmtNoticeDetailResp      = pb.GetAppMgmtNoticeDetailResp
	GetAppMgmtRichArticleDetailReq  = pb.GetAppMgmtRichArticleDetailReq
	GetAppMgmtRichArticleDetailResp = pb.GetAppMgmtRichArticleDetailResp
	GetAppMgmtShieldWordDetailReq   = pb.GetAppMgmtShieldWordDetailReq
	GetAppMgmtShieldWordDetailResp  = pb.GetAppMgmtShieldWordDetailResp
	GetAppMgmtVersionDetailReq      = pb.GetAppMgmtVersionDetailReq
	GetAppMgmtVersionDetailResp     = pb.GetAppMgmtVersionDetailResp
	GetAppMgmtVpnDetailReq          = pb.GetAppMgmtVpnDetailReq
	GetAppMgmtVpnDetailResp         = pb.GetAppMgmtVpnDetailResp
	GetLatestVersionReq             = pb.GetLatestVersionReq
	GetLatestVersionResp            = pb.GetLatestVersionResp
	GetUploadInfoReq                = pb.GetUploadInfoReq
	GetUploadInfoResp               = pb.GetUploadInfoResp
	UpdateAppAddressBookReq         = pb.UpdateAppAddressBookReq
	UpdateAppAddressBookResp        = pb.UpdateAppAddressBookResp
	UpdateAppMgmtConfigReq          = pb.UpdateAppMgmtConfigReq
	UpdateAppMgmtConfigResp         = pb.UpdateAppMgmtConfigResp
	UpdateAppMgmtEmojiGroupReq      = pb.UpdateAppMgmtEmojiGroupReq
	UpdateAppMgmtEmojiGroupResp     = pb.UpdateAppMgmtEmojiGroupResp
	UpdateAppMgmtEmojiReq           = pb.UpdateAppMgmtEmojiReq
	UpdateAppMgmtEmojiResp          = pb.UpdateAppMgmtEmojiResp
	UpdateAppMgmtLinkReq            = pb.UpdateAppMgmtLinkReq
	UpdateAppMgmtLinkResp           = pb.UpdateAppMgmtLinkResp
	UpdateAppMgmtNoticeReq          = pb.UpdateAppMgmtNoticeReq
	UpdateAppMgmtNoticeResp         = pb.UpdateAppMgmtNoticeResp
	UpdateAppMgmtRichArticleReq     = pb.UpdateAppMgmtRichArticleReq
	UpdateAppMgmtRichArticleResp    = pb.UpdateAppMgmtRichArticleResp
	UpdateAppMgmtShieldWordReq      = pb.UpdateAppMgmtShieldWordReq
	UpdateAppMgmtShieldWordResp     = pb.UpdateAppMgmtShieldWordResp
	UpdateAppMgmtVersionReq         = pb.UpdateAppMgmtVersionReq
	UpdateAppMgmtVersionResp        = pb.UpdateAppMgmtVersionResp
	UpdateAppMgmtVpnReq             = pb.UpdateAppMgmtVpnReq
	UpdateAppMgmtVpnResp            = pb.UpdateAppMgmtVpnResp

	AppMgmtService interface {
		GetAllAppMgmtConfig(ctx context.Context, in *GetAllAppMgmtConfigReq, opts ...grpc.CallOption) (*GetAllAppMgmtConfigResp, error)
		UpdateAppMgmtConfig(ctx context.Context, in *UpdateAppMgmtConfigReq, opts ...grpc.CallOption) (*UpdateAppMgmtConfigResp, error)
		GetAllAppMgmtVersion(ctx context.Context, in *GetAllAppMgmtVersionReq, opts ...grpc.CallOption) (*GetAllAppMgmtVersionResp, error)
		GetLatestVersion(ctx context.Context, in *GetLatestVersionReq, opts ...grpc.CallOption) (*GetLatestVersionResp, error)
		GetAppMgmtVersionDetail(ctx context.Context, in *GetAppMgmtVersionDetailReq, opts ...grpc.CallOption) (*GetAppMgmtVersionDetailResp, error)
		AddAppMgmtVersion(ctx context.Context, in *AddAppMgmtVersionReq, opts ...grpc.CallOption) (*AddAppMgmtVersionResp, error)
		UpdateAppMgmtVersion(ctx context.Context, in *UpdateAppMgmtVersionReq, opts ...grpc.CallOption) (*UpdateAppMgmtVersionResp, error)
		DeleteAppMgmtVersion(ctx context.Context, in *DeleteAppMgmtVersionReq, opts ...grpc.CallOption) (*DeleteAppMgmtVersionResp, error)
		GetAllAppMgmtShieldWord(ctx context.Context, in *GetAllAppMgmtShieldWordReq, opts ...grpc.CallOption) (*GetAllAppMgmtShieldWordResp, error)
		GetAppMgmtShieldWordDetail(ctx context.Context, in *GetAppMgmtShieldWordDetailReq, opts ...grpc.CallOption) (*GetAppMgmtShieldWordDetailResp, error)
		AddAppMgmtShieldWord(ctx context.Context, in *AddAppMgmtShieldWordReq, opts ...grpc.CallOption) (*AddAppMgmtShieldWordResp, error)
		UpdateAppMgmtShieldWord(ctx context.Context, in *UpdateAppMgmtShieldWordReq, opts ...grpc.CallOption) (*UpdateAppMgmtShieldWordResp, error)
		DeleteAppMgmtShieldWord(ctx context.Context, in *DeleteAppMgmtShieldWordReq, opts ...grpc.CallOption) (*DeleteAppMgmtShieldWordResp, error)
		GetAllAppMgmtVpn(ctx context.Context, in *GetAllAppMgmtVpnReq, opts ...grpc.CallOption) (*GetAllAppMgmtVpnResp, error)
		GetAppMgmtVpnDetail(ctx context.Context, in *GetAppMgmtVpnDetailReq, opts ...grpc.CallOption) (*GetAppMgmtVpnDetailResp, error)
		AddAppMgmtVpn(ctx context.Context, in *AddAppMgmtVpnReq, opts ...grpc.CallOption) (*AddAppMgmtVpnResp, error)
		UpdateAppMgmtVpn(ctx context.Context, in *UpdateAppMgmtVpnReq, opts ...grpc.CallOption) (*UpdateAppMgmtVpnResp, error)
		DeleteAppMgmtVpn(ctx context.Context, in *DeleteAppMgmtVpnReq, opts ...grpc.CallOption) (*DeleteAppMgmtVpnResp, error)
		GetAllAppMgmtEmoji(ctx context.Context, in *GetAllAppMgmtEmojiReq, opts ...grpc.CallOption) (*GetAllAppMgmtEmojiResp, error)
		GetAppMgmtEmojiDetail(ctx context.Context, in *GetAppMgmtEmojiDetailReq, opts ...grpc.CallOption) (*GetAppMgmtEmojiDetailResp, error)
		AddAppMgmtEmoji(ctx context.Context, in *AddAppMgmtEmojiReq, opts ...grpc.CallOption) (*AddAppMgmtEmojiResp, error)
		UpdateAppMgmtEmoji(ctx context.Context, in *UpdateAppMgmtEmojiReq, opts ...grpc.CallOption) (*UpdateAppMgmtEmojiResp, error)
		DeleteAppMgmtEmoji(ctx context.Context, in *DeleteAppMgmtEmojiReq, opts ...grpc.CallOption) (*DeleteAppMgmtEmojiResp, error)
		GetAllAppMgmtEmojiGroup(ctx context.Context, in *GetAllAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*GetAllAppMgmtEmojiGroupResp, error)
		GetAppMgmtEmojiGroupDetail(ctx context.Context, in *GetAppMgmtEmojiGroupDetailReq, opts ...grpc.CallOption) (*GetAppMgmtEmojiGroupDetailResp, error)
		AddAppMgmtEmojiGroup(ctx context.Context, in *AddAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*AddAppMgmtEmojiGroupResp, error)
		UpdateAppMgmtEmojiGroup(ctx context.Context, in *UpdateAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*UpdateAppMgmtEmojiGroupResp, error)
		DeleteAppMgmtEmojiGroup(ctx context.Context, in *DeleteAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*DeleteAppMgmtEmojiGroupResp, error)
		GetAllAppMgmtNotice(ctx context.Context, in *GetAllAppMgmtNoticeReq, opts ...grpc.CallOption) (*GetAllAppMgmtNoticeResp, error)
		GetAppMgmtNoticeDetail(ctx context.Context, in *GetAppMgmtNoticeDetailReq, opts ...grpc.CallOption) (*GetAppMgmtNoticeDetailResp, error)
		AddAppMgmtNotice(ctx context.Context, in *AddAppMgmtNoticeReq, opts ...grpc.CallOption) (*AddAppMgmtNoticeResp, error)
		UpdateAppMgmtNotice(ctx context.Context, in *UpdateAppMgmtNoticeReq, opts ...grpc.CallOption) (*UpdateAppMgmtNoticeResp, error)
		DeleteAppMgmtNotice(ctx context.Context, in *DeleteAppMgmtNoticeReq, opts ...grpc.CallOption) (*DeleteAppMgmtNoticeResp, error)
		GetAllAppMgmtLink(ctx context.Context, in *GetAllAppMgmtLinkReq, opts ...grpc.CallOption) (*GetAllAppMgmtLinkResp, error)
		GetAppMgmtLinkDetail(ctx context.Context, in *GetAppMgmtLinkDetailReq, opts ...grpc.CallOption) (*GetAppMgmtLinkDetailResp, error)
		AddAppMgmtLink(ctx context.Context, in *AddAppMgmtLinkReq, opts ...grpc.CallOption) (*AddAppMgmtLinkResp, error)
		UpdateAppMgmtLink(ctx context.Context, in *UpdateAppMgmtLinkReq, opts ...grpc.CallOption) (*UpdateAppMgmtLinkResp, error)
		DeleteAppMgmtLink(ctx context.Context, in *DeleteAppMgmtLinkReq, opts ...grpc.CallOption) (*DeleteAppMgmtLinkResp, error)
		AppGetAllConfig(ctx context.Context, in *AppGetAllConfigReq, opts ...grpc.CallOption) (*AppGetAllConfigResp, error)
		GetUploadInfo(ctx context.Context, in *GetUploadInfoReq, opts ...grpc.CallOption) (*GetUploadInfoResp, error)
		GetAllAppMgmtRichArticle(ctx context.Context, in *GetAllAppMgmtRichArticleReq, opts ...grpc.CallOption) (*GetAllAppMgmtRichArticleResp, error)
		GetAppMgmtRichArticleDetail(ctx context.Context, in *GetAppMgmtRichArticleDetailReq, opts ...grpc.CallOption) (*GetAppMgmtRichArticleDetailResp, error)
		AddAppMgmtRichArticle(ctx context.Context, in *AddAppMgmtRichArticleReq, opts ...grpc.CallOption) (*AddAppMgmtRichArticleResp, error)
		UpdateAppMgmtRichArticle(ctx context.Context, in *UpdateAppMgmtRichArticleReq, opts ...grpc.CallOption) (*UpdateAppMgmtRichArticleResp, error)
		DeleteAppMgmtRichArticle(ctx context.Context, in *DeleteAppMgmtRichArticleReq, opts ...grpc.CallOption) (*DeleteAppMgmtRichArticleResp, error)
		AppGetRichArticleList(ctx context.Context, in *AppGetRichArticleListReq, opts ...grpc.CallOption) (*AppGetRichArticleListResp, error)
		UpdateAppAddressBook(ctx context.Context, in *UpdateAppAddressBookReq, opts ...grpc.CallOption) (*UpdateAppAddressBookResp, error)
		GetAppAddressBook(ctx context.Context, in *GetAppAddressBookReq, opts ...grpc.CallOption) (*GetAppAddressBookResp, error)
		GetAppAddressBookUrl(ctx context.Context, in *GetAppAddressBookUrlReq, opts ...grpc.CallOption) (*GetAppAddressBookUrlResp, error)
	}

	defaultAppMgmtService struct {
		cli zrpc.Client
	}
)

func NewAppMgmtService(cli zrpc.Client) AppMgmtService {
	return &defaultAppMgmtService{
		cli: cli,
	}
}

func (m *defaultAppMgmtService) GetAllAppMgmtConfig(ctx context.Context, in *GetAllAppMgmtConfigReq, opts ...grpc.CallOption) (*GetAllAppMgmtConfigResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtConfig(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtConfig(ctx context.Context, in *UpdateAppMgmtConfigReq, opts ...grpc.CallOption) (*UpdateAppMgmtConfigResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtConfig(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtVersion(ctx context.Context, in *GetAllAppMgmtVersionReq, opts ...grpc.CallOption) (*GetAllAppMgmtVersionResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtVersion(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetLatestVersion(ctx context.Context, in *GetLatestVersionReq, opts ...grpc.CallOption) (*GetLatestVersionResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetLatestVersion(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtVersionDetail(ctx context.Context, in *GetAppMgmtVersionDetailReq, opts ...grpc.CallOption) (*GetAppMgmtVersionDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtVersionDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtVersion(ctx context.Context, in *AddAppMgmtVersionReq, opts ...grpc.CallOption) (*AddAppMgmtVersionResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtVersion(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtVersion(ctx context.Context, in *UpdateAppMgmtVersionReq, opts ...grpc.CallOption) (*UpdateAppMgmtVersionResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtVersion(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtVersion(ctx context.Context, in *DeleteAppMgmtVersionReq, opts ...grpc.CallOption) (*DeleteAppMgmtVersionResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtVersion(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtShieldWord(ctx context.Context, in *GetAllAppMgmtShieldWordReq, opts ...grpc.CallOption) (*GetAllAppMgmtShieldWordResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtShieldWord(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtShieldWordDetail(ctx context.Context, in *GetAppMgmtShieldWordDetailReq, opts ...grpc.CallOption) (*GetAppMgmtShieldWordDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtShieldWordDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtShieldWord(ctx context.Context, in *AddAppMgmtShieldWordReq, opts ...grpc.CallOption) (*AddAppMgmtShieldWordResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtShieldWord(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtShieldWord(ctx context.Context, in *UpdateAppMgmtShieldWordReq, opts ...grpc.CallOption) (*UpdateAppMgmtShieldWordResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtShieldWord(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtShieldWord(ctx context.Context, in *DeleteAppMgmtShieldWordReq, opts ...grpc.CallOption) (*DeleteAppMgmtShieldWordResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtShieldWord(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtVpn(ctx context.Context, in *GetAllAppMgmtVpnReq, opts ...grpc.CallOption) (*GetAllAppMgmtVpnResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtVpn(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtVpnDetail(ctx context.Context, in *GetAppMgmtVpnDetailReq, opts ...grpc.CallOption) (*GetAppMgmtVpnDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtVpnDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtVpn(ctx context.Context, in *AddAppMgmtVpnReq, opts ...grpc.CallOption) (*AddAppMgmtVpnResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtVpn(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtVpn(ctx context.Context, in *UpdateAppMgmtVpnReq, opts ...grpc.CallOption) (*UpdateAppMgmtVpnResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtVpn(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtVpn(ctx context.Context, in *DeleteAppMgmtVpnReq, opts ...grpc.CallOption) (*DeleteAppMgmtVpnResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtVpn(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtEmoji(ctx context.Context, in *GetAllAppMgmtEmojiReq, opts ...grpc.CallOption) (*GetAllAppMgmtEmojiResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtEmoji(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtEmojiDetail(ctx context.Context, in *GetAppMgmtEmojiDetailReq, opts ...grpc.CallOption) (*GetAppMgmtEmojiDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtEmojiDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtEmoji(ctx context.Context, in *AddAppMgmtEmojiReq, opts ...grpc.CallOption) (*AddAppMgmtEmojiResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtEmoji(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtEmoji(ctx context.Context, in *UpdateAppMgmtEmojiReq, opts ...grpc.CallOption) (*UpdateAppMgmtEmojiResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtEmoji(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtEmoji(ctx context.Context, in *DeleteAppMgmtEmojiReq, opts ...grpc.CallOption) (*DeleteAppMgmtEmojiResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtEmoji(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtEmojiGroup(ctx context.Context, in *GetAllAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*GetAllAppMgmtEmojiGroupResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtEmojiGroup(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtEmojiGroupDetail(ctx context.Context, in *GetAppMgmtEmojiGroupDetailReq, opts ...grpc.CallOption) (*GetAppMgmtEmojiGroupDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtEmojiGroupDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtEmojiGroup(ctx context.Context, in *AddAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*AddAppMgmtEmojiGroupResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtEmojiGroup(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtEmojiGroup(ctx context.Context, in *UpdateAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*UpdateAppMgmtEmojiGroupResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtEmojiGroup(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtEmojiGroup(ctx context.Context, in *DeleteAppMgmtEmojiGroupReq, opts ...grpc.CallOption) (*DeleteAppMgmtEmojiGroupResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtEmojiGroup(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtNotice(ctx context.Context, in *GetAllAppMgmtNoticeReq, opts ...grpc.CallOption) (*GetAllAppMgmtNoticeResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtNotice(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtNoticeDetail(ctx context.Context, in *GetAppMgmtNoticeDetailReq, opts ...grpc.CallOption) (*GetAppMgmtNoticeDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtNoticeDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtNotice(ctx context.Context, in *AddAppMgmtNoticeReq, opts ...grpc.CallOption) (*AddAppMgmtNoticeResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtNotice(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtNotice(ctx context.Context, in *UpdateAppMgmtNoticeReq, opts ...grpc.CallOption) (*UpdateAppMgmtNoticeResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtNotice(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtNotice(ctx context.Context, in *DeleteAppMgmtNoticeReq, opts ...grpc.CallOption) (*DeleteAppMgmtNoticeResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtNotice(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtLink(ctx context.Context, in *GetAllAppMgmtLinkReq, opts ...grpc.CallOption) (*GetAllAppMgmtLinkResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtLink(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtLinkDetail(ctx context.Context, in *GetAppMgmtLinkDetailReq, opts ...grpc.CallOption) (*GetAppMgmtLinkDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtLinkDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtLink(ctx context.Context, in *AddAppMgmtLinkReq, opts ...grpc.CallOption) (*AddAppMgmtLinkResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtLink(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtLink(ctx context.Context, in *UpdateAppMgmtLinkReq, opts ...grpc.CallOption) (*UpdateAppMgmtLinkResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtLink(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtLink(ctx context.Context, in *DeleteAppMgmtLinkReq, opts ...grpc.CallOption) (*DeleteAppMgmtLinkResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtLink(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AppGetAllConfig(ctx context.Context, in *AppGetAllConfigReq, opts ...grpc.CallOption) (*AppGetAllConfigResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AppGetAllConfig(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetUploadInfo(ctx context.Context, in *GetUploadInfoReq, opts ...grpc.CallOption) (*GetUploadInfoResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetUploadInfo(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAllAppMgmtRichArticle(ctx context.Context, in *GetAllAppMgmtRichArticleReq, opts ...grpc.CallOption) (*GetAllAppMgmtRichArticleResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAllAppMgmtRichArticle(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppMgmtRichArticleDetail(ctx context.Context, in *GetAppMgmtRichArticleDetailReq, opts ...grpc.CallOption) (*GetAppMgmtRichArticleDetailResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppMgmtRichArticleDetail(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AddAppMgmtRichArticle(ctx context.Context, in *AddAppMgmtRichArticleReq, opts ...grpc.CallOption) (*AddAppMgmtRichArticleResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AddAppMgmtRichArticle(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppMgmtRichArticle(ctx context.Context, in *UpdateAppMgmtRichArticleReq, opts ...grpc.CallOption) (*UpdateAppMgmtRichArticleResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppMgmtRichArticle(ctx, in, opts...)
}

func (m *defaultAppMgmtService) DeleteAppMgmtRichArticle(ctx context.Context, in *DeleteAppMgmtRichArticleReq, opts ...grpc.CallOption) (*DeleteAppMgmtRichArticleResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.DeleteAppMgmtRichArticle(ctx, in, opts...)
}

func (m *defaultAppMgmtService) AppGetRichArticleList(ctx context.Context, in *AppGetRichArticleListReq, opts ...grpc.CallOption) (*AppGetRichArticleListResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.AppGetRichArticleList(ctx, in, opts...)
}

func (m *defaultAppMgmtService) UpdateAppAddressBook(ctx context.Context, in *UpdateAppAddressBookReq, opts ...grpc.CallOption) (*UpdateAppAddressBookResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.UpdateAppAddressBook(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppAddressBook(ctx context.Context, in *GetAppAddressBookReq, opts ...grpc.CallOption) (*GetAppAddressBookResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppAddressBook(ctx, in, opts...)
}

func (m *defaultAppMgmtService) GetAppAddressBookUrl(ctx context.Context, in *GetAppAddressBookUrlReq, opts ...grpc.CallOption) (*GetAppAddressBookUrlResp, error) {
	client := pb.NewAppMgmtServiceClient(m.cli.Conn())
	return client.GetAppAddressBookUrl(ctx, in, opts...)
}
