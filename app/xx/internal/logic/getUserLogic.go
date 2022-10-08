package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUser 获取用户信息
func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.GetUserResp, error) {
	if len(in.UserIdList) == 0 {
		return &pb.GetUserResp{}, nil
	}
	var users []*dbmodel.User
	err := l.svcCtx.UserCollection().Find(l.ctx, bson.M{
		"_id": bson.M{
			"$in": in.UserIdList,
		},
	}).All(&users)
	if err != nil {
		l.Errorf("GetUser error: %v", err)
		return nil, err
	}
	var resp []*pb.UserData
	for _, user := range users {
		resp = append(resp, user.ToPbUser())
	}
	return &pb.GetUserResp{
		UserDataList: resp,
	}, nil
}
