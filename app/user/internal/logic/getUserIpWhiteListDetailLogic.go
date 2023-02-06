package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserIpWhiteListDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserIpWhiteListDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserIpWhiteListDetailLogic {
	return &GetUserIpWhiteListDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserIpWhiteListDetailLogic) GetUserIpWhiteListDetail(in *pb.GetUserIpWhiteListDetailReq) (*pb.GetUserIpWhiteListDetailResp, error) {
	// 查询原模型
	model := &usermodel.IpWhiteList{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetUserIpWhiteListDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetUserIpWhiteListDetailResp{UserIpList: model.ToPB()}, nil
}
