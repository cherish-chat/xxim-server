package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSRoleListLogic {
	return &GetAllMSRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSRoleListLogic) GetAllMSRoleList(in *pb.GetAllMSRoleListReq) (*pb.GetAllMSRoleListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*mgmtmodel.Role
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.Role{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSRoleList err: %v", err)
		return &pb.GetAllMSRoleListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSRole
	var mapMembers = make(map[string]int64)
	{
		var getMembersFs []func()
		for _, model := range models {
			m := model
			getMembersFs = append(getMembersFs, func() {
				var count int64
				l.svcCtx.Mysql().Model(&mgmtmodel.User{}).Where("roleId = ?", m.Id).Count(&count)
				mapMembers[m.Id] = count
			})
		}
		mr.FinishVoid(getMembersFs...)
	}
	for _, model := range models {
		role := model.ToPB()
		role.Member, _ = mapMembers[model.Id]
		resp = append(resp, role)
	}
	return &pb.GetAllMSRoleListResp{
		Roles: resp,
		Total: count,
	}, nil
}
