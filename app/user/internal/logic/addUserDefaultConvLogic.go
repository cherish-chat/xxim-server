package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserDefaultConvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserDefaultConvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserDefaultConvLogic {
	return &AddUserDefaultConvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserDefaultConvLogic) AddUserDefaultConv(in *pb.AddUserDefaultConvReq) (*pb.AddUserDefaultConvResp, error) {
	model := &usermodel.DefaultConv{
		Id:             usermodel.GetId(l.svcCtx.Mysql(), &usermodel.DefaultConv{}, 10000),
		ConvType:       in.UserDefaultConv.ConvType,
		FilterType:     in.UserDefaultConv.FilterType,
		InvitationCode: in.UserDefaultConv.InvitationCode,
		ConvId:         in.UserDefaultConv.ConvId,
		TextMsg:        in.UserDefaultConv.TextMsg,
		CreateTime:     time.Now().UnixMilli(),
	}
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := model.Insert(tx)
		if err != nil {
			l.Errorf("insert err: %v", err)
		}
		return err
	})
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddUserDefaultConvResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddUserDefaultConvResp{}, nil
}
