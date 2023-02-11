package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSUserDetailLogic {
	return &GetMSUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSUserDetailLogic) GetMSUserDetail(in *pb.GetMSUserDetailReq) (*pb.GetMSUserDetailResp, error) {
	// 查询原模型
	model := &mgmtmodel.User{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.GetMSUserDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var loginRecord = &mgmtmodel.LoginRecord{}
	l.svcCtx.Mysql().Model(loginRecord).Where("userId = ?", model.Id).Order("time DESC").First(loginRecord)
	return &pb.GetMSUserDetailResp{
		User: model.ToPBMSUser(loginRecord),
	}, nil
}
