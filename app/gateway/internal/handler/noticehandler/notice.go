package noticehandler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/wrapper"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/grpc"
	"time"
)

// AckNoticeDataConfig ...
func AckNoticeDataConfig[REQ *pb.AckNoticeDataReq, RESP *pb.AckNoticeDataResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp] {
	return wrapper.Config[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp]{
		Do: func(ctx context.Context, in *pb.AckNoticeDataReq, opts ...grpc.CallOption) (*pb.AckNoticeDataResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.NoticeService().AckNoticeData(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "AckNoticeData", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.AckNoticeDataReq {
			return &pb.AckNoticeDataReq{}
		},
	}
}
