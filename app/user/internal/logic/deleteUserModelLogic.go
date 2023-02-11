package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserModelLogic {
	return &DeleteUserModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserModelLogic) DeleteUserModel(in *pb.DeleteUserModelReq) (*pb.DeleteUserModelResp, error) {
	// 查询models
	models := make([]*usermodel.User, 0)
	err := l.svcCtx.Mysql().Model(&usermodel.User{}).Where("id in (?)", in.Ids).Find(&models).Error
	if err != nil {
		l.Errorf("find error: %v", err)
		return &pb.DeleteUserModelResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	var userRecycleBins []*usermodel.UserRecycleBin
	for _, model := range models {
		userRecycleBin := &usermodel.UserRecycleBin{
			Id:         utils.GenId(),
			UserModel:  model,
			CreateTime: time.Now().UnixMilli(),
			Creator:    in.CommonReq.UserId,
			Ip:         in.CommonReq.Ip,
			IpRegion:   ip2region.Ip2Region(in.CommonReq.Ip).String(),
		}
		userRecycleBins = append(userRecycleBins, userRecycleBin)
	}
	model := &usermodel.User{}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := tx.Model(model).Where("id in (?)", in.Ids).Delete(model).Error
		if err != nil {
			return err
		}
		// 存储到回收站
		err = tx.Model(&usermodel.UserRecycleBin{}).Create(userRecycleBins).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteUserModelResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteUserModelResp{}, nil
}
