package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserModelLogic {
	return &AddUserModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserModelLogic) AddUserModel(in *pb.AddUserModelReq) (*pb.AddUserModelResp, error) {
	salt := utils.GenId()
	pwd := xpwd.GeneratePwd(in.Password, salt)
	now := time.Now()
	ip := in.CommonReq.Ip
	region := ip2region.Ip2Region(ip)
	model := &usermodel.User{
		Id:                strings.ToLower(in.UserModel.Id),
		Password:          pwd,
		PasswordSalt:      salt,
		InvitationCode:    in.UserModel.InvitationCode,
		Mobile:            in.UserModel.Mobile,
		MobileCountryCode: in.UserModel.MobileCountryCode,
		Nickname:          in.UserModel.Nickname,
		Avatar:            in.UserModel.Avatar,
		RegInfo: &usermodel.LoginInfo{
			Time:        now.UnixMilli(),
			Ip:          ip,
			IpCountry:   region.Country,
			IpProvince:  region.Province,
			IpCity:      region.City,
			IpISP:       region.ISP,
			AppVersion:  in.CommonReq.AppVersion,
			UserAgent:   in.CommonReq.UserAgent,
			OsVersion:   in.CommonReq.OsVersion,
			Platform:    in.CommonReq.Platform,
			DeviceId:    in.CommonReq.DeviceId,
			DeviceModel: in.CommonReq.DeviceModel,
		},
		Xb:       pb.XB(in.UserModel.Xb),
		Birthday: nil,
		InfoMap:  make(xorm.M),
		LevelInfo: usermodel.LevelInfo{
			Level:        0,
			Exp:          0,
			NextLevelExp: 0,
		},
		Role:          usermodel.Role(in.UserModel.Role),
		UnblockTime:   0,
		BlockRecordId: "",
		AdminRemark:   in.UserModel.AdminRemark,
		CreateTime:    time.Now().UnixMilli(),
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
		return &pb.AddUserModelResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "AfterRegister", func(ctx context.Context) {
		NewRegisterLogic(l.ctx, l.svcCtx).afterRegister(ctx, &pb.RegisterReq{
			CommonReq: in.CommonReq,
			Id:        in.UserModel.Id,
		}, model, nil)
	}, propagation.MapCarrier{
		"user_id": model.Id,
	})
	return &pb.AddUserModelResp{}, nil
}
