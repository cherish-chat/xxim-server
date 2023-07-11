package internalservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLongConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLongConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLongConnectionLogic {
	return &ListLongConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListLongConnection 获取长连接列表
func (l *ListLongConnectionLogic) ListLongConnection(in *peerpb.ListLongConnectionReq) (*peerpb.ListLongConnectionResp, error) {
	// todo: add your logic here and delete this line
	l.Infof("获取长连接列表")
	return &peerpb.ListLongConnectionResp{}, nil
}
