package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSUserListLogic {
	return &GetAllMSUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSUserListLogic) GetAllMSUserList(in *pb.GetAllMSUserListReq) (*pb.GetAllMSUserListResp, error) {
	var models []*mgmtmodel.User
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
		for k, v := range in.Filter {
			switch k {
			case "username":
				wheres = append(wheres, xorm.Where("id = ?", v))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.User{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSRoleList err: %v", err)
		return &pb.GetAllMSUserListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSUser
	var mapLatestLoginRecord = make(map[string]*mgmtmodel.LoginRecord)
	{
		var latestLoginRecords []*mgmtmodel.LoginRecord
		var getLatestLoginRecordFs []func()
		for _, model := range models {
			m := model
			getLatestLoginRecordFs = append(getLatestLoginRecordFs, func() {
				var loginRecord = &mgmtmodel.LoginRecord{}
				l.svcCtx.Mysql().Model(loginRecord).Where("userId = ?", m.Id).Order("time DESC").First(loginRecord)
				if loginRecord.Id != "" {
					latestLoginRecords = append(latestLoginRecords, loginRecord)
				}
			})
		}
		mr.FinishVoid(getLatestLoginRecordFs...)
		for _, loginRecord := range latestLoginRecords {
			mapLatestLoginRecord[loginRecord.UserId] = loginRecord
		}
	}
	var mapRole = make(map[string]*mgmtmodel.Role)
	{
		var roles []*mgmtmodel.Role
		l.svcCtx.Mysql().Model(&mgmtmodel.Role{}).Find(&roles)
		for _, role := range roles {
			mapRole[role.Id] = role
		}
	}
	for _, model := range models {
		record, ok := mapLatestLoginRecord[model.Id]
		if !ok {
			record = &mgmtmodel.LoginRecord{}
		}
		role, ok := mapRole[model.RoleId]
		if !ok {
			role = &mgmtmodel.Role{}
		}
		user := model.ToPBMSUser(record)
		user.Role = role.Name
		resp = append(resp, user)
	}
	return &pb.GetAllMSUserListResp{
		Users: resp,
		Total: count,
	}, nil
}
