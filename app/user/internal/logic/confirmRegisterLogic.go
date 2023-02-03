package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmRegisterLogic {
	return &ConfirmRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmRegisterLogic) ConfirmRegister(in *pb.ConfirmRegisterReq) (*pb.ConfirmRegisterResp, error) {
	userTmp := &usermodel.UserTmp{}
	// 使用id查询用户信息
	err := l.svcCtx.Mysql().Where("userId = ?", in.Id).Order("createdAt desc").First(userTmp).Error
	if err != nil {
		l.Errorf("ConfirmRegisterLogic ConfirmRegister err: %v", err)
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	// 用户存在 判断密码是否正确
	if !xpwd.VerifyPwd(in.Password, userTmp.Password, userTmp.PasswordSalt) {
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "登录失败"), l.svcCtx.T(in.CommonReq.Language, "密码错误"))}, nil
	}
	// 插入用户表
	user := &usermodel.User{
		Id:           userTmp.UserId,
		Password:     userTmp.Password,
		PasswordSalt: userTmp.PasswordSalt,
		Nickname:     l.svcCtx.SystemConfigMgr.Get("nickname_default"),
		Avatar:       utils.AnyRandomInSlice(l.svcCtx.SystemConfigMgr.GetSlice("avatars_default"), ""),
		RegInfo:      userTmp.RegInfo,
		CreateTime:   time.Now().UnixMilli(),
	}
	err = xorm.InsertOne(l.svcCtx.Mysql(), user)
	if err != nil {
		// id已被占用
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewAlertErrorResp(l.svcCtx.T(in.CommonReq.Language, "注册失败"), l.svcCtx.T(in.CommonReq.Language, "用户名已存在"))}, nil
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
		return &pb.ConfirmRegisterResp{CommonResp: resp.CommonResp}, err
	}
	return &pb.ConfirmRegisterResp{CommonResp: resp.CommonResp, Token: resp.Token, UserId: in.Id}, nil
}
