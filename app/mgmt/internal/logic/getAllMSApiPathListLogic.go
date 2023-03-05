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

type GetAllMSApiPathListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSApiPathListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSApiPathListLogic {
	return &GetAllMSApiPathListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSApiPathListLogic) GetAllMSApiPathList(in *pb.GetAllMSApiPathListReq) (*pb.GetAllMSApiPathListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*mgmtmodel.ApiPath
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			if v == "" {
				continue
			}
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			case "logEnable":
				if v == "true" || v == "1" {
					wheres = append(wheres, xorm.Where("logEnable = ?", true))
				} else {
					wheres = append(wheres, xorm.Where("logEnable = ?", false))
				}
			case "title":
				wheres = append(wheres, xorm.Where("title like ?", "%"+v+"%"))
			case "path":
				wheres = append(wheres, xorm.Where("path like ?", "%"+v+"%"))
			case "time_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime >= ?", val))
			case "time_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime <= ?", val))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.ApiPath{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSApiPathList err: %v", err)
		return &pb.GetAllMSApiPathListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSApiPath
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllMSApiPathListResp{
		ApiPaths: resp,
		Total:    count,
	}, nil
}
