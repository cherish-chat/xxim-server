package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
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

func (l *SetConnParamsLogic) SetConnParams(ctx context.Context, req *pb.SetCxnParamsReq, opts ...grpc.CallOption) (*pb.SetCxnParamsResp, error) {
	return &pb.SetCxnParamsResp{CxnParams: req.GetCxnParams()}, nil
}

func (l *SetConnParamsLogic) Callback(ctx context.Context, resp *pb.SetCxnParamsResp, c *types.UserConn) {
	if resp == nil || resp.CxnParams == nil {
		return
	}
	param := resp.GetCxnParams()
	c.SetConnParams(&pb.ConnParam{
		UserId:      c.ConnParam.UserId,
		Token:       c.ConnParam.Token,
		DeviceId:    param.DeviceId,
		Platform:    param.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: param.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		PodIp:       utils.GetPodIp(),
		DeviceModel: param.DeviceModel,
		OsVersion:   param.OsVersion,
		AppVersion:  param.AppVersion,
		Language:    param.Language,
	})
}
