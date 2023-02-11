package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSOperationLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSOperationLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSOperationLogDetailLogic {
	return &GetMSOperationLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSOperationLogDetailLogic) GetMSOperationLogDetail(in *pb.GetMSOperationLogDetailReq) (*pb.GetMSOperationLogDetailResp, error) {
	// 查询原模型
	model := &mgmtmodel.OperationLog{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetMSOperationLogDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSOperationLogDetailResp{
		OperationLog: model.ToPB(),
	}, nil
}
