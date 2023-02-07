package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupModelDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupModelDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupModelDetailLogic {
	return &GetGroupModelDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupModelDetail 获取群组详情
func (l *GetGroupModelDetailLogic) GetGroupModelDetail(in *pb.GetGroupModelDetailReq) (*pb.GetGroupModelDetailResp, error) {
	// 查询原模型
	model := &groupmodel.Group{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetGroupModelDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetGroupModelDetailResp{GroupModel: model.ToPB()}, nil
}
