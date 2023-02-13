package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllGroupModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllGroupModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllGroupModelLogic {
	return &GetAllGroupModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAllGroupModel 获取所有群组
func (l *GetAllGroupModelLogic) GetAllGroupModel(in *pb.GetAllGroupModelReq) (*pb.GetAllGroupModelResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*groupmodel.Group
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			switch k {
			case "id":
				wheres = append(wheres, xorm.Where("id = ?", v))
			case "name":
				wheres = append(wheres, xorm.Where("name LIKE ?", v+"%"))
			case "owner":
				wheres = append(wheres, xorm.Where("owner = ?", v))
			case "invitationCode":
				wheres = append(wheres, xorm.Where("invitation_code = ?", v))
			case "status":
				if v == "normal" {
					wheres = append(wheres, xorm.Where("dismissTime = ? AND allMute = ?", 0, false))
				} else if v == "dismiss" {
					wheres = append(wheres, xorm.Where("dismissTime > ?", 0))
				} else if v == "mute" {
					wheres = append(wheres, xorm.Where("allMute = ?", true))
				}
			case "memberCount_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("memberCount >= ?", val))
			case "memberCount_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("memberCount <= ?", val))
			case "createTime_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime >= ?", val))
			case "createTime_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime <= ?", val))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &groupmodel.Group{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.GroupModel
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllGroupModelResp{
		GroupModels: resp,
		Total:       count,
	}, nil
}
