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
	return &pb.SetCxnParamsResp{
		Platform:    req.GetPlatform(),
		DeviceId:    req.GetDeviceId(),
		DeviceModel: req.GetDeviceModel(),
		OsVersion:   req.GetOsVersion(),
		AppVersion:  req.GetAppVersion(),
		Language:    req.GetLanguage(),
		NetworkUsed: req.GetNetworkUsed(),
		Ext:         req.GetExt(),
	}, nil
}

func (l *SetConnParamsLogic) Callback(ctx context.Context, resp *pb.SetCxnParamsResp, c *types.UserConn) {
	c.SetConnParams(&pb.ConnParam{
		UserId:      c.ConnParam.UserId,
		Token:       c.ConnParam.Token,
		DeviceId:    resp.DeviceId,
		Platform:    resp.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: resp.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		PodIp:       utils.GetPodIp(),
		DeviceModel: resp.DeviceModel,
		OsVersion:   resp.OsVersion,
		AppVersion:  resp.AppVersion,
		Language:    resp.Language,
	})
}
