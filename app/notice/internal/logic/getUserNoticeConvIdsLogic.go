package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNoticeConvIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNoticeConvIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNoticeConvIdsLogic {
	return &GetUserNoticeConvIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserNoticeConvIds 获取用户所有的通知号
func (l *GetUserNoticeConvIdsLogic) GetUserNoticeConvIds(in *pb.GetUserNoticeConvIdsReq) (*pb.GetUserNoticeConvIdsResp, error) {
	return &pb.GetUserNoticeConvIdsResp{
		ConvIds: noticemodel.DefaultConvIds,
	}, nil
}
