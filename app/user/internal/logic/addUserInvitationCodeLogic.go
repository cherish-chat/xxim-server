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

type AddUserInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserInvitationCodeLogic {
	return &AddUserInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserInvitationCodeLogic) AddUserInvitationCode(in *pb.AddUserInvitationCodeReq) (*pb.AddUserInvitationCodeResp, error) {
	model := &usermodel.InvitationCode{
		Code:        in.UserInvitationCode.Code,
		Remark:      in.UserInvitationCode.Remark,
		Creator:     in.CommonReq.UserId,
		CreatorType: in.UserInvitationCode.CreatorType,
		IsEnable:    in.UserInvitationCode.IsEnable,
		CreateTime:  time.Now().UnixMilli(),
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
		return &pb.AddUserInvitationCodeResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddUserInvitationCodeResp{}, nil
}
