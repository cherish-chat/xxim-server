package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtVpnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtVpnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtVpnLogic {
	return &AddAppMgmtVpnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtVpnLogic) AddAppMgmtVpn(in *pb.AddAppMgmtVpnReq) (*pb.AddAppMgmtVpnResp, error) {
	model := &appmgmtmodel.Vpn{
		Id:         appmgmtmodel.GetId(l.svcCtx.Mysql(), &appmgmtmodel.Vpn{}, 10000),
		Name:       in.AppMgmtVpn.Name,
		Platform:   in.AppMgmtVpn.Platform,
		Type:       in.AppMgmtVpn.Type,
		Ip:         in.AppMgmtVpn.Ip,
		Port:       int(in.AppMgmtVpn.Port),
		Username:   in.AppMgmtVpn.Username,
		Password:   in.AppMgmtVpn.Password,
		SecretKey:  in.AppMgmtVpn.SecretKey,
		CreateTime: time.Now().UnixMilli(),
	}
	err := model.Insert(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddAppMgmtVpnResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddAppMgmtVpnResp{}, nil
}
