package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"

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
	err := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.UserId).Order("connectTime desc").First(model).Error
	if err != nil {
		if xorm.RecordNotFound(err) {
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
