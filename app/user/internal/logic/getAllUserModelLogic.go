package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/mr"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUserModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserModelLogic {
	return &GetAllUserModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUserModelLogic) GetAllUserModel(in *pb.GetAllUserModelReq) (*pb.GetAllUserModelResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*usermodel.User
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
			case "nickname":
				wheres = append(wheres, xorm.Where("nickname LIKE ?", v+"%"))
			case "role":
				role := int32(utils.AnyToInt64(v))
				switch role {
				case 1:
					role = int32(usermodel.RoleUser)
				case 2:
					role = int32(usermodel.RoleService)
				case 3:
					role = int32(usermodel.RoleGuest)
				case 4:
					role = int32(usermodel.RoleZombie)
				}
				wheres = append(wheres, xorm.Where("role = ?", role))
			case "invitationCode":
				wheres = append(wheres, xorm.Where("invitation_code = ?", v))
			case "status":
				if v == "normal" {
					wheres = append(wheres, xorm.Where("unblockTime < ?", time.Now().UnixMilli()))
				} else if v == "block" {
					wheres = append(wheres, xorm.Where("unblockTime > ?", time.Now().UnixMilli()))
				}
			case "createTime_gte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime >= ?", val))
			case "createTime_lte":
				val := utils.AnyToInt64(v)
				wheres = append(wheres, xorm.Where("createTime <= ?", val))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &usermodel.User{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllUserModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var lastLoginRecordMap = make(map[string]*usermodel.LoginRecord)
	var lastLoginRecords []*usermodel.LoginRecord
	{
		var getLoginRecordFunc []func()
		for _, model := range models {
			m := model
			getLoginRecordFunc = append(getLoginRecordFunc, func() {
				var lastLoginRecord usermodel.LoginRecord
				l.svcCtx.Mysql().Where("userId = ?", m.Id).Order("time DESC").First(&lastLoginRecord)
				lastLoginRecords = append(lastLoginRecords, &lastLoginRecord)
			})
		}
		mr.FinishVoid(getLoginRecordFunc...)
		for _, record := range lastLoginRecords {
			lastLoginRecordMap[record.UserId] = record
		}
	}
	var resp []*pb.UserModel
	for _, model := range models {
		toPB := model.ToPB()
		if lastLoginRecord, ok := lastLoginRecordMap[model.Id]; !ok {
			toPB.LastLoginRecord = nil
		} else {
			toPB.LastLoginRecord = lastLoginRecord.ToPB()
		}
		resp = append(resp, toPB)
	}
	return &pb.GetAllUserModelResp{
		UserModelList: resp,
		Total:         count,
	}, nil
}
