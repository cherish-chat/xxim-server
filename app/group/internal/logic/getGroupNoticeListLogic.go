package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupNoticeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupNoticeListLogic {
	return &GetGroupNoticeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupNoticeList 获取群公告列表
func (l *GetGroupNoticeListLogic) GetGroupNoticeList(in *pb.GetGroupNoticeListReq) (*pb.GetGroupNoticeListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupNoticeListResp{}, nil
}
