package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppAddressBookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppAddressBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppAddressBookLogic {
	return &GetAppAddressBookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppAddressBookLogic) GetAppAddressBook(in *pb.GetAppAddressBookReq) (*pb.GetAppAddressBookResp, error) {
	addressBook, err := l.svcCtx.Redis().Get(rediskey.AppAddressBook())
	if err != nil {
		l.Errorf("set app address book error: %v", err)
		return &pb.GetAppAddressBookResp{}, err
	}
	return &pb.GetAppAddressBookResp{
		AddressBook: addressBook,
	}, nil
}
