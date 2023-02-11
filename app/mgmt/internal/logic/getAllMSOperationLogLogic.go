package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSOperationLogLogic {
	return &GetAllMSOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSOperationLogLogic) GetAllMSOperationLog(in *pb.GetAllMSOperationLogReq) (*pb.GetAllMSOperationLogResp, error) {
	var models []*mgmtmodel.OperationLog
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.OperationLog{}, in.Page.Page, in.Page.Size, "reqTime DESC", wheres...)
	if err != nil {
		l.Errorf("ListWithPagingOrder err: %v", err)
		return &pb.GetAllMSOperationLogResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSOperationLog
	for _, model := range models {
		resp = append(resp, model.ToPB())
	}
	return &pb.GetAllMSOperationLogResp{
		OperationLogs: resp,
		Total:         count,
	}, nil
}
