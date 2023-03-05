package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupHomeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupHomeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupHomeLogic {
	return &GetGroupHomeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupHome 获取群聊首页
func (l *GetGroupHomeLogic) GetGroupHome(in *pb.GetGroupHomeReq) (*pb.GetGroupHomeResp, error) {
	var groupByIds *pb.MapGroupByIdsResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "", func(ctx context.Context) {
		groupByIds, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
			CommonReq: in.CommonReq,
			Ids:       []string{in.GroupId},
		})
	})
	if err != nil {
		l.Errorf("MapGroupByIds error: %s", err.Error())
		return &pb.GetGroupHomeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	bytes, ok := groupByIds.GroupMap[in.GroupId]
	if !ok {
		l.Errorf("group not found: %s", in.GroupId)
		return &pb.GetGroupHomeResp{CommonResp: pb.NewToastErrorResp("群聊不存在")}, nil
	}
	group := groupmodel.GroupFromBytes(bytes)
	return &pb.GetGroupHomeResp{
		GroupId:          group.Id,
		Name:             group.Name,
		Avatar:           group.Avatar,
		CreatedAt:        utils.TimeFormat(group.CreateTime),
		MemberCount:      int32(group.MemberCount),
		Introduction:     "群主很懒，什么都没写",
		MemberStatistics: nil,
	}, nil
}
