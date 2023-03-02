package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterIm(svcCtx *svc.ServiceContext) {
	// group
	{
		// UpdateConvSettingReq UpdateConvSettingResp
		{
			route := conngateway.Route[*pb.UpdateConvSettingReq, *pb.UpdateConvSettingResp]{
				NewRequest: func() *pb.UpdateConvSettingReq {
					return &pb.UpdateConvSettingReq{}
				},
				Do: svcCtx.ImService().UpdateConvSetting,
			}
			conngateway.AddRoute("/v1/im/updateConvSetting", route)
		}
	}
}
