package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSIpWhiteListDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSIpWhiteListDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSIpWhiteListDetailLogic {
	return &GetMSIpWhiteListDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSIpWhiteListDetailLogic) GetMSIpWhiteListDetail(in *pb.GetMSIpWhiteListDetailReq) (*pb.GetMSIpWhiteListDetailResp, error) {
	// 查询原模型
	model := &mgmtmodel.MSIPWhitelist{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetMSIpWhiteListDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSIpWhiteListDetailResp{
		IpWhiteList: model.ToPB(),
	}, nil
}
