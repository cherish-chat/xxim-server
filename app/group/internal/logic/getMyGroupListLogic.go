package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyGroupListLogic {
	return &GetMyGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMyGroupList 获取我的群聊列表
func (l *GetMyGroupListLogic) GetMyGroupList(in *pb.GetMyGroupListReq) (*pb.GetMyGroupListResp, error) {
	if in.Opt == pb.GetMyGroupListReq_DEFAULT {
		return l.getMyGroupListDefault(in)
	} else if in.Opt == pb.GetMyGroupListReq_ONLY_ID {
		return l.getMyGroupListOnlyId(in)
	}
	return &pb.GetMyGroupListResp{}, nil
}

func (l *GetMyGroupListLogic) getMyGroupListDefault(in *pb.GetMyGroupListReq) (*pb.GetMyGroupListResp, error) {
	// todo: add your logic here and delete this line
	return &pb.GetMyGroupListResp{}, nil
}

func (l *GetMyGroupListLogic) getMyGroupListOnlyId(in *pb.GetMyGroupListReq) (*pb.GetMyGroupListResp, error) {
	type res struct {
		Id string `bson:"groupId"`
	}
	var result []res
	filter := bson.M{
		"userId": in.Requester.Id,
	}
	if in.Filter != nil {
		if in.Filter.FilterFold {
			filter["fold"] = bson.M{
				"$ne": true,
			}
		}
		if in.Filter.FilterShield {
			filter["shield"] = bson.M{
				"$ne": true,
			}
		}
	}
	err := l.svcCtx.Mongo().Collection(&groupmodel.GroupMember{}).Find(l.ctx, filter).Select(bson.M{
		"groupId": 1,
	}).All(&result)
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{}, err
	}
	var ids []string
	for _, v := range result {
		ids = append(ids, v.Id)
	}
	return &pb.GetMyGroupListResp{
		Ids: ids,
	}, nil
}
