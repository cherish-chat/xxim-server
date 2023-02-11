package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserIpWhiteListLogic {
	return &DeleteUserIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserIpWhiteListLogic) DeleteUserIpWhiteList(in *pb.DeleteUserIpWhiteListReq) (*pb.DeleteUserIpWhiteListResp, error) {
	model := &usermodel.IpWhiteList{}
	err := l.svcCtx.Mysql().Model(model).Where("id in (?)", in.Ids).Delete(model).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteUserIpWhiteListResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteUserIpWhiteListResp{}, nil
}
