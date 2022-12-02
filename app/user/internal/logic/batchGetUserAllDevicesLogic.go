package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserAllDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserAllDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserAllDevicesLogic {
	return &BatchGetUserAllDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUserAllDevicesLogic) BatchGetUserAllDevices(in *pb.BatchGetUserAllDevicesReq) (*pb.BatchGetUserAllDevicesResp, error) {
	var resp []*pb.BatchGetUserAllDevicesResp_AllDevices
	for _, userId := range in.UserIds {
		// get all
		val, err := l.svcCtx.Redis().ZrangeCtx(l.ctx, rediskey.UserDeviceMapping(userId), 0, -1)
		if err != nil {
			l.Logger.Errorf("get user %s all devices from redis error: %s", userId, err.Error())
			return &pb.BatchGetUserAllDevicesResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		resp = append(resp, &pb.BatchGetUserAllDevicesResp_AllDevices{
			UserId:    userId,
			DeviceIds: val,
		})
	}
	return &pb.BatchGetUserAllDevicesResp{
		AllDevices: resp,
	}, nil
}
