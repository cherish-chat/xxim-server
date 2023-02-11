package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MapGroupByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMapGroupByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MapGroupByIdsLogic {
	return &MapGroupByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MapGroupByIds 获取群聊信息
func (l *MapGroupByIdsLogic) MapGroupByIds(in *pb.MapGroupByIdsReq) (*pb.MapGroupByIdsResp, error) {
	var groups []*groupmodel.Group
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetGroupByIds", func(ctx context.Context) {
		groups, err = groupmodel.ListGroupByIdsFromRedis(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.Ids)
	})
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.MapGroupByIdsResp{}, err
	}
	var groupMap = make(map[string][]byte)
	for _, group := range groups {
		groupMap[group.Id] = group.Bytes()
	}
	return &pb.MapGroupByIdsResp{GroupMap: groupMap}, nil
}
