package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserLatestConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserLatestConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserLatestConnLogic {
	return &BatchGetUserLatestConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUserLatestConnLogic) BatchGetUserLatestConn(in *pb.BatchGetUserLatestConnReq) (*pb.BatchGetUserLatestConnResp, error) {
	rediskeys := make([]string, 0, len(in.UserIds))
	for _, userId := range in.UserIds {
		rediskeys = append(rediskeys, rediskey.LatestConnectRecord(userId))
	}
	if len(rediskeys) == 0 {
		return &pb.BatchGetUserLatestConnResp{}, nil
	}
	latestConnRecords, err := l.svcCtx.Redis().MgetCtx(l.ctx, rediskeys...)
	if err != nil {
		l.Errorf("BatchGetUserLatestConnLogic BatchGetUserLatestConn err: %v", err)
		return &pb.BatchGetUserLatestConnResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	var resp = make([]*pb.GetUserLatestConnResp, 0)
	for _, latestConnRecord := range latestConnRecords {
		if latestConnRecord == "" {
			continue
		}
		model := &immodel.UserConnectRecord{}
		err := json.Unmarshal([]byte(latestConnRecord), model)
		if err != nil {
			l.Errorf("BatchGetUserLatestConnLogic BatchGetUserLatestConn err: %v", err)
			continue
		}
		resp = append(resp, &pb.GetUserLatestConnResp{
			UserId:         model.UserId,
			Ip:             model.Ips,
			IpRegion:       model.IpRegion.Pb(),
			ConnectedAt:    utils.AnyToString(model.ConnectTime),
			DisconnectedAt: utils.AnyToString(model.DisconnectTime),
			Platform:       model.Platform,
			DeviceId:       model.DeviceId,
		})
	}
	return &pb.BatchGetUserLatestConnResp{UserLatestConns: resp}, nil
}
