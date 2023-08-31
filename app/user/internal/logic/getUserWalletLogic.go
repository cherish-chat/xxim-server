package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWalletLogic {
	return &GetUserWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserWalletLogic) GetUserWallet(in *pb.GetUserWalletReq) (*pb.GetUserWalletResp, error) {
	dest := &usermodel.UserWallet{}
	l.svcCtx.Mysql().Model(dest).Where("userId = ?", in.UserId).First(dest)
	if dest.UserId == "" {
		// 创建一个钱包
		dest.UserId = in.UserId
		l.svcCtx.Mysql().Model(dest).Create(dest)
	}
	return &pb.GetUserWalletResp{
		UserWallet: dest.ToProto(),
	}, nil
}
