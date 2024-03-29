package emailservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCodeVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailCodeVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCodeVerifyLogic {
	return &EmailCodeVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EmailCodeVerify 验证邮件
func (l *EmailCodeVerifyLogic) EmailCodeVerify(in *pb.EmailCodeVerifyReq) (*pb.EmailCodeVerifyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.EmailCodeVerifyResp{Success: true}, nil
}
