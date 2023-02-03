package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSApiPathDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSApiPathDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSApiPathDetailLogic {
	return &GetMSApiPathDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSApiPathDetailLogic) GetMSApiPathDetail(in *pb.GetMSApiPathDetailReq) (*pb.GetMSApiPathDetailResp, error) {
	// 查询原模型
	model := &mgmtmodel.ApiPath{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetMSApiPathDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSApiPathDetailResp{ApiPath: model.ToPB()}, nil
}
