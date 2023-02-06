package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/mr"

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
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &usermodel.User{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllUserModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var lastLoginRecordMap = make(map[string]*usermodel.LoginRecord)
	{
		var getLoginRecordFunc []func()
		for _, model := range models {
			m := model
			getLoginRecordFunc = append(getLoginRecordFunc, func() {
				var lastLoginRecord usermodel.LoginRecord
				l.svcCtx.Mysql().Where("userId = ?", m.Id).Order("time DESC").First(&lastLoginRecord)
				lastLoginRecordMap[m.Id] = &lastLoginRecord
			})
		}
		mr.FinishVoid(getLoginRecordFunc...)
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
