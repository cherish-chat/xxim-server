package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MapGroupMemberInfoByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMapGroupMemberInfoByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MapGroupMemberInfoByIdsLogic {
	return &MapGroupMemberInfoByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MapGroupMemberInfoByIds 批量获取群成员信息
func (l *MapGroupMemberInfoByIdsLogic) MapGroupMemberInfoByIds(in *pb.MapGroupMemberInfoByIdsReq) (*pb.MapGroupMemberInfoByIdsResp, error) {
	memberInfoMap := make(map[string]*pb.GroupMemberInfo)
	if len(in.MemberIds) == 0 {
		return &pb.MapGroupMemberInfoByIdsResp{
			CommonResp:         pb.NewSuccessResp(),
			GroupMemberInfoMap: memberInfoMap,
		}, nil
	}
	members, err := groupmodel.ListGroupMemberFromRedis(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.GroupId, in.MemberIds)
	if err != nil {
		l.Errorf("getGroupMemberInfoLogic err: %v", err)
		return &pb.MapGroupMemberInfoByIdsResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(members) == 0 {
		return &pb.MapGroupMemberInfoByIdsResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "群成员不存在"),
		)}, gorm.ErrRecordNotFound
	}
	for _, member := range members {
		memberInfoMap[member.UserId] = member.Pb()
	}
	return &pb.MapGroupMemberInfoByIdsResp{
		CommonResp:         pb.NewSuccessResp(),
		GroupMemberInfoMap: memberInfoMap,
	}, nil
}
