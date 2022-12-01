package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AfterConnect conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	exp := rediskey.UserDeviceMappingExpire()
	// 用户与用户设备的对应关系
	err := xredis.ZaddEx(l.ctx, l.svcCtx.Redis(), rediskey.UserDeviceMapping(in.ConnParam.UserId), in.ConnParam.DeviceId, utils.AnyToInt64(in.ConnectedAt), exp)
	if err != nil {
		l.Logger.Errorf("set user %s device %s to redis error: %s", in.ConnParam.UserId, in.ConnParam.DeviceId, err.Error())
		return pb.NewRetryErrorResp(), err
	}
	// 查询有没有超过 exp 的设备，有的话删除
	_, err = l.svcCtx.Redis().ZremrangebyscoreCtx(l.ctx, rediskey.UserDeviceMapping(in.ConnParam.UserId), 0, utils.AnyToInt64(in.ConnectedAt)-exp)
	if err != nil {
		l.Logger.Errorf("set user %s device %s to redis error: %s", in.ConnParam.UserId, in.ConnParam.DeviceId, err.Error())
		return pb.NewRetryErrorResp(), err
	}
	// 查询设备与用户的对应关系
	val, err := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.DeviceUserMapping(in.ConnParam.DeviceId))
	if err != nil {
		// nil
		if !errors.Is(err, redis.Nil) {
			l.Logger.Errorf("get device %s user from redis error: %s", in.ConnParam.DeviceId, err.Error())
			return pb.NewRetryErrorResp(), err
		}
	}
	if val != "" && val != in.ConnParam.UserId {
		// 说明设备已经被其他用户登录，需要踢掉
		_, err = l.svcCtx.Redis().ZremCtx(l.ctx, rediskey.UserDeviceMapping(val), in.ConnParam.DeviceId)
		if err != nil {
			l.Logger.Errorf("remove device %s from user %s error: %s", in.ConnParam.DeviceId, val, err.Error())
			return pb.NewRetryErrorResp(), err
		}
	}
	// 设备与用户的对应关系
	err = l.svcCtx.Redis().SetexCtx(l.ctx, rediskey.DeviceUserMapping(in.ConnParam.DeviceId), in.ConnParam.UserId, int(exp))
	if err != nil {
		l.Logger.Errorf("set device %s user %s to redis error: %s", in.ConnParam.DeviceId, in.ConnParam.UserId, err.Error())
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
