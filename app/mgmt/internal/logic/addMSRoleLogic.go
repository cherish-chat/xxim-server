package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSRoleLogic {
	return &AddMSRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSRoleLogic) AddMSRole(in *pb.AddMSRoleReq) (*pb.AddMSRoleResp, error) {
	model := &mgmtmodel.Role{
		Id:         mgmtmodel.GetId(l.svcCtx.Mysql(), &mgmtmodel.Role{}, 1000),
		Name:       in.Role.Name,
		Remark:     in.Role.Remark,
		IsDisable:  in.Role.IsDisable,
		Sort:       in.Role.Sort,
		MenuIds:    strings.Join(in.Role.MenuIds, ","),
		ApiPathIds: strings.Join(in.Role.ApiPathIds, ","),
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}
	err := l.svcCtx.Mysql().Model(model).Create(model).Error
	if err != nil {
		l.Errorf("添加失败: %v", err)
		return &pb.AddMSRoleResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AddMSRoleResp{}, nil
}
