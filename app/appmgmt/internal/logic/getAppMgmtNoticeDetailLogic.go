package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtNoticeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtNoticeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtNoticeDetailLogic {
	return &GetAppMgmtNoticeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtNoticeDetailLogic) GetAppMgmtNoticeDetail(in *pb.GetAppMgmtNoticeDetailReq) (*pb.GetAppMgmtNoticeDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Notice{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtNoticeDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtNoticeDetailResp{AppMgmtNotice: model.ToPB()}, nil
}
