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

type GetAllMSLuaConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSLuaConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSLuaConfigLogic {
	return &GetAllMSLuaConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSLuaConfigLogic) GetAllMSLuaConfig(in *pb.GetAllMSLuaConfigReq) (*pb.GetAllMSLuaConfigResp, error) {
	var models []*mgmtmodel.LuaConfig
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			if v == "" {
				continue
			}
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			case "isEnable":
				if v == "true" || v == "1" {
					wheres = append(wheres, xorm.Where("isEnable = ?", true))
				} else {
					wheres = append(wheres, xorm.Where("isEnable = ?", false))
				}
			case "name":
				wheres = append(wheres, xorm.Where("name like ?", "%"+v+"%"))
			case "desc":
				wheres = append(wheres, xorm.Where("desc like ?", "%"+v+"%"))
			case "time_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime >= ?", val))
			case "type":
				wheres = append(wheres, xorm.Where("type = ?", v))
			case "time_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime <= ?", val))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.LuaConfig{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("ListWithPagingOrder err: %v", err)
		return &pb.GetAllMSLuaConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.LuaConfig
	for _, model := range models {
		resp = append(resp, model.ToPB())
	}
	return &pb.GetAllMSLuaConfigResp{
		LuaConfigs: resp,
		Total:      count,
	}, nil
}
