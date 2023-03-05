package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"strings"

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
	if len(in.Filter) == 0 {
		in.Filter = make(map[string]string)
	}
	var groupIds []string
	var total int64
	model := &groupmodel.GroupMember{}
	var members []*groupmodel.GroupMember
	tx := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.UserId)
	var name *string
	for k, v := range in.Filter {
		if v == "" {
			continue
		}
		switch k {
		case "name":
			v := v
			name = utils.AnyPtr(v)
		}
	}
	if gid, ok := in.Filter["id"]; ok && gid != "" {
		tx = tx.Where("groupId = ?", gid)
	} else {
		if name != nil {
			// 获取所有群
			myGroupList, err := NewGetMyGroupListLogic(l.ctx, l.svcCtx).GetMyGroupList(&pb.GetMyGroupListReq{
				CommonReq: &pb.CommonReq{UserId: in.UserId},
				Page:      &pb.Page{Page: 1, Size: 10000},
				Opt:       pb.GetMyGroupListReq_DEFAULT,
			})
			if err != nil {
				l.Errorf("GetGroupListByUserId err: %v", err)
				return &pb.GetGroupListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
			var ids []string
			for _, group := range myGroupList.GroupMap {
				if strings.Contains(group.Name, *name) {
					ids = append(ids, group.Id)
				}
			}
			if len(ids) == 0 {
				return &pb.GetGroupListByUserIdResp{CommonResp: pb.NewSuccessResp(), Total: total}, nil
			}
			tx = tx.Where("groupId in (?)", ids)
		}
	}
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
