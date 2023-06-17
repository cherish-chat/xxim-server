package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeConnectLogic {
	return &UserBeforeConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserBeforeConnect 用户连接前的回调
func (l *UserBeforeConnectLogic) UserBeforeConnect(in *pb.UserBeforeConnectReq) (*pb.UserBeforeConnectResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserBeforeConnectResp{
		Success:     true,
		CloseCode:   0,
		CloseReason: "",
	}, nil
}
