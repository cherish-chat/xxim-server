package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtShieldWordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtShieldWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtShieldWordLogic {
	return &GetAllAppMgmtShieldWordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtShieldWordLogic) GetAllAppMgmtShieldWord(in *pb.GetAllAppMgmtShieldWordReq) (*pb.GetAllAppMgmtShieldWordResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*appmgmtmodel.ShieldWord
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			if v == "" {
				continue
			}
			switch k {
			case "word":
				wheres = append(wheres, xorm.Where("word LIKE ?", v+"%"))
			case "time_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime >= ?", val))
			case "time_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime <= ?", val))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.ShieldWord{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllAppMgmtShieldWordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtShieldWord
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllAppMgmtShieldWordResp{
		AppMgmtShieldWords: resp,
		Total:              count,
	}, nil
}
