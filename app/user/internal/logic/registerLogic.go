package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
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
	// ip频率限制
	period, quota := l.svcCtx.ConfigMgr.RegisterIpLimit(l.ctx)
	if period == 0 || quota == 0 {
		// 不限制
	} else {
		redisKey := rediskey.RegisterIpLimitKey(in.CommonReq.Ip)
		// incrby 1
		val, err := l.svcCtx.Redis().IncrbyCtx(l.ctx, redisKey, 1)
		if err != nil {
			l.Errorf("incrby register ip limit key err: %v", err)
			return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		// 如果是第一次设置过期时间
		if val == 1 {
			err = l.svcCtx.Redis().ExpireCtx(l.ctx, redisKey, period)
			if err != nil {
				l.Errorf("expire register ip limit key err: %v", err)
				return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
		}
		// 判断是否超过限制
		if val > int64(quota) {
			l.Errorf("register ip limit key over quota: %v", val)
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册频率过高，请稍后再试")}, nil
		}
	}
	// 是否开启ip白名单限制
	if l.svcCtx.ConfigMgr.RegisterMustIpInWhiteList(l.ctx) {
		var count int64
		err := l.svcCtx.Mysql().Model(&usermodel.IpWhiteList{}).
			Where("isEnable = ?", true).
			Where("startIp <= ? and endIp >= ?", in.CommonReq.Ip, in.CommonReq.Ip).
			Count(&count).Error
		if err != nil {
			l.Errorf("check ip in whitelist err: %v", err)
			return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if count == 0 {
			l.Errorf("ip not in whitelist: %v", in.CommonReq.Ip)
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，ip不在白名单中")}, nil
		}
	}
	// 注册开关
	if !l.svcCtx.ConfigMgr.EnablePlatformRegister(l.ctx, in.CommonReq.Platform) {
		return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "不允许在"+in.CommonReq.Platform+"上注册")}, nil
	}
	var inviCode = ""
	var invitationCode *usermodel.InvitationCode
	// 是否必须要邀请码
	if l.svcCtx.ConfigMgr.RegisterMustInviteCode(l.ctx) {
		// 请求中必须带邀请码
		if in.InvitationCode == nil {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，必须要邀请码")}, nil
		}
		if *in.InvitationCode == "" {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，邀请码不能为空")}, nil
		}
	}
	if in.InvitationCode != nil && *in.InvitationCode != "" {
		inviCode = *in.InvitationCode
		// 邀请码是否存在
		var ic = &usermodel.InvitationCode{}
		err := l.svcCtx.Mysql().Model(&usermodel.InvitationCode{}).Where("code = ?", inviCode).First(ic).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，邀请码不存在")}, nil
			}
			l.Errorf("check invitation code err: %v", err)
			return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if !ic.IsEnable {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，邀请码已失效")}, nil
		}
		invitationCode = ic
	}
	var mobile = ""
	var mobileCountryCode = ""
	// 是否必须要手机号
	if l.svcCtx.ConfigMgr.RegisterMustMobile(l.ctx) {
		// 请求中必须带手机号
		if in.Mobile == nil || in.MobileCountryCode == nil {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，手机号为空")}, nil
		}
		if *in.Mobile == "" || *in.MobileCountryCode == "" {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，手机号为空")}, nil
		}
		mobile = *in.Mobile
		mobileCountryCode = *in.MobileCountryCode
		// 手机号是否存在
		var count int64
		err := l.svcCtx.Mysql().Model(&usermodel.User{}).Where("mobile = ? and mobileCountryCode = ?", mobile, mobileCountryCode).Count(&count).Error
		if err != nil {
			l.Errorf("check mobile err: %v", err)
			return &pb.RegisterResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if count > 0 {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，手机号已存在")}, nil
		}
	}
	// 是否必须要smsCode
	if l.svcCtx.ConfigMgr.RegisterMustSmsCode(l.ctx) {
		// 请求中必须带smsCode
		if in.SmsCode == nil {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，短信验证码为空")}, nil
		}
		if *in.SmsCode == "" {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，短信验证码为空")}, nil
		}
		//TODO: 验证smsCode
	}
	// 是否必填头像
	var avatar = utils.AnyRandomInSlice(l.svcCtx.ConfigMgr.AvatarsDefault(l.ctx), "")
	if l.svcCtx.ConfigMgr.RegisterMustAvatar(l.ctx) {
		// 请求中必须带头像
		if in.Avatar == nil {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，头像为空")}, nil
		}
		if *in.Avatar == "" {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，头像为空")}, nil
		}
		avatar = *in.Avatar
	}
	// 是否必填昵称
	var nickname = l.svcCtx.ConfigMgr.NicknameDefault(l.ctx)
	if l.svcCtx.ConfigMgr.RegisterMustNickname(l.ctx) {
		// 请求中必须带昵称
		if in.Nickname == nil {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，昵称为空")}, nil
		}
		if *in.Nickname == "" {
			return &pb.RegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "注册失败，昵称为空")}, nil
		}
		nickname = *in.Nickname
	}

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
	ipRegion := ip2region.Ip2Region(in.CommonReq.Ip)
	user := &usermodel.User{
		Id:           in.Id,
		Password:     password,
		PasswordSalt: passwordSalt,
		Nickname:     nickname,
		Avatar:       avatar,
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
		InvitationCode: inviCode,
		CreateTime:     time.Now().UnixMilli(),
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
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "AfterRegister", func(ctx context.Context) {
		l.afterRegister(ctx, in, user, invitationCode)
	}, propagation.MapCarrier{
		"user_id": user.Id,
	})
	return &pb.RegisterResp{CommonResp: resp.CommonResp, Token: resp.Token, UserId: in.Id}, nil
}

