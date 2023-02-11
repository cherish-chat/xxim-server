package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupListByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListByUserIdLogic {
	return &GetGroupListByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupListByUserId 分页获取某人的群列表
func (l *GetGroupListByUserIdLogic) GetGroupListByUserId(in *pb.GetGroupListByUserIdReq) (*pb.GetGroupListByUserIdResp, error) {
	var groupIds []string
	var total int64
	model := &groupmodel.GroupMember{}
	var members []*groupmodel.GroupMember
	tx := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.UserId)
	tx.Count(&total)
	err := xorm.Paging(tx.Order("createTime DESC"), in.Page.Page, in.Page.Size).Find(&members).Error
	if err != nil {
		l.Errorf("GetGroupListByUserId err: %v", err)
		return &pb.GetGroupListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, member := range members {
		groupIds = append(groupIds, member.GroupId)
	}
	if len(groupIds) == 0 {
		return &pb.GetGroupListByUserIdResp{CommonResp: pb.NewSuccessResp(), Total: total}, nil
	}
	// 获取群信息去
	groupByIds, err := NewMapGroupByIdsLogic(l.ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
		CommonReq: in.CommonReq,
		Ids:       groupIds,
	})
	if err != nil {
		l.Errorf("GetGroupListByUserId err: %v", err)
		return &pb.GetGroupListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.GetGroupListByUserIdItem
	for _, member := range members {
		groupId := member.GroupId
		group := &groupmodel.Group{}
		bytes, ok := groupByIds.GroupMap[groupId]
		if ok {
			group = groupmodel.GroupFromBytes(bytes)
		}
		resp = append(resp, &pb.GetGroupListByUserIdItem{
			GroupId:       group.Id,
			Avatar:        group.Avatar,
			Name:          group.Name,
			MemberCount:   int64(group.MemberCount),
			JoinTime:      member.CreateTime,
			JoinTimeStr:   utils.TimeFormat(member.CreateTime),
			Owner:         group.Owner,
			CreateTime:    group.CreateTime,
			CreateTimeStr: utils.TimeFormat(group.CreateTime),
		})
	}
	return &pb.GetGroupListByUserIdResp{
		GroupList: resp,
		Total:     total,
	}, nil
}
