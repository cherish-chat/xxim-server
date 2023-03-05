package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// 验证
	_, ok := pb.PlatformMap[c.ConnParam.Platform]
	if !ok {
		logx.WithContext(ctx).Errorf("platform err: %v, conn: %s, ip: %s", c.ConnParam.Platform, utils.AnyToString(c.ConnParam), ip2region.Ip2Region(c.ConnParam.Ips).String())
		c.Conn.Close(types.WebsocketStatusCodePlatformFailed(), "platform failed")
		return
	}
	if c.ConnParam.DeviceId == "" {
		logx.WithContext(ctx).Errorf("deviceId err: conn: %s, ip: %s", utils.AnyToString(c.ConnParam), ip2region.Ip2Region(c.ConnParam.Ips).String())
		c.Conn.Close(types.WebsocketStatusCodeDeviceIdFailed(), "deviceId failed")
		return
	}
	if resp == nil || resp.UserId == "" || resp.Token == "" {
		c.Conn.Close(types.WebsocketStatusCodeAuthFailed(3000), "鉴权失败: id或token为空")
		return
	}
	if resp.UserId == c.ConnParam.UserId && resp.Token == c.ConnParam.Token {
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
		statusError, ok := status.FromError(err)
		if ok && statusError.Code() == codes.Unauthenticated {
			// 被封禁
			c.Conn.Close(types.WebsocketStatusCodeAuthFailed(3000), statusError.Message())
			return
		}
		// 断开连接
		c.Conn.Close(types.WebsocketStatusCodeAuthFailed(1), "认证失败，请重新登录")
		return
	}
	if code != 0 {
		logx.WithContext(ctx).Errorf("BeforeConnect code: %d", code)
		// 断开连接
		c.Conn.Close(types.WebsocketStatusCodeAuthFailed(code), "您的账号在其他设备登录")
		return
	}
	kickResp, err := l.svcCtx.ImService().KickUserConn(ctx, &pb.KickUserConnReq{GetUserConnReq: &pb.GetUserConnReq{
		UserIds:   []string{c.ConnParam.UserId},
		Platforms: []string{c.ConnParam.Platform},
		Devices:   []string{c.ConnParam.DeviceId},
	}})
	if err != nil {
		logx.WithContext(ctx).Errorf("KickUserConn err: %v", err)
		c.Conn.Close(types.WebsocketStatusCodeAuthFailed(1), "服务错误，请稍后再试")
		return
	}
	if kickResp.GetCommonResp().GetCode() != pb.CommonResp_Success {
		logx.WithContext(ctx).Errorf("KickUserConn code: %d", kickResp.GetCommonResp().GetCode())
		c.Conn.Close(types.WebsocketStatusCodeAuthFailed(1), "服务错误，请稍后再试")
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
		AesKey:      c.ConnParam.AesKey,
		AesIv:       c.ConnParam.AesIv,
	})
	GetConnLogic().AddSubscriber(c)
}