func (l *RegisterLogic) afterRegister(ctx context.Context, in *pb.RegisterReq, user *usermodel.User, invitationCode *usermodel.InvitationCode) {
	// 查询预设会话
	var defaultConvs []*usermodel.DefaultConv
	var tx = l.svcCtx.Mysql().Model(&usermodel.DefaultConv{})
	if user.InvitationCode != "" {
		tx = tx.Where("((filterType = ? AND invitationCode = ?) OR filterType = ?)", 1, user.InvitationCode, 0)
	} else {
		tx = tx.Where("filterType = ?", 0)
	}
	if invitationCode != nil {
		if invitationCode.DefaultConvMode == 0 {
			// 添加所有预设会话
		} else if invitationCode.DefaultConvMode == 1 {
			// 1:只添加一个会话(轮询)
			// 查询上次该邀请码轮询到的id
			val, _ := l.svcCtx.Redis().GetCtx(ctx, rediskey.LatestTurnConvIdKey(invitationCode.Code))
			// 这个id是不是目前最大的
			var count int64
			l.svcCtx.Mysql().Model(&usermodel.DefaultConv{}).Where("id > ?", val).Count(&count)
			if count == 0 {
				// 是的话,就从头开始
			} else {
				// 不是的话,就从下一个开始
				tx = tx.Where("id > ?", val)
			}
		} else if invitationCode.DefaultConvMode == 2 {
			// 2:只添加一个会话(随机)
			val, _ := l.svcCtx.Redis().GetCtx(ctx, rediskey.LatestTurnConvIdKey(invitationCode.Code))
			// 这个id是不是目前最大的
			var count int64
			l.svcCtx.Mysql().Model(&usermodel.DefaultConv{}).Where("id > ?", val).Count(&count)
			if count == 0 {
				// 是的话,就从头开始
			} else {
				// 不是的话,就从下一个开始
				tx = tx.Where("id > ?", val)
			}
		} else {
			// 不添加
			return
		}
	}
	err := tx.Order("id asc").Limit(100).Find(&defaultConvs).Error
	if err != nil {
		l.Errorf("查询预设会话失败: %s", err.Error())
		return
	}
	if len(defaultConvs) == 0 {
		return
	}
	var friendCount int64
	var successFriendId string
	for _, conv := range defaultConvs {
		if conv.ConvType == 0 {
			// 好友
			if invitationCode != nil && invitationCode.DefaultConvMode == 1 {
				// 只添加一个会话
				if friendCount > 0 {
					// 更新 latestTurnConvId
					l.svcCtx.Redis().SetCtx(ctx, rediskey.LatestTurnConvIdKey(invitationCode.Code), successFriendId)
					continue
				}
			}
			_, err := l.svcCtx.RelationService().AcceptAddFriend(ctx, &pb.AcceptAddFriendReq{
				CommonReq: &pb.CommonReq{
					UserId: conv.ConvId,
				},
				ApplyUserId: user.Id,
				RequestId:   nil,
				SendTextMsg: utils.AnyPtr(conv.TextMsg),
			})
			if err != nil {
				l.Errorf("自动添加好友失败: %s", err.Error())
			} else {
				friendCount++
				successFriendId = conv.ConvId
			}
		} else if conv.ConvType == 1 {
			_, err := l.svcCtx.GroupService().AddGroupMember(ctx, &pb.AddGroupMemberReq{
				CommonReq: in.CommonReq,
				GroupId:   conv.ConvId,
				UserId:    user.Id,
			})
			if err != nil {
				l.Errorf("自动加入群组失败: %s", err.Error())
			}
		}
	}
}
