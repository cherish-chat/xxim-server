package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
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
	err := xorm.DetailByWhere(l.svcCtx.Mysql(), user, xorm.Where("id = ?", in.Id))
	if err != nil {
		if xorm.RecordNotFound(err) {
			// 用户不存在 注册流程
			return l.register(in)
		} else {
			// 报错
			return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	// 用户密码输入错误次数是否超过限制
	// 获取用户密码输入错误次数
	val, _ := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.UserPasswordErrorCountKey(in.Id))
	if utils.AnyToInt64(val) > l.svcCtx.ConfigMgr.UserPasswordErrorMaxCount(l.ctx, in.CommonReq.UserId) {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp("登录失败", "登录失败，密码输入错误次数超过限制")}, nil
	}
	// 用户存在 判断密码是否正确
	if !xpwd.VerifyPwd(in.Password, user.Password, user.PasswordSalt) {
		// 记录用户密码输入错误次数
		_, _ = l.svcCtx.Redis().IncrbyCtx(l.ctx, rediskey.UserPasswordErrorCountKey(in.Id), 1)
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "密码错误"))}, nil
	}
	// 判断用户角色
	{
		if user.Role == usermodel.RoleUser {
			// 用户能否在该平台功能
			if !l.svcCtx.ConfigMgr.LoginUserOnPlatform(l.ctx, in.CommonReq.Platform) {
				return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "用户不能在"+in.CommonReq.Platform+"登录"))}, nil
			}
			// 用户登录时是否需要ip在白名单
			if l.svcCtx.ConfigMgr.LoginUserNeedIpWhiteList(l.ctx, user.Id) {
				resp, err := l.checkIpWhiteList(in)
				if err != nil {
					return resp, err
				}
				if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
					return resp, nil
				}
			}
		} else if user.Role == usermodel.RoleService {
			// 客服能否在该平台功能
			if !l.svcCtx.ConfigMgr.LoginServiceOnPlatform(l.ctx, in.CommonReq.Platform, user.Id) {
				return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "客服不能在"+in.CommonReq.Platform+"登录"))}, nil
			}
			// 客服登录时是否需要ip在白名单
			if l.svcCtx.ConfigMgr.LoginServiceNeedIpWhiteList(l.ctx, user.Id) {
				resp, err := l.checkIpWhiteList(in)
				if err != nil {
					return resp, err
				}
				if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
					return resp, nil
				}
			}
		} else if user.Role == usermodel.RoleGuest {
			// 游客能否在该平台功能
			if !l.svcCtx.ConfigMgr.LoginGuestOnPlatform(l.ctx, in.CommonReq.Platform, user.Id) {
				return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "游客不能在"+in.CommonReq.Platform+"登录"))}, nil
			}
			// 游客登录时是否需要ip在白名单
			if l.svcCtx.ConfigMgr.LoginGuestNeedIpWhiteList(l.ctx, user.Id) {
				resp, err := l.checkIpWhiteList(in)
				if err != nil {
					return resp, err
				}
				if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
					return resp, nil
				}
			}
		} else {
			// 未知角色
			return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "账号角色异常"))}, nil
		}
	}
	// 密码正确
	if user.DestroyTime > 0 {
		return &pb.LoginResp{
			IsNewUser:   false,
			Token:       "",
			UserId:      user.Id,
			IsDestroyed: true,
		}, nil
	}
	// 生成token
	// 是否允许同平台多设备登录
	uniqueSuffix := fmt.Sprintf("%s", in.CommonReq.Platform)
	if l.svcCtx.Config.EnableMultiDeviceLogin {
		uniqueSuffix = fmt.Sprintf("%s:%s", in.CommonReq.Platform, in.CommonReq.DeviceId)
	}
	tokenObj := xjwt.GenerateToken(user.Id, uniqueSuffix,
		xjwt.WithPlatform(in.CommonReq.Platform),
		xjwt.WithDeviceId(in.CommonReq.DeviceId),
		xjwt.WithDeviceModel(in.CommonReq.DeviceModel),
	)
	// 断开设备连接
	getUserConnReq := &pb.GetUserConnReq{
		UserIds:   []string{user.Id},
		Platforms: []string{in.CommonReq.Platform},
	}
	if l.svcCtx.Config.EnableMultiDeviceLogin {
		getUserConnReq.Devices = []string{in.CommonReq.DeviceId}
	}
	_, err = l.svcCtx.ImService().KickUserConn(l.ctx, &pb.KickUserConnReq{GetUserConnReq: getUserConnReq})
	if err != nil {
		l.Errorf("kick user conn failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	err = xjwt.SaveToken(l.ctx, l.svcCtx.Redis(), tokenObj)
	if err != nil {
		l.Errorf("save token failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "AfterLogin", func(ctx context.Context) {
		NewAfterLogic(ctx, l.svcCtx).AfterLogin(user.Id, in.CommonReq)
	}, propagation.MapCarrier{
		"user_id": user.Id,
	})
	return &pb.LoginResp{
		IsNewUser: false,
		Token:     tokenObj.Token,
		UserId:    user.Id,
	}, nil
}

func (l *LoginLogic) register(in *pb.LoginReq) (*pb.LoginResp, error) {
	// 检查用户id是否符合规则 只能是字母数字下划线
	reg := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !reg.MatchString(in.Id) {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, l.svcCtx.T(in.CommonReq.Language, "注册失败")), l.svcCtx.T(in.CommonReq.Language, "用户名违规"))}, nil
	}
	// 检查用户id是否符合规则 长度在6-20之间
	if len(in.Id) < 6 || len(in.Id) > 20 {
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "注册失败"), l.svcCtx.T(in.CommonReq.Language, "用户名违规"))}, nil
	}
	region := ip2region.Ip2Region(in.CommonReq.Ip)
	// 注册
	salt := utils.GenId()
	userTmp := &usermodel.UserTmp{
		UserId:       in.Id,
		Password:     xpwd.GeneratePwd(in.Password, salt),
		PasswordSalt: salt,
		RegInfo: &usermodel.LoginInfo{
			Time:        time.Now().UnixMilli(),
			Ip:          in.CommonReq.Ip,
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
	}
	// 保存用户信息
	err := xorm.InsertOne(l.svcCtx.Mysql(), userTmp)
	if err != nil {
		l.Errorf("insert user failed, err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.LoginResp{
		CommonResp: pb.NewSuccessResp(),
		IsNewUser:  true,
		Token:      "",
		UserId:     in.Id,
	}, nil
}

func (l *LoginLogic) checkIpWhiteList(in *pb.LoginReq) (*pb.LoginResp, error) {
	var count int64
	err := l.svcCtx.Mysql().Model(&usermodel.IpWhiteList{}).
		Where("isEnable = ?", true).
		Where("startIp <= ? and endIp >= ?", in.CommonReq.Ip, in.CommonReq.Ip).
		Count(&count).Error
	if err != nil {
		l.Errorf("check ip in whitelist err: %v", err)
		return &pb.LoginResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if count == 0 {
		l.Errorf("ip not in whitelist: %v", in.CommonReq.Ip)
		return &pb.LoginResp{CommonResp: pb.NewAlertErrorResp("登录失败", "登录失败，ip不在白名单中")}, nil
	}
	return nil, nil
}
