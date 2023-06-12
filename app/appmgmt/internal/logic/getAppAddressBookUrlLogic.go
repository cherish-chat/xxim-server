package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppAddressBookUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppAddressBookUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppAddressBookUrlLogic {
	return &GetAppAddressBookUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppAddressBookUrlLogic) GetAppAddressBookUrl(in *pb.GetAppAddressBookUrlReq) (*pb.GetAppAddressBookUrlResp, error) {
	url, err := l.svcCtx.Redis().Get(rediskey.AppAddressBookUrl())
	if err != nil {
		l.Errorf("set app address book error: %v", err)
		return &pb.GetAppAddressBookUrlResp{}, err
	}
	return &pb.GetAppAddressBookUrlResp{
		Url: url,
	}, nil
}
