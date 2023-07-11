package callbackservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeRequestLogic {
	return &UserBeforeRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserBeforeRequest 用户请求前的回调
func (l *UserBeforeRequestLogic) UserBeforeRequest(in *pb.UserBeforeRequestReq) (*pb.UserBeforeRequestResp, error) {
	return &pb.UserBeforeRequestResp{
		Header: i18n.NewOkHeader(),
	}, nil
}
