package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRobotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRobotLogic {
	return &CreateRobotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateRobot 创建机器人
func (l *CreateRobotLogic) CreateRobot(in *pb.CreateRobotReq) (*pb.CreateRobotResp, error) {
	user := &usermodel.User{
		UserId:       in.RobotId,
		RegisterTime: primitive.NewDateTimeFromTime(time.Now()),
		DestroyTime:  0,
		AccountMap: bson.M{
			pb.AccountTypeStatus: usermodel.AccountStatusNormal,
			pb.AccountTypeRole:   usermodel.AccountRoleRobot,
		},
		Nickname: "",
		Avatar:   "",
		ExtraMap: bson.M{
			pb.AccountExtraTypeRobotCreatedBy: in.Header.UserId,
		},
	}
	// 验证请求
	{
		if !l.svcCtx.Config.Account.Robot.AllowCreate {
			return &pb.CreateRobotResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "robot_not_allow_create")),
			}, nil
		}
		if in.Nickname == nil {
			// 是否允许空昵称
			if l.svcCtx.Config.Account.Robot.RequireNickname {
				return &pb.CreateRobotResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "robot_nickname_required")),
				}, nil
			} else {
				user.Nickname = l.svcCtx.Config.Account.Robot.DefaultNickname
			}
		} else {
			user.Nickname = *in.Nickname
		}
		if in.Avatar == nil {
			// 是否允许空头像
			if l.svcCtx.Config.Account.Robot.RequireAvatar {
				return &pb.CreateRobotResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "robot_avatar_required")),
				}, nil
			}
		} else {
			user.Avatar = *in.Avatar
		}

		// userId是否已存在
		{
			count, err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{
				"userId": user.UserId,
			}).Count()
			if err != nil {
				l.Errorf("Find err: %v", err)
				return nil, err
			}
			if count > 0 {
				return &pb.CreateRobotResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.Get(in.Header.Language, "robot_id_exist")),
				}, nil
			}
		}
	}

	_, err := l.svcCtx.UserCollection.InsertOne(l.ctx, user)
	if err != nil {
		l.Errorf("InsertOne err: %v", err)
		return nil, err
	}

	l.Debugf("create robot success: %v", user)

	userAccessTokenResp := NewUserAccessTokenLogic(l.ctx, l.svcCtx).generateToken(&pb.UserAccessTokenReq{
		Header: in.Header,
	}, user)

	return &pb.CreateRobotResp{
		AccessToken: userAccessTokenResp.AccessToken,
	}, nil
}
