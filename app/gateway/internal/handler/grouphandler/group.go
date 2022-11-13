package grouphandler

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

// CreateGroupConfig ...
func CreateGroupConfig[REQ *pb.CreateGroupReq, RESP *pb.CreateGroupResp](svcCtx *svc.ServiceContext) wrapper.Config[*pb.CreateGroupReq, *pb.CreateGroupResp] {
	return wrapper.Config[*pb.CreateGroupReq, *pb.CreateGroupResp]{
		Do: func(ctx context.Context, in *pb.CreateGroupReq, opts ...grpc.CallOption) (*pb.CreateGroupResp, error) {
			requestTime := time.Now()
			resp, err := svcCtx.GroupService().CreateGroup(ctx, in, opts...)
			go logic.NewApiLogLogic(ctx, svcCtx).ApiLog(in.GetRequester(), "CreateGroup", resp.GetCommonResp(), utils.AnyToString(in), utils.AnyToString(resp), requestTime, time.Now(), err)
			return resp, err
		},
		NewRequest: func() *pb.CreateGroupReq {
			return &pb.CreateGroupReq{}
		},
	}
}
