// Code generated by goctl. DO NOT EDIT!
// Source: appmgmt.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type AppMgmtServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedAppMgmtServiceServer
}

func NewAppMgmtServiceServer(svcCtx *svc.ServiceContext) *AppMgmtServiceServer {
	return &AppMgmtServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *AppMgmtServiceServer) GetAllAppMgmtConfig(ctx context.Context, in *pb.GetAllAppMgmtConfigReq) (*pb.GetAllAppMgmtConfigResp, error) {
	l := logic.NewGetAllAppMgmtConfigLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtConfig(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtConfig(ctx context.Context, in *pb.UpdateAppMgmtConfigReq) (*pb.UpdateAppMgmtConfigResp, error) {
	l := logic.NewUpdateAppMgmtConfigLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtConfig(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtVersion(ctx context.Context, in *pb.GetAllAppMgmtVersionReq) (*pb.GetAllAppMgmtVersionResp, error) {
	l := logic.NewGetAllAppMgmtVersionLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtVersion(in)
}

func (s *AppMgmtServiceServer) GetLatestVersion(ctx context.Context, in *pb.GetLatestVersionReq) (*pb.GetLatestVersionResp, error) {
	l := logic.NewGetLatestVersionLogic(ctx, s.svcCtx)
	return l.GetLatestVersion(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtVersionDetail(ctx context.Context, in *pb.GetAppMgmtVersionDetailReq) (*pb.GetAppMgmtVersionDetailResp, error) {
	l := logic.NewGetAppMgmtVersionDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtVersionDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtVersion(ctx context.Context, in *pb.AddAppMgmtVersionReq) (*pb.AddAppMgmtVersionResp, error) {
	l := logic.NewAddAppMgmtVersionLogic(ctx, s.svcCtx)
	return l.AddAppMgmtVersion(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtVersion(ctx context.Context, in *pb.UpdateAppMgmtVersionReq) (*pb.UpdateAppMgmtVersionResp, error) {
	l := logic.NewUpdateAppMgmtVersionLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtVersion(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtVersion(ctx context.Context, in *pb.DeleteAppMgmtVersionReq) (*pb.DeleteAppMgmtVersionResp, error) {
	l := logic.NewDeleteAppMgmtVersionLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtVersion(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtShieldWord(ctx context.Context, in *pb.GetAllAppMgmtShieldWordReq) (*pb.GetAllAppMgmtShieldWordResp, error) {
	l := logic.NewGetAllAppMgmtShieldWordLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtShieldWord(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtShieldWordDetail(ctx context.Context, in *pb.GetAppMgmtShieldWordDetailReq) (*pb.GetAppMgmtShieldWordDetailResp, error) {
	l := logic.NewGetAppMgmtShieldWordDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtShieldWordDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtShieldWord(ctx context.Context, in *pb.AddAppMgmtShieldWordReq) (*pb.AddAppMgmtShieldWordResp, error) {
	l := logic.NewAddAppMgmtShieldWordLogic(ctx, s.svcCtx)
	return l.AddAppMgmtShieldWord(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtShieldWord(ctx context.Context, in *pb.UpdateAppMgmtShieldWordReq) (*pb.UpdateAppMgmtShieldWordResp, error) {
	l := logic.NewUpdateAppMgmtShieldWordLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtShieldWord(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtShieldWord(ctx context.Context, in *pb.DeleteAppMgmtShieldWordReq) (*pb.DeleteAppMgmtShieldWordResp, error) {
	l := logic.NewDeleteAppMgmtShieldWordLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtShieldWord(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtVpn(ctx context.Context, in *pb.GetAllAppMgmtVpnReq) (*pb.GetAllAppMgmtVpnResp, error) {
	l := logic.NewGetAllAppMgmtVpnLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtVpn(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtVpnDetail(ctx context.Context, in *pb.GetAppMgmtVpnDetailReq) (*pb.GetAppMgmtVpnDetailResp, error) {
	l := logic.NewGetAppMgmtVpnDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtVpnDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtVpn(ctx context.Context, in *pb.AddAppMgmtVpnReq) (*pb.AddAppMgmtVpnResp, error) {
	l := logic.NewAddAppMgmtVpnLogic(ctx, s.svcCtx)
	return l.AddAppMgmtVpn(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtVpn(ctx context.Context, in *pb.UpdateAppMgmtVpnReq) (*pb.UpdateAppMgmtVpnResp, error) {
	l := logic.NewUpdateAppMgmtVpnLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtVpn(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtVpn(ctx context.Context, in *pb.DeleteAppMgmtVpnReq) (*pb.DeleteAppMgmtVpnResp, error) {
	l := logic.NewDeleteAppMgmtVpnLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtVpn(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtEmoji(ctx context.Context, in *pb.GetAllAppMgmtEmojiReq) (*pb.GetAllAppMgmtEmojiResp, error) {
	l := logic.NewGetAllAppMgmtEmojiLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtEmoji(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtEmojiDetail(ctx context.Context, in *pb.GetAppMgmtEmojiDetailReq) (*pb.GetAppMgmtEmojiDetailResp, error) {
	l := logic.NewGetAppMgmtEmojiDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtEmojiDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtEmoji(ctx context.Context, in *pb.AddAppMgmtEmojiReq) (*pb.AddAppMgmtEmojiResp, error) {
	l := logic.NewAddAppMgmtEmojiLogic(ctx, s.svcCtx)
	return l.AddAppMgmtEmoji(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtEmoji(ctx context.Context, in *pb.UpdateAppMgmtEmojiReq) (*pb.UpdateAppMgmtEmojiResp, error) {
	l := logic.NewUpdateAppMgmtEmojiLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtEmoji(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtEmoji(ctx context.Context, in *pb.DeleteAppMgmtEmojiReq) (*pb.DeleteAppMgmtEmojiResp, error) {
	l := logic.NewDeleteAppMgmtEmojiLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtEmoji(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtEmojiGroup(ctx context.Context, in *pb.GetAllAppMgmtEmojiGroupReq) (*pb.GetAllAppMgmtEmojiGroupResp, error) {
	l := logic.NewGetAllAppMgmtEmojiGroupLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtEmojiGroup(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtEmojiGroupDetail(ctx context.Context, in *pb.GetAppMgmtEmojiGroupDetailReq) (*pb.GetAppMgmtEmojiGroupDetailResp, error) {
	l := logic.NewGetAppMgmtEmojiGroupDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtEmojiGroupDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtEmojiGroup(ctx context.Context, in *pb.AddAppMgmtEmojiGroupReq) (*pb.AddAppMgmtEmojiGroupResp, error) {
	l := logic.NewAddAppMgmtEmojiGroupLogic(ctx, s.svcCtx)
	return l.AddAppMgmtEmojiGroup(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtEmojiGroup(ctx context.Context, in *pb.UpdateAppMgmtEmojiGroupReq) (*pb.UpdateAppMgmtEmojiGroupResp, error) {
	l := logic.NewUpdateAppMgmtEmojiGroupLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtEmojiGroup(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtEmojiGroup(ctx context.Context, in *pb.DeleteAppMgmtEmojiGroupReq) (*pb.DeleteAppMgmtEmojiGroupResp, error) {
	l := logic.NewDeleteAppMgmtEmojiGroupLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtEmojiGroup(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtNotice(ctx context.Context, in *pb.GetAllAppMgmtNoticeReq) (*pb.GetAllAppMgmtNoticeResp, error) {
	l := logic.NewGetAllAppMgmtNoticeLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtNotice(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtNoticeDetail(ctx context.Context, in *pb.GetAppMgmtNoticeDetailReq) (*pb.GetAppMgmtNoticeDetailResp, error) {
	l := logic.NewGetAppMgmtNoticeDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtNoticeDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtNotice(ctx context.Context, in *pb.AddAppMgmtNoticeReq) (*pb.AddAppMgmtNoticeResp, error) {
	l := logic.NewAddAppMgmtNoticeLogic(ctx, s.svcCtx)
	return l.AddAppMgmtNotice(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtNotice(ctx context.Context, in *pb.UpdateAppMgmtNoticeReq) (*pb.UpdateAppMgmtNoticeResp, error) {
	l := logic.NewUpdateAppMgmtNoticeLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtNotice(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtNotice(ctx context.Context, in *pb.DeleteAppMgmtNoticeReq) (*pb.DeleteAppMgmtNoticeResp, error) {
	l := logic.NewDeleteAppMgmtNoticeLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtNotice(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtLink(ctx context.Context, in *pb.GetAllAppMgmtLinkReq) (*pb.GetAllAppMgmtLinkResp, error) {
	l := logic.NewGetAllAppMgmtLinkLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtLink(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtLinkDetail(ctx context.Context, in *pb.GetAppMgmtLinkDetailReq) (*pb.GetAppMgmtLinkDetailResp, error) {
	l := logic.NewGetAppMgmtLinkDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtLinkDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtLink(ctx context.Context, in *pb.AddAppMgmtLinkReq) (*pb.AddAppMgmtLinkResp, error) {
	l := logic.NewAddAppMgmtLinkLogic(ctx, s.svcCtx)
	return l.AddAppMgmtLink(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtLink(ctx context.Context, in *pb.UpdateAppMgmtLinkReq) (*pb.UpdateAppMgmtLinkResp, error) {
	l := logic.NewUpdateAppMgmtLinkLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtLink(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtLink(ctx context.Context, in *pb.DeleteAppMgmtLinkReq) (*pb.DeleteAppMgmtLinkResp, error) {
	l := logic.NewDeleteAppMgmtLinkLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtLink(in)
}

func (s *AppMgmtServiceServer) AppGetAllConfig(ctx context.Context, in *pb.AppGetAllConfigReq) (*pb.AppGetAllConfigResp, error) {
	l := logic.NewAppGetAllConfigLogic(ctx, s.svcCtx)
	return l.AppGetAllConfig(in)
}

func (s *AppMgmtServiceServer) GetUploadInfo(ctx context.Context, in *pb.GetUploadInfoReq) (*pb.GetUploadInfoResp, error) {
	l := logic.NewGetUploadInfoLogic(ctx, s.svcCtx)
	return l.GetUploadInfo(in)
}

func (s *AppMgmtServiceServer) GetAllAppMgmtRichArticle(ctx context.Context, in *pb.GetAllAppMgmtRichArticleReq) (*pb.GetAllAppMgmtRichArticleResp, error) {
	l := logic.NewGetAllAppMgmtRichArticleLogic(ctx, s.svcCtx)
	return l.GetAllAppMgmtRichArticle(in)
}

func (s *AppMgmtServiceServer) GetAppMgmtRichArticleDetail(ctx context.Context, in *pb.GetAppMgmtRichArticleDetailReq) (*pb.GetAppMgmtRichArticleDetailResp, error) {
	l := logic.NewGetAppMgmtRichArticleDetailLogic(ctx, s.svcCtx)
	return l.GetAppMgmtRichArticleDetail(in)
}

func (s *AppMgmtServiceServer) AddAppMgmtRichArticle(ctx context.Context, in *pb.AddAppMgmtRichArticleReq) (*pb.AddAppMgmtRichArticleResp, error) {
	l := logic.NewAddAppMgmtRichArticleLogic(ctx, s.svcCtx)
	return l.AddAppMgmtRichArticle(in)
}

func (s *AppMgmtServiceServer) UpdateAppMgmtRichArticle(ctx context.Context, in *pb.UpdateAppMgmtRichArticleReq) (*pb.UpdateAppMgmtRichArticleResp, error) {
	l := logic.NewUpdateAppMgmtRichArticleLogic(ctx, s.svcCtx)
	return l.UpdateAppMgmtRichArticle(in)
}

func (s *AppMgmtServiceServer) DeleteAppMgmtRichArticle(ctx context.Context, in *pb.DeleteAppMgmtRichArticleReq) (*pb.DeleteAppMgmtRichArticleResp, error) {
	l := logic.NewDeleteAppMgmtRichArticleLogic(ctx, s.svcCtx)
	return l.DeleteAppMgmtRichArticle(in)
}

func (s *AppMgmtServiceServer) AppGetRichArticleList(ctx context.Context, in *pb.AppGetRichArticleListReq) (*pb.AppGetRichArticleListResp, error) {
	l := logic.NewAppGetRichArticleListLogic(ctx, s.svcCtx)
	return l.AppGetRichArticleList(in)
}
