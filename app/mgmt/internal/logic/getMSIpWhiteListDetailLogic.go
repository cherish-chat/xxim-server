package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &pb.GetMSIpWhiteListDetailResp{}, nil
}
