package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
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
			if v == "" {
				continue
			}
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			case "operationType":
				wheres = append(wheres, xorm.Where("operationType = ?", v))
			case "operationTitle":
				wheres = append(wheres, xorm.Where("operationTitle LIKE ?", v+"%"))
			case "resultSuccess":
				if v == "1" || v == "true" {
					wheres = append(wheres, xorm.Where("resultSuccess = ?", true))
				} else {
					wheres = append(wheres, xorm.Where("resultSuccess = ?", false))
				}
			case "reqIp":
				wheres = append(wheres, xorm.Where("reqIp LIKE ?", v+"%"))
			case "operator":
				wheres = append(wheres, xorm.Where("operator = ?", v))
			case "reqTime_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("reqTime >= ?", val))
			case "reqTime_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("reqTime <= ?", val))
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
