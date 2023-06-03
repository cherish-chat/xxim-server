package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterAppMgmt(svcCtx *svc.ServiceContext) {
	// AppGetAllConfigReq AppGetAllConfigResp
	{
		route := conngateway.Route[*pb.AppGetAllConfigReq, *pb.AppGetAllConfigResp]{
			NewRequest: func() *pb.AppGetAllConfigReq {
				return &pb.AppGetAllConfigReq{}
			},
			Do: svcCtx.AppMgmtService().AppGetAllConfig,
		}
		conngateway.AddRoute("/v1/appmgmt/white/appGetAllConfig", route)
	}
	// GetLatestVersionReq GetLatestVersionResp
	{
		route := conngateway.Route[*pb.GetLatestVersionReq, *pb.GetLatestVersionResp]{
			NewRequest: func() *pb.GetLatestVersionReq {
				return &pb.GetLatestVersionReq{}
			},
			Do: svcCtx.AppMgmtService().GetLatestVersion,
		}
		conngateway.AddRoute("/v1/appmgmt/white/getLatestVersion", route)
	}
	// GetUploadInfoReq GetUploadInfoResp
	{
		route := conngateway.Route[*pb.GetUploadInfoReq, *pb.GetUploadInfoResp]{
			NewRequest: func() *pb.GetUploadInfoReq {
				return &pb.GetUploadInfoReq{}
			},
			Do: svcCtx.AppMgmtService().GetUploadInfo,
		}
		conngateway.AddRoute("/v1/appmgmt/getUploadInfo", route)
	}
	// GetAllAppMgmtLinkReq GetAllAppMgmtLinkResp
	{
		route := conngateway.Route[*pb.GetAllAppMgmtLinkReq, *pb.GetAllAppMgmtLinkResp]{
			NewRequest: func() *pb.GetAllAppMgmtLinkReq {
				return &pb.GetAllAppMgmtLinkReq{}
			},
			Do: svcCtx.AppMgmtService().GetAllAppMgmtLink,
		}
		conngateway.AddRoute("/v1/appmgmt/getAllAppMgmtLink", route)
	}
	// AppGetRichArticleListReq AppGetRichArticleListResp
	{
		route := conngateway.Route[*pb.AppGetRichArticleListReq, *pb.AppGetRichArticleListResp]{
			NewRequest: func() *pb.AppGetRichArticleListReq {
				return &pb.AppGetRichArticleListReq{}
			},
			Do: svcCtx.AppMgmtService().AppGetRichArticleList,
		}
		conngateway.AddRoute("/v1/appmgmt/white/appGetRichArticleList", route)
	}
}
