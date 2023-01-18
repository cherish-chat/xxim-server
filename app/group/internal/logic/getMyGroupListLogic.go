package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"

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
	myGroupListOnlyId, err := l.getMyGroupListOnlyId(in)
	if err != nil {
		return &pb.GetMyGroupListResp{}, err
	}
	var groupMap = make(map[string]*pb.GroupBaseInfo)
	// 使用id获取群聊信息
	var mapGroupByIdsResp *pb.MapGroupByIdsResp
	xtrace.StartFuncSpan(l.ctx, "MapGroupByIds", func(ctx context.Context) {
		mapGroupByIdsResp, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
			Ids: myGroupListOnlyId.Ids,
		})
	})
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{}, err
	}
	for _, id := range myGroupListOnlyId.Ids {
		group, ok := mapGroupByIdsResp.GroupMap[id]
		if ok {
			groupMap[id] = groupmodel.GroupFromBytes(group).GroupBaseInfo()
		}
	}
	return &pb.GetMyGroupListResp{
		GroupMap: groupMap,
		Ids:      myGroupListOnlyId.Ids,
	}, nil
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
