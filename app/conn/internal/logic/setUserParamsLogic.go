package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

type SetUserParamsLogic struct {
	svcCtx *svc.ServiceContext
}

var singletonSetUserParamsLogic *SetUserParamsLogic

func NewSetUserParamsLogic(svcCtx *svc.ServiceContext) *SetUserParamsLogic {
	if singletonSetUserParamsLogic == nil {
		singletonSetUserParamsLogic = &SetUserParamsLogic{
			svcCtx: svcCtx,
		}
	}
	return singletonSetUserParamsLogic
}

func (l *SetUserParamsLogic) SetUserParams(ctx context.Context, req *pb.SetUserParamsReq, opts ...grpc.CallOption) (*pb.SetUserParamsResp, error) {
	return &pb.SetUserParamsResp{
		UserId: req.GetUserId(),
		Token:  req.GetToken(),
		Ext:    req.GetExt(),
	}, nil
}

func (l *SetUserParamsLogic) Callback(ctx context.Context, resp *pb.SetUserParamsResp, c *types.UserConn) {
	if resp == nil || resp.UserId == "" || resp.Token == "" {
		return
	}
	// 鉴权
	code, err := GetConnLogic().BeforeConnect(ctx, types.ConnParam{
		UserId:      resp.UserId,
		Token:       resp.Token,
		DeviceId:    c.ConnParam.DeviceId,
		DeviceModel: c.ConnParam.DeviceModel,
		OsVersion:   c.ConnParam.OsVersion,
		AppVersion:  c.ConnParam.AppVersion,
		Language:    c.ConnParam.Language,
		Platform:    c.ConnParam.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: c.ConnParam.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		Timestamp:   c.ConnParam.Timestamp,
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("BeforeConnect err: %v", err)
		return
	}
	if code != 0 {
		logx.WithContext(ctx).Errorf("BeforeConnect code: %d", code)
		return
	}
	c.SetConnParams(&pb.ConnParam{
		UserId:      resp.UserId,
		Token:       resp.Token,
		DeviceId:    c.ConnParam.DeviceId,
		Platform:    c.ConnParam.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: c.ConnParam.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		PodIp:       utils.GetPodIp(),
		DeviceModel: c.ConnParam.DeviceModel,
		OsVersion:   c.ConnParam.OsVersion,
		AppVersion:  c.ConnParam.AppVersion,
		Language:    c.ConnParam.Language,
	})
	GetConnLogic().AddSubscriber(c)
}
