package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendSms 发送短信
func (l *SendSmsLogic) SendSms(in *pb.ThirdSendSmsReq) (*pb.ThirdSendSmsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ThirdSendSmsResp{}, nil
}
