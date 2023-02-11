package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtShieldWordDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtShieldWordDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtShieldWordDetailLogic {
	return &GetAppMgmtShieldWordDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtShieldWordDetailLogic) GetAppMgmtShieldWordDetail(in *pb.GetAppMgmtShieldWordDetailReq) (*pb.GetAppMgmtShieldWordDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.ShieldWord{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtShieldWordDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtShieldWordDetailResp{AppMgmtShieldWord: model.ToPB()}, nil
}
