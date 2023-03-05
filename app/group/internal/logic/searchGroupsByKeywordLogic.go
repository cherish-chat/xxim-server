package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"strconv"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGroupsByKeywordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchGroupsByKeywordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupsByKeywordLogic {
	return &SearchGroupsByKeywordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchGroupsByKeyword 搜索群组
func (l *SearchGroupsByKeywordLogic) SearchGroupsByKeyword(in *pb.SearchGroupsByKeywordReq) (*pb.SearchGroupsByKeywordResp, error) {
	// 如果 keyword 是纯数字，那么就是群号
	_, e := strconv.Atoi(in.Keyword)
	var groups []*pb.GroupBaseInfo
	if e == nil {
		// 使用id查询
		var groupByIds *pb.MapGroupByIdsResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "MapGroupByIds", func(ctx context.Context) {
			groupByIds, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{Ids: []string{in.Keyword}})
			if err != nil {
				l.Errorf("map group by ids failed, err: %v", err)
			}
		})
		if err != nil {
			return &pb.SearchGroupsByKeywordResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		for _, bytes := range groupByIds.GroupMap {
			groups = append(groups, groupmodel.GroupFromBytes(bytes).GroupBaseInfo())
		}
	}
	if e != nil || len(groups) == 0 {
		// 使用名字查询
		err := l.svcCtx.Mysql().Model(&groupmodel.Group{}).Where("name like ?", in.Keyword+"%").Find(&groups).Error
		if err != nil {
			l.Errorf("search groups by keyword failed, err: %v", err)
			return &pb.SearchGroupsByKeywordResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.SearchGroupsByKeywordResp{Groups: groups}, nil
}
