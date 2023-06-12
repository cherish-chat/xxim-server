package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppAddressBookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppAddressBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppAddressBookLogic {
	return &UpdateAppAddressBookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppAddressBookLogic) UpdateAppAddressBook(in *pb.UpdateAppAddressBookReq) (*pb.UpdateAppAddressBookResp, error) {
	err := l.svcCtx.Redis().Set(rediskey.AppAddressBook(), in.AddressBook)
	if err != nil {
		l.Errorf("set app address book error: %v", err)
		return &pb.UpdateAppAddressBookResp{}, err
	}
	// 上传到对象存储
	url, err := NewUploadFileLogic(l.ctx, l.svcCtx).UploadFile(utils.GenId(), []byte(in.AddressBook))
	if err != nil {
		l.Errorf("upload file error: %v", err)
		return &pb.UpdateAppAddressBookResp{}, err
	}
	// 保存到数据库
	err = l.svcCtx.Redis().Set(rediskey.AppAddressBookUrl(), url)
	if err != nil {
		l.Errorf("set app address book url error: %v", err)
		return &pb.UpdateAppAddressBookResp{}, err
	}
	return &pb.UpdateAppAddressBookResp{}, nil
}
