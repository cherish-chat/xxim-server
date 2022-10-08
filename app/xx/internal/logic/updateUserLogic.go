package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"github.com/cherish-chat/xxim-server/common/utils/pwd"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUser 更新用户信息
func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	userLogic := NewGetUserLogic(l.ctx, l.svcCtx)
	getUser, err := userLogic.GetUser(&pb.GetUserReq{UserIdList: []string{in.Base.SelfId}})
	if err != nil {
		l.Errorf("GetUser error: %v", err)
		return nil, err
	}
	if len(getUser.UserDataList) == 0 {
		l.Errorf("GetUser error: user not found")
		return nil, err
	}
	user := getUser.UserDataList[0]
	updateMap := bson.M{}
	// 对比要更新的字段
	if in.UserData.Nickname != "" {
		if in.UserData.Nickname != user.Nickname {
			updateMap["nickname"] = in.UserData.Nickname
		}
	}
	if in.UserData.Avatar != "" {
		if in.UserData.Avatar != user.Avatar {
			updateMap["avatar"] = in.UserData.Avatar
		}
	}
	if in.UserData.Xb != "" {
		if in.UserData.Xb != user.Xb {
			updateMap["xb"] = in.UserData.Xb
		}
	}
	if in.UserData.Birthday != "" {
		if in.UserData.Birthday != user.Birthday {
			updateMap["birthday"] = in.UserData.Birthday
		}
	}
	if in.UserData.Signature != "" {
		if in.UserData.Signature != user.Signature {
			updateMap["signature"] = in.UserData.Signature
		}
	}
	if in.UserData.Tags != nil {
		if len(in.UserData.Tags) != len(user.Tags) {
			updateMap["tags"] = in.UserData.Tags
		} else {
			for i, v := range in.UserData.Tags {
				if v != user.Tags[i] {
					updateMap["tags"] = in.UserData.Tags
					break
				}
			}
		}
	}
	if in.UserData.Password != "" {
		if !pwd.VerifyPwd(in.UserData.Password, user.Password, []byte(user.RegisterInfo.Salt)) {
			updateMap["password"] = pwd.GeneratePwd(in.UserData.Password, []byte(user.RegisterInfo.Salt))
		}
	}
	if in.UserData.Ex != nil && len(in.UserData.Ex) != 0 {
		updateMap["ex"] = dbmodel.NewUserEx(in.UserData.Ex)
	}
	if len(updateMap) == 0 {
		return &pb.UpdateUserResp{}, nil
	}
	err = l.svcCtx.UserCollection().UpdateOne(l.ctx, bson.M{"_id": in.Base.SelfId}, bson.M{"$set": updateMap})
	if err != nil {
		l.Errorf("UpdateUser error: %v", err)
		return nil, err
	}
	getUser, err = userLogic.GetUser(&pb.GetUserReq{UserIdList: []string{in.Base.SelfId}})
	if err != nil {
		l.Errorf("GetUser error: %v", err)
		return nil, err
	}
	if len(getUser.UserDataList) == 0 {
		l.Errorf("GetUser error: user not found")
		return nil, err
	}
	user = getUser.UserDataList[0]
	return &pb.UpdateUserResp{
		UserData: user,
	}, nil
}
