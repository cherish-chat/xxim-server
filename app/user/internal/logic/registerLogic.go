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
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
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

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// id是否被占用
	count, err := xorm.Count(l.svcCtx.Mysql(), &usermodel.User{Id: in.Id}, "id = ?", in.Id)
	if err != nil {
		l.Errorf("查询用户是否存在失败: %s", err.Error())
		return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if count > 0 {
		l.Errorf("用户已存在: %s", in.Id)
		return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "注册失败"), l.svcCtx.T(in.CommonReq.Language, "用户名已存在"))}, nil
	}
	// 插入用户表
	password := in.Password
	passwordSalt := utils.GenId()
	password = xpwd.GeneratePwd(password, passwordSalt)
	if in.Nickname == nil {
		in.Nickname = utils.AnyPtr(l.svcCtx.ConfigMgr.NicknameDefault(l.ctx))
	}
	ipRegion := ip2region.Ip2Region(in.CommonReq.Ip)
	user := &usermodel.User{
		Id:           in.Id,
		Password:     password,
		PasswordSalt: passwordSalt,
		Nickname:     *in.Nickname,
		Avatar:       utils.AnyRandomInSlice(l.svcCtx.ConfigMgr.AvatarsDefault(l.ctx), ""),
		RegInfo: &usermodel.LoginInfo{
			Time:        time.Now().UnixMilli(),
			Ip:          in.CommonReq.Ip,
			IpCountry:   ipRegion.Country,
			IpProvince:  ipRegion.Province,
			IpCity:      ipRegion.City,
			IpISP:       ipRegion.ISP,
			AppVersion:  in.CommonReq.AppVersion,
			UserAgent:   in.CommonReq.UserAgent,
			OsVersion:   in.CommonReq.OsVersion,
			Platform:    in.CommonReq.Platform,
			DeviceId:    in.CommonReq.DeviceId,
			DeviceModel: in.CommonReq.DeviceModel,
		},
		CreateTime: time.Now().UnixMilli(),
	}
	err = xorm.InsertOne(l.svcCtx.Mysql(), user)
	if err != nil {
		// id已被占用
		return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "注册失败"), l.svcCtx.T(in.CommonReq.Language, "用户名已存在"))}, nil
	} else {
		_ = usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{user.Id})
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "AfterRegister", func(ctx context.Context) {
			NewAfterLogic(ctx, l.svcCtx).AfterRegister(user.Id, in.CommonReq)
		}, propagation.MapCarrier{
			"user_id": user.Id,
		})
	}
	var resp *pb.LoginResp
	// 密码正确
	xtrace.StartFuncSpan(l.ctx, "login", func(ctx context.Context) {
		resp, err = NewLoginLogic(ctx, l.svcCtx).Login(&pb.LoginReq{
			CommonReq: in.CommonReq,
			Id:        in.Id,
			Password:  in.Password,
		})
	})
	if err != nil {
		l.Errorf("ConfirmRegisterLogic ConfirmRegister err: %v", err)
		return &pb.RegisterResp{CommonResp: resp.CommonResp}, err
	}
	return &pb.RegisterResp{CommonResp: resp.CommonResp, Token: resp.Token, UserId: in.Id}, nil
}
