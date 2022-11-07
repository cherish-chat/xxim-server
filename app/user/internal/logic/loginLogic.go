package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user := &usermodel.User{}
	// 使用id查询用户信息
	err := l.svcCtx.Mongo().Collection(&usermodel.User{}).Find(l.ctx, bson.M{
		"_id": in.Id,
	}).One(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 用户不存在 注册流程
			return l.register(in)
		} else {
			// 报错
			return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	// 用户存在 判断密码是否正确
	if !xpwd.VerifyPwd(in.Password, user.Password, user.PasswordSalt) {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp("登录失败", "密码错误")}, nil
	}
	// 密码正确
	// 生成token
	//uniqueSuffix := fmt.Sprintf("%s:%s", in.Requester.Platform, in.Requester.DeviceId) // 如果你不限制同设备登录多次一个账号，可以使用这行代码
	uniqueSuffix := fmt.Sprintf("%s", in.Requester.Platform)
	tokenObj := xjwt.GenerateToken(user.Id, uniqueSuffix,
		xjwt.WithPlatform(in.Requester.Platform),
		xjwt.WithDeviceId(in.Requester.DeviceId),
		xjwt.WithDeviceModel(in.Requester.DeviceModel),
	)
	// 断开设备连接
	_, err = l.svcCtx.ImService().KickUserConn(l.ctx, &pb.KickUserConnReq{GetUserConnReq: &pb.GetUserConnReq{
		UserIds:   []string{user.Id},
		Platforms: []string{in.Requester.Platform},
	}})
	if err != nil {
		l.Errorf("kick user conn failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	err = xjwt.SaveToken(l.ctx, l.svcCtx.Redis(), tokenObj)
	if err != nil {
		l.Errorf("save token failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SaveLoginRecord", func(ctx context.Context) {
		region := ip2region.Ip2Region(in.Requester.Ip)
		record := &usermodel.LoginRecord{
			Id:     utils.GenId(),
			UserId: user.Id,
			LoginInfo: usermodel.LoginInfo{
				Time:        time.Now().UnixMilli(),
				Ip:          in.Requester.Ip,
				IpCountry:   region.Country,
				IpProvince:  region.Province,
				IpCity:      region.City,
				IpISP:       region.ISP,
				AppVersion:  in.Requester.AppVersion,
				Ua:          in.Requester.Ua,
				OsVersion:   in.Requester.OsVersion,
				Platform:    in.Requester.Platform,
				DeviceId:    in.Requester.DeviceId,
				DeviceModel: in.Requester.DeviceModel,
			},
		}
		_, err := l.svcCtx.Mongo().Collection(&usermodel.LoginRecord{}).InsertOne(ctx, record)
		if err != nil {
			l.Errorf("save login record failed, err: %v", err)
		}
	}, nil)
	return &pb.LoginResp{
		IsNewUser: false,
		Token:     tokenObj.Token,
	}, nil
}

func (l *LoginLogic) register(in *pb.LoginReq) (*pb.LoginResp, error) {
	// 检查用户id是否符合规则 只能是字母数字下划线
	reg := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !reg.MatchString(in.Id) {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp("注册失败", "id只能包含字母数字下划线")}, nil
	}
	// 检查用户id是否符合规则 长度在6-20之间
	if len(in.Id) < 6 || len(in.Id) > 20 {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp("注册失败", "id长度在6-20之间")}, nil
	}
	region := ip2region.Ip2Region(in.Requester.Ip)
	// 注册
	salt := utils.GenId()
	userTmp := &usermodel.UserTmp{
		UserId:       in.Id,
		Password:     xpwd.GeneratePwd(in.Password, salt),
		PasswordSalt: salt,
		RegInfo: &usermodel.LoginInfo{
			Time:        time.Now().UnixMilli(),
			Ip:          in.Requester.Ip,
			IpCountry:   region.Country,
			IpProvince:  region.Province,
			IpCity:      region.City,
			IpISP:       region.ISP,
			AppVersion:  in.Requester.AppVersion,
			Ua:          in.Requester.Ua,
			OsVersion:   in.Requester.OsVersion,
			Platform:    in.Requester.Platform,
			DeviceId:    in.Requester.DeviceId,
			DeviceModel: in.Requester.DeviceModel,
		},
	}
	// 保存用户信息
	_, err := l.svcCtx.Mongo().Collection(&usermodel.UserTmp{}).InsertOne(l.ctx, userTmp)
	if err != nil {
		l.Errorf("insert user failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.LoginResp{
		CommonResp: pb.NewSuccessResp(),
		IsNewUser:  true,
		Token:      "",
	}, nil
}
