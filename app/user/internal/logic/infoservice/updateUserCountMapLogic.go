package infoservicelogic

import (
	"context"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserCountMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserCountMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserCountMapLogic {
	return &UpdateUserCountMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserCountMap 更新用户计数信息
func (l *UpdateUserCountMapLogic) UpdateUserCountMap(in *pb.UpdateUserCountMapReq) (*pb.UpdateUserCountMapResp, error) {
	if in.Statistics {
		return l.statistic(in)
	}
	var bsonM bson.M
	switch in.Algorithm {
	case pb.UpdateUserCountMapReq_add:
		bsonM = bson.M{
			"$inc": bson.M{
				"countMap." + in.CountType.String(): in.Count,
			},
		}
	case pb.UpdateUserCountMapReq_sub:
		bsonM = bson.M{
			"$inc": bson.M{
				"countMap." + in.CountType.String(): -in.Count,
			},
		}
	case pb.UpdateUserCountMapReq_fixed:
		bsonM = bson.M{
			"$set": bson.M{
				"countMap." + in.CountType.String(): in.Count,
			},
		}
	default:
		return nil, nil
	}
	err := l.svcCtx.UserCollection.UpdateOne(l.ctx, bson.M{
		"userId": in.Header.UserId,
	}, bsonM)
	if err != nil {
		l.Errorf("update user count map error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	return &pb.UpdateUserCountMapResp{}, nil
}

func (l *UpdateUserCountMapLogic) statistic(in *pb.UpdateUserCountMapReq) (*pb.UpdateUserCountMapResp, error) {
	switch in.CountType {
	case pb.UpdateUserCountMapReq_friendCount:
		return l.statisticFriendCount(in)
	case pb.UpdateUserCountMapReq_joinGroupCount:
		return l.statisticJoinGroupCount(in)
	case pb.UpdateUserCountMapReq_createGroupCount:
		return l.statisticCreateGroupCount(in)
	default:
		return nil, nil
	}
}

// statisticFriendCount 统计好友数量
func (l *UpdateUserCountMapLogic) statisticFriendCount(in *pb.UpdateUserCountMapReq) (*pb.UpdateUserCountMapResp, error) {
	countFriendResp, err := l.svcCtx.FriendService.CountFriend(l.ctx, &pb.CountFriendReq{
		Header: in.Header,
	})
	if err != nil {
		l.Errorf("statistic friend count error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	err = l.svcCtx.UserCollection.UpdateOne(l.ctx, bson.M{
		"userId": in.Header.UserId,
	}, bson.M{
		"$set": bson.M{
			"countMap." + in.CountType.String(): countFriendResp.Count,
		},
	})
	if err != nil {
		l.Errorf("update user count map error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	return &pb.UpdateUserCountMapResp{}, nil
}

// statisticJoinGroupCount 统计加入群组数量
func (l *UpdateUserCountMapLogic) statisticJoinGroupCount(in *pb.UpdateUserCountMapReq) (*pb.UpdateUserCountMapResp, error) {
	countJoinGroupResp, err := l.svcCtx.GroupService.CountJoinGroup(l.ctx, &pb.CountJoinGroupReq{
		Header: in.Header,
	})
	if err != nil {
		l.Errorf("statistic join group count error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	err = l.svcCtx.UserCollection.UpdateOne(l.ctx, bson.M{
		"userId": in.Header.UserId,
	}, bson.M{
		"$set": bson.M{
			"countMap." + in.CountType.String(): countJoinGroupResp.Count,
		},
	})
	if err != nil {
		l.Errorf("update user count map error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	return &pb.UpdateUserCountMapResp{}, nil
}

// statisticCreateGroupCount 统计创建群组数量
func (l *UpdateUserCountMapLogic) statisticCreateGroupCount(in *pb.UpdateUserCountMapReq) (*pb.UpdateUserCountMapResp, error) {
	countCreateGroupResp, err := l.svcCtx.GroupService.CountCreateGroup(l.ctx, &pb.CountCreateGroupReq{
		Header: in.Header,
	})
	if err != nil {
		l.Errorf("statistic create group count error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	err = l.svcCtx.UserCollection.UpdateOne(l.ctx, bson.M{
		"userId": in.Header.UserId,
	}, bson.M{
		"$set": bson.M{
			"countMap." + in.CountType.String(): countCreateGroupResp.Count,
		},
	}, opts.UpdateOptions{
		// 防止 cannot create field 'countMap' in element {countMap: null} 的错误
		UpdateOptions: options.Update().SetUpsert(true),
	})
	if err != nil {
		l.Errorf("update user count map error: %v", err)
		return &pb.UpdateUserCountMapResp{}, err
	}
	return &pb.UpdateUserCountMapResp{}, nil
}
