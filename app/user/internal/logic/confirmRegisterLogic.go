package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"

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
	err := l.svcCtx.Mongo().Collection(&usermodel.UserTmp{}).Find(l.ctx, bson.M{
		"userId": in.Id,
	}).One(userTmp)
	if err != nil {
		l.Errorf("ConfirmRegisterLogic ConfirmRegister err: %v", err)
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	// 用户存在 判断密码是否正确
	if !xpwd.VerifyPwd(in.Password, userTmp.Password, userTmp.PasswordSalt) {
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewAlertErrorResp("登录失败", "密码错误")}, nil
	}
	// 插入用户表
	user := &usermodel.User{
		Id:           userTmp.UserId,
		Password:     userTmp.Password,
		PasswordSalt: userTmp.PasswordSalt,
		Nickname:     "", // TODO random nickname
		Avatar:       "", // TODO random avatar
		RegInfo:      userTmp.RegInfo,
	}
	_, err = l.svcCtx.Mongo().Collection(&usermodel.User{}).InsertOne(l.ctx, user)
	if err != nil {
		// id已被占用
		return &pb.ConfirmRegisterResp{CommonResp: pb.NewAlertErrorResp("注册失败", "id已被占用")}, nil
	}
	var resp *pb.LoginResp
	// 密码正确
	xtrace.StartFuncSpan(l.ctx, "login", func(ctx context.Context) {
		resp, err = NewLoginLogic(ctx, l.svcCtx).Login(&pb.LoginReq{
			Requester: in.Requester,
			Id:        in.Id,
			Password:  in.Password,
		})
	})
	if err != nil {
		l.Errorf("ConfirmRegisterLogic ConfirmRegister err: %v", err)
		return &pb.ConfirmRegisterResp{CommonResp: resp.CommonResp}, err
	}
	return &pb.ConfirmRegisterResp{CommonResp: resp.CommonResp, Token: resp.Token}, nil
}
