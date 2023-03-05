package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterNotice(svcCtx *svc.ServiceContext) {
	// AckNoticeDataReq AckNoticeDataResp
	{
		route := conngateway.Route[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp]{
			NewRequest: func() *pb.AckNoticeDataReq {
				return &pb.AckNoticeDataReq{}
			},
			Do: svcCtx.NoticeService().AckNoticeData,
		}
		conngateway.AddRoute("/v1/notice/ackNoticeData", route)
	}
}
