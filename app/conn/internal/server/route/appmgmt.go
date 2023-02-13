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
}
