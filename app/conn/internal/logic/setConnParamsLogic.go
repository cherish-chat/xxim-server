package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"google.golang.org/grpc"
)

type SetConnParamsLogic struct {
	svcCtx *svc.ServiceContext
}

var singletonSetConnParamsLogic *SetConnParamsLogic

func NewSetConnParamsLogic(svcCtx *svc.ServiceContext) *SetConnParamsLogic {
	if singletonSetConnParamsLogic == nil {
		singletonSetConnParamsLogic = &SetConnParamsLogic{
			svcCtx: svcCtx,
		}
	}
	return singletonSetConnParamsLogic
}

func (l *SetConnParamsLogic) SetConnParams(ctx context.Context, req *pb.SetConnParamsReq, opts ...grpc.CallOption) (*pb.SetConnParamsResp, error) {
	return &pb.SetConnParamsResp{ConnParam: req.GetConnParam()}, nil
}

func (l *SetConnParamsLogic) Callback(ctx context.Context, resp *pb.SetConnParamsResp, c *types.UserConn) {
	if resp == nil || resp.ConnParam == nil {
		return
	}
	param := resp.GetConnParam()
	param.UserId = c.ConnParam.UserId
	param.Token = c.ConnParam.Token
	c.SetConnParams(param)
}
