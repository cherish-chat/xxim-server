package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginMSLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginMSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginMSLogic {
	return &LoginMSLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginMSLogic) LoginMS(in *pb.LoginMSReq) (*pb.LoginMSResp, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginMSResp{}, nil
}
