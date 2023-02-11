package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDefaultConvDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDefaultConvDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDefaultConvDetailLogic {
	return &GetUserDefaultConvDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserDefaultConvDetailLogic) GetUserDefaultConvDetail(in *pb.GetUserDefaultConvDetailReq) (*pb.GetUserDefaultConvDetailResp, error) {
	// 查询原模型
	model := &usermodel.DefaultConv{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetUserDefaultConvDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetUserDefaultConvDetailResp{UserDefaultConv: model.ToPB()}, nil
}
