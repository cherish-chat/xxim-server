package imhandler

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

// GetAppSystemConfigConfig ...
func GetAppSystemConfigConfig[REQ *pb.GetAppSystemConfigReq, RESP *pb.GetAppSystemConfigResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.GetAppSystemConfigReq, *pb.GetAppSystemConfigResp] {
	return wrapper.Config[*pb.GetAppSystemConfigReq, *pb.GetAppSystemConfigResp]{
		Do: func(ctx context.Context, in *pb.GetAppSystemConfigReq, opts ...grpc.CallOption) (*pb.GetAppSystemConfigResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.ImService().GetAppSystemConfig(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetCommonReq(), "GetAppSystemConfig", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.GetAppSystemConfigReq {
			return &pb.GetAppSystemConfigReq{}
		},
	}
}
