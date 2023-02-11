package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtVersionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtVersionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtVersionDetailLogic {
	return &GetAppMgmtVersionDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtVersionDetailLogic) GetAppMgmtVersionDetail(in *pb.GetAppMgmtVersionDetailReq) (*pb.GetAppMgmtVersionDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Version{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtVersionDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtVersionDetailResp{AppMgmtVersion: model.ToPB()}, nil
}
