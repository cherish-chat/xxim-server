package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtLinkDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtLinkDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtLinkDetailLogic {
	return &GetAppMgmtLinkDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtLinkDetailLogic) GetAppMgmtLinkDetail(in *pb.GetAppMgmtLinkDetailReq) (*pb.GetAppMgmtLinkDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Link{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtLinkDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtLinkDetailResp{AppMgmtLink: model.ToPB()}, nil
}
