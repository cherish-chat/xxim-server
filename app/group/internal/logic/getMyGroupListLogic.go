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
		resp, err := l.getMyGroupListDefault(in)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		}
		return resp, err
	} else if in.Opt == pb.GetMyGroupListReq_ONLY_ID {
		resp, err := l.getMyGroupListOnlyId(in)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		}
		return resp, err
	} else if in.Opt == pb.GetMyGroupListReq_WITH_MY_MEMBER_INFO {
		getMyGroupListResp, err := l.getMyGroupListDefault(in)
		if err != nil {
			return &pb.GetMyGroupListResp{}, err
		}
		resp, err := l.getMyGroupListWithMyMemberInfo(in, getMyGroupListResp)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		}
		return resp, err
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
	var groupIds []string
	var err error
	groupIds, err = groupmodel.ListGroupsByUserIdFromRedis(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.CommonReq.UserId)
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{}, err
	}
	var mapGroupByIdsResp *pb.MapGroupByIdsResp
	xtrace.StartFuncSpan(l.ctx, "MapGroupByIds", func(ctx context.Context) {
		mapGroupByIdsResp, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
			Ids: groupIds,
		})
	})
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 只获取有效的群聊
	var validGroupIds []string
	for _, id := range groupIds {
		_, ok := mapGroupByIdsResp.GroupMap[id]
		if ok {
			validGroupIds = append(validGroupIds, id)
		}
	}
	return &pb.GetMyGroupListResp{
		Ids: validGroupIds,
	}, nil
}

func (l *GetMyGroupListLogic) getMyGroupListWithMyMemberInfo(in *pb.GetMyGroupListReq, getMyGroupListResp *pb.GetMyGroupListResp) (*pb.GetMyGroupListResp, error) {
	var groupIds []string
	var mapGruopMemberInfoByIdsResp *pb.MapGroupMemberInfoByIdsResp
	var err error
	for _, info := range getMyGroupListResp.GroupMap {
		groupIds = append(groupIds, info.Id)
	}
	xtrace.StartFuncSpan(l.ctx, "MapGruopMemberInfoByIds", func(ctx context.Context) {
		mapGruopMemberInfoByIdsResp, err = NewMapGroupMemberInfoByGroupIdsLogic(ctx, l.svcCtx).MapGroupMemberInfoByGroupIds(&pb.MapGroupMemberInfoByGroupIdsReq{
			CommonReq: &pb.CommonReq{},
			GroupIds:  groupIds,
			MemberId:  in.GetCommonReq().GetUserId(),
			Opt:       nil,
		})
	})
	if err != nil {
		l.Errorf("get group list error: %v", err)
		return &pb.GetMyGroupListResp{}, err
	}
	for _, info := range getMyGroupListResp.GroupMap {
		memberInfo, ok := mapGruopMemberInfoByIdsResp.GroupMemberInfoMap[info.Id]
		if ok {
			info.MyMemberInfo = memberInfo
		}
	}
	return getMyGroupListResp, nil
}

func (l *GetMyGroupListLogic) withConvSetting(in *pb.GetMyGroupListReq, resp *pb.GetMyGroupListResp) (*pb.GetMyGroupListResp, error) {
	convIds := make([]string, 0)
	for _, groupId := range resp.Ids {
		convIds = append(convIds, pb.GroupConvId(groupId))
	}
	resp.ConvSettingMap = make(map[string]*pb.ConvSetting)
	resp.ConvSetting2Map = make(map[string]*pb.ConvSettingProto2)
	if len(convIds) == 0 {
		return resp, nil
	}
	getConvSettingResp, err := l.svcCtx.ImService().GetConvSetting(l.ctx, &pb.GetConvSettingReq{
		CommonReq: in.GetCommonReq(),
		ConvIds:   convIds,
	})
	if err != nil {
		l.Errorf("GetConvSetting err: %v", err)
		return resp, err
	}
	for _, setting := range getConvSettingResp.ConvSettings {
		resp.ConvSettingMap[setting.ConvId] = setting
	}
	for _, setting := range getConvSettingResp.ConvSetting2S {
		resp.ConvSetting2Map[setting.ConvId] = setting
	}
	return resp, nil
}
