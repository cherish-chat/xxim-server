package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchUser 搜索用户
func (l *SearchUserLogic) SearchUser(in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	if in.Keyword == "" {
		return nil, nil
	}
	filter := bson.M{}
	filter["$or"] = []bson.M{
		{"_id": bson.M{"$regex": in.Keyword}},
		{"nickname": bson.M{"$regex": in.Keyword}},
	}
	var users []*dbmodel.User
	err := l.svcCtx.UserCollection().
		Find(l.ctx, filter).
		Skip(int64(in.PageSize * (in.Page - 1))).
		Limit(int64(in.PageSize)).
		Sort("nickname").
		All(&users)
	if err != nil {
		l.Errorf("SearchUser error: %v", err)
		return nil, err
	}
	var resp []*pb.UserData
	for _, user := range users {
		resp = append(resp, user.ToPbUser())
	}
	return &pb.SearchUserResp{
		UserDataList: resp,
	}, nil
}
