package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MapUserByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMapUserByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MapUserByIdsLogic {
	return &MapUserByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MapUserByIdsLogic) MapUserByIds(in *pb.MapUserByIdsReq) (*pb.MapUserByIdsResp, error) {
	usersByIds, err := usermodel.GetUsersByIds(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&usermodel.User{}), in.Ids)
	if err != nil {
		l.Errorf("MapUserByIdsLogic MapUserByIds err: %v", err)
		return &pb.MapUserByIdsResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	return &pb.MapUserByIdsResp{Users: utils.Slice2MapBytes[*usermodel.User](usersByIds)}, nil
}
