package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSUserLogic {
	return &AddMSUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSUserLogic) AddMSUser(in *pb.AddMSUserReq) (*pb.AddMSUserResp, error) {
	salt := utils.GenId()
	password := xpwd.GeneratePwd(in.User.Password, salt)
	// 插入用户表
	region := ip2region.Ip2Region(in.CommonReq.Ip)
	user := &mgmtmodel.User{
		Id:           in.User.Username,
		Password:     password,
		PasswordSalt: salt,
		Nickname:     in.User.Nickname,
		Avatar:       in.User.Avatar,
		RegInfo: &mgmtmodel.LoginInfo{
			Time:       time.Now().UnixMilli(),
			Ip:         in.CommonReq.Ip,
			IpCountry:  region.Country,
			IpProvince: region.Province,
			IpCity:     region.City,
			IpISP:      region.ISP,
			UserAgent:  in.CommonReq.UserAgent,
		},
		RoleId:     in.User.Role,
		CreateTime: time.Now().UnixMilli(),
		IsDisable:  in.User.IsDisable,
	}
	err := xorm.InsertOne(l.svcCtx.Mysql(), user)
	if err != nil {
		l.Errorf("AddMSUser err: %v", err)
		return &pb.AddMSUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AddMSUserResp{}, nil
}
