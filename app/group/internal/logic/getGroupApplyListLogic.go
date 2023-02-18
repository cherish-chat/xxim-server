package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupApplyListLogic {
	return &GetGroupApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupApplyList 获取群聊申请列表
func (l *GetGroupApplyListLogic) GetGroupApplyList(in *pb.GetGroupApplyListReq) (*pb.GetGroupApplyListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 20}
	}
	// 1. 先获取用户管理的群聊id
	var manageGroupIds []string
	err := l.svcCtx.Mysql().Model(&groupmodel.GroupMember{}).
		Where("userId = ?", in.CommonReq.UserId).
		Where("(role = ? OR role = ?)", groupmodel.RoleType_OWNER, groupmodel.RoleType_MANAGER).
		Pluck("groupId", &manageGroupIds).Error
	if err != nil {
		l.Errorf("GetGroupApplyList failed, err: %v", err)
		return &pb.GetGroupApplyListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(manageGroupIds) == 0 {
		return &pb.GetGroupApplyListResp{}, nil
	}
	// 2. 获取群聊申请列表
	orBuilder := ""
	orArgs := make([]interface{}, 0)
	for _, groupId := range manageGroupIds {
		orBuilder += "groupId = ? OR "
		orArgs = append(orArgs, groupId)
	}
	orBuilder = orBuilder[:len(orBuilder)-4]
	orBuilder = "(" + orBuilder + ")"
	var groupApplyList []*groupmodel.GroupApply
	tx := l.svcCtx.Mysql().Model(&groupmodel.GroupApply{}).
		Where(orBuilder, orArgs...)
	if in.Filter != nil {
		if in.Filter.Result != nil {
			tx = tx.Where("result = ?", in.Filter.Result)
		}
	}
	var total int64
	tx.Count(&total)
	tx = tx.Order("createTime desc")
	tx = xorm.Paging(tx, in.Page.Page, in.Page.Size)
	err = tx.Find(&groupApplyList).Error
	if len(groupApplyList) == 0 {
		return &pb.GetGroupApplyListResp{}, nil
	}
	var resp []*pb.GroupApplyInfo
	var groupIds []string
	var userIds []string
	for _, groupApply := range groupApplyList {
		groupIds = append(groupIds, groupApply.GroupId)
		userIds = append(userIds, groupApply.UserId)
	}
	// 3. 获取群聊信息
	var groupMap *pb.MapGroupByIdsResp
	{
		xtrace.StartFuncSpan(l.ctx, "MapGroupByIds", func(ctx context.Context) {
			groupMap, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
				CommonReq: in.CommonReq,
				Ids:       groupIds,
			})
		})
		if err != nil {
			l.Errorf("GetGroupApplyList failed, err: %v", err)
			return &pb.GetGroupApplyListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	// 4. 获取用户信息
	var userMap *pb.MapUserByIdsResp
	{
		userMap, err = l.svcCtx.UserService().MapUserByIds(l.ctx, &pb.MapUserByIdsReq{
			CommonReq: in.CommonReq,
			Ids:       userIds,
		})
		if err != nil {
			l.Errorf("GetGroupApplyList failed, err: %v", err)
			return &pb.GetGroupApplyListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	for _, groupApply := range groupApplyList {
		applyInfo := groupApply.ToPB()
		groupBuf, ok := groupMap.GroupMap[applyInfo.GroupId]
		if !ok {
			applyInfo.GroupBaseInfo = &pb.GroupBaseInfo{Id: applyInfo.GroupId, Name: "群聊已解散"}
		} else {
			applyInfo.GroupBaseInfo = groupmodel.GroupFromBytes(groupBuf).GroupBaseInfo()
		}
		userBuf, ok := userMap.Users[applyInfo.UserId]
		if !ok {
			applyInfo.UserBaseInfo = &pb.UserBaseInfo{Id: applyInfo.UserId, Nickname: "用户已注销"}
		} else {
			applyInfo.UserBaseInfo = usermodel.UserFromBytes(userBuf).BaseInfo()
		}
		handleUserBuf, ok := userMap.Users[applyInfo.HandleUserId]
		if !ok {
			applyInfo.HandleUserBaseInfo = &pb.UserBaseInfo{Id: applyInfo.HandleUserId, Nickname: "用户已注销"}
		} else {
			applyInfo.HandleUserBaseInfo = usermodel.UserFromBytes(handleUserBuf).BaseInfo()
		}
		resp = append(resp, applyInfo)
	}
	return &pb.GetGroupApplyListResp{
		GroupApplyList: resp,
		Total:          total,
	}, nil
}
