package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLatestConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLatestConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLatestConnLogic {
	return &GetUserLatestConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLatestConnLogic) GetUserLatestConn(in *pb.GetUserLatestConnReq) (*pb.GetUserLatestConnResp, error) {
	model := &immodel.UserConnectRecord{}
	err := l.svcCtx.Mongo().Collection(model).Find(l.ctx, bson.M{
		"userId": in.UserId,
	}).Sort("-connectTime").One(model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &pb.GetUserLatestConnResp{}, nil
		} else {
			l.Errorf("find user latest connect record failed, err: %v", err)
			return &pb.GetUserLatestConnResp{}, err
		}
	}
	return &pb.GetUserLatestConnResp{
		UserId:         model.UserId,
		Ip:             model.Ips,
		IpRegion:       model.IpRegion.Pb(),
		ConnectedAt:    utils.AnyToString(model.ConnectTime),
		DisconnectedAt: utils.AnyToString(model.DisconnectTime),
		Platform:       model.Platform,
		DeviceId:       model.DeviceId,
	}, nil
}
