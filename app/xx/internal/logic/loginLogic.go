package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/jwt"
	"github.com/cherish-chat/xxim-server/common/utils/pwd"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
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

// Login 登录
func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user := &dbmodel.User{}
	if in.UserId == "" && in.Base == nil && in.Base.DeviceId == "" {
		return &pb.LoginResp{FailedReason: "参数错误"}, nil
	}
	if in.UserId == "" {
		// 游客登录
		err := l.svcCtx.UserCollection().Find(l.ctx, bson.D{
			{"isGuest", true},
			{"registryInfo.deviceId", in.Base.DeviceId},
		}).One(user)
		if err != nil {
			// 创建游客
			genId := utils.GenId()
			_, err = NewRegisterLogic(l.ctx, l.svcCtx).Register(&pb.RegisterReq{
				Base: in.Base,
				UserData: &pb.UserData{
					Id:           genId,
					Nickname:     "游客" + utils.RandString(6),
					Avatar:       "https://cdn.xx.com/xx.png",
					Xb:           "",
					Birthday:     "",
					Signature:    "",
					Tags:         nil,
					Password:     "",
					RegisterInfo: nil,
					IsRobot:      false,
					IsGuest:      true,
					IsAdmin:      false,
					IsOfficial:   false,
					UnbanTime:    "",
					AdminRemark:  "",
					Ex:           nil,
				},
			})
			if err != nil {
				l.Errorf("游客登录失败: %v", err)
				return &pb.LoginResp{FailedReason: "登录失败"}, nil
			}
			user.Id = genId
		}
	} else {
		// 用户登录
		err := l.svcCtx.UserCollection().Find(l.ctx, bson.D{
			{"_id", in.UserId},
		}).One(user)
		if err != nil {
			l.Errorf("用户id查询失败: %v", err)
			return &pb.LoginResp{FailedReason: "登录失败"}, nil
		}
		// 密码校验
		if !pwd.VerifyPwd(in.Password, user.Password, []byte(user.RegistryInfo.Salt)) {
			return &pb.LoginResp{FailedReason: "密码错误"}, nil
		}
	}
	// 生成token
	token := jwt.GenerateTokenButNotSet(user.Id)
	getUser, err := NewGetUserLogic(l.ctx, l.svcCtx).GetUser(&pb.GetUserReq{
		UserIdList: []string{user.Id},
	})
	if err != nil {
		l.Errorf("获取用户信息失败: %v", err)
		return &pb.LoginResp{FailedReason: "登录失败"}, nil
	}
	if len(getUser.UserDataList) == 0 {
		return &pb.LoginResp{FailedReason: "登录失败"}, nil
	}
	// 替换token
	err = jwt.SetToken(l.ctx, l.svcCtx.Redis(), token, in.Base.Platform)
	if err != nil {
		l.Errorf("token替换失败: %v", err)
		return &pb.LoginResp{FailedReason: err.Error()}, nil
	}
	return &pb.LoginResp{
		UserData: getUser.UserDataList[0],
		Token:    token,
	}, nil
}
