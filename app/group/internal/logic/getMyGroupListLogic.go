package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
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
	model := &groupmodel.GroupMember{}
	var groupIds []string
	err := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.CommonReq.UserId).Pluck("groupId", &groupIds).Error
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{}, err
	}
	return &pb.GetMyGroupListResp{
		Ids: groupIds,
	}, nil
}
