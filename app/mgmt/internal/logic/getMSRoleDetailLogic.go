package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSRoleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSRoleDetailLogic {
	return &GetMSRoleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSRoleDetailLogic) GetMSRoleDetail(in *pb.GetMSRoleDetailReq) (*pb.GetMSRoleDetailResp, error) {
	// 查询原模型
	model := &mgmtmodel.Role{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetMSRoleDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSRoleDetailResp{Role: model.ToPB()}, nil
}
