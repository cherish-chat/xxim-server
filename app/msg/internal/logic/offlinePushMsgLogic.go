package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OfflinePushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOfflinePushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OfflinePushMsgLogic {
	return &OfflinePushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// OfflinePushMsg 离线推送消息
func (l *OfflinePushMsgLogic) OfflinePushMsg(in *pb.OfflinePushMsgReq) (*pb.OfflinePushMsgResp, error) {
	// 是否已推送
	val, err := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.OfflinePushMsgListKey(in.UniqueId))
	if err == nil && val != "" {
		// 推送过
		return &pb.OfflinePushMsgResp{}, nil
	}
	if l.svcCtx.Config.MobAlias == "deviceId" {
		{
			batchGetUserAllDevicesResp, err := l.svcCtx.UserService().BatchGetUserAllDevices(l.ctx, &pb.BatchGetUserAllDevicesReq{
				UserIds: in.UserIds,
			})
			if err != nil {
				l.Errorf("BatchGetUserAllDevices err: %v", err)
				return &pb.OfflinePushMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
			deviceIds := make([]string, 0)
			for _, device := range batchGetUserAllDevicesResp.AllDevices {
				deviceIds = append(deviceIds, device.DeviceIds...)
			}
			if len(deviceIds) == 0 {
				return &pb.OfflinePushMsgResp{}, nil
			}
			resp, err := l.svcCtx.MobPush.Push(l.ctx, utils.Set(deviceIds), in.Title, in.Content)
			if err != nil {
				l.Errorf("MobPush err: %v", err)
				return &pb.OfflinePushMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
			l.Infof("MobPush resp: %v", resp)
		}
	} else if l.svcCtx.Config.MobAlias == "userId" {
		if len(in.UserIds) > 0 {
			resp, err := l.svcCtx.MobPush.Push(l.ctx, utils.Set(in.UserIds), in.Title, in.Content)
			if err != nil {
				l.Errorf("MobPush err: %v", err)
				return &pb.OfflinePushMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
			l.Infof("MobPush resp: %v", resp)
		}
	}
	err = l.svcCtx.Redis().SetexCtx(l.ctx, rediskey.OfflinePushMsgListKey(in.UniqueId), in.Content, 10)
	if err != nil {
		l.Errorf("Redis SetexCtx err: %v", err)
		return &pb.OfflinePushMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.OfflinePushMsgResp{}, nil
}
