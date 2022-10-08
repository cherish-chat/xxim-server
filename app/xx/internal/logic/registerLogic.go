package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/utils/pwd"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 注册用户
func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	region := ip2region.RegionSplit(in.Base.Ips)
	salt := utils.GenId()
	model := &dbmodel.User{
		Id:          in.UserData.Id,
		Nickname:    in.UserData.Nickname,
		Avatar:      in.UserData.Avatar,
		Xb:          in.UserData.Xb,
		Birthday:    in.UserData.Birthday,
		Signature:   in.UserData.Signature,
		Tags:        in.UserData.Tags,
		Password:    pwd.GeneratePwd(in.UserData.Password, []byte(salt)),
		IsRobot:     in.UserData.IsRobot,
		IsGuest:     in.UserData.IsGuest,
		IsAdmin:     in.UserData.IsAdmin,
		IsOfficial:  in.UserData.IsOfficial,
		UnbanTime:   in.UserData.UnbanTime,
		AdminRemark: in.UserData.AdminRemark,
		RegistryInfo: dbmodel.RegistryInfo{
			Platform:       in.Base.Platform,
			Time:           time.Now().Format("2006-01-02 15:04:05"),
			DeviceModel:    in.Base.DeviceModel,
			DeviceId:       in.Base.DeviceId,
			Ips:            strings.Split(in.Base.Ips, ","),
			RegisterSource: "",
			Salt:           salt,
			IpCountry:      region.Country,
			IpProvince:     region.Province,
			IpCity:         region.City,
			IpDistrict:     region.District,
		},
		Ex: dbmodel.NewUserEx(in.UserData.Ex),
	}
	err := l.svcCtx.UserCollection().Find(l.ctx, bson.M{
		"_id": in.UserData.Id,
	}).One(&dbmodel.User{})
	if err == nil {
		return &pb.RegisterResp{
			FailedReason: "id已存在",
		}, nil
	}
	_, err = l.svcCtx.UserCollection().InsertOne(l.ctx, model)
	if err != nil {
		return &pb.RegisterResp{
			FailedReason: "注册失败",
		}, nil
	}
	return &pb.RegisterResp{}, nil
}
