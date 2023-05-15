package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNoticeInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeInsertLogic {
	return &NoticeInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// NoticeInsert 插入公告
func (l *NoticeInsertLogic) NoticeInsert(in *pb.NoticeInsertReq) (*pb.NoticeInsertResp, error) {
	// todo: add your logic here and delete this line

	return &pb.NoticeInsertResp{}, nil
}
