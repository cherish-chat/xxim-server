package connectionservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteDataLogic {
	return &WriteDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// WriteData 向用户推送数据
func (l *WriteDataLogic) WriteData(in *peerpb.GatewayWriteDataContent) (*peerpb.GatewayWriteDataContent, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GatewayWriteDataContent{}, nil
}
