package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSLoginRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSLoginRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSLoginRecordLogic {
	return &GetAllMSLoginRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSLoginRecordLogic) GetAllMSLoginRecord(in *pb.GetAllMSLoginRecordReq) (*pb.GetAllMSLoginRecordResp, error) {
	var models []*mgmtmodel.LoginRecord
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.LoginRecord{}, in.Page.Page, in.Page.Size, "time DESC", wheres...)
	if err != nil {
		l.Errorf("ListWithPagingOrder err: %v", err)
		return &pb.GetAllMSLoginRecordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSLoginRecord
	for _, model := range models {
		resp = append(resp, model.ToPB())
	}
	return &pb.GetAllMSLoginRecordResp{
		LoginRecords: resp,
		Total:        count,
	}, nil
}
