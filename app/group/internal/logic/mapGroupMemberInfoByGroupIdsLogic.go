package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MapGroupMemberInfoByGroupIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMapGroupMemberInfoByGroupIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MapGroupMemberInfoByGroupIdsLogic {
	return &MapGroupMemberInfoByGroupIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MapGroupMemberInfoByGroupIds 批量获取群成员信息
func (l *MapGroupMemberInfoByGroupIdsLogic) MapGroupMemberInfoByGroupIds(in *pb.MapGroupMemberInfoByGroupIdsReq) (*pb.MapGroupMemberInfoByIdsResp, error) {
	memberInfoMap := make(map[string]*pb.GroupMemberInfo)
	if len(in.GetGroupIds()) == 0 {
		return &pb.MapGroupMemberInfoByIdsResp{
			CommonResp:         pb.NewSuccessResp(),
			GroupMemberInfoMap: memberInfoMap,
		}, nil
	}
	members, err := groupmodel.ListGroupMemberFromRedisByGroupIds(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.MemberId, in.GroupIds)
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
	var userIds []string
	for _, member := range members {
		userIds = append(userIds, member.UserId)
	}
	userIds = utils.Set(userIds)
	var userMap = make(map[string]*usermodel.User)
	if in.GetOpt().GetUserBaseInfo() {
		mapUserByIds, err := l.svcCtx.UserService().MapUserByIds(l.ctx, &pb.MapUserByIdsReq{
			CommonReq: in.CommonReq,
			Ids:       userIds,
		})
		if err != nil {
			l.Errorf("getGroupMemberInfoLogic err: %v", err)
			return &pb.MapGroupMemberInfoByIdsResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		for _, bytes := range mapUserByIds.Users {
			user := usermodel.UserFromBytes(bytes)
			userMap[user.Id] = user
		}
	}
	for _, member := range members {
		info := member.Pb()
		if in.GetOpt().GetUserBaseInfo() {
			if user, ok := userMap[member.UserId]; ok {
				info.UserBaseInfo = user.BaseInfo()
			} else {
				info.UserBaseInfo = &pb.UserBaseInfo{
					Id:       member.UserId,
					Nickname: "用户已注销",
					Avatar:   "",
				}
			}
		}
		memberInfoMap[member.GroupId] = info
	}
	return &pb.MapGroupMemberInfoByIdsResp{
		CommonResp:         pb.NewSuccessResp(),
		GroupMemberInfoMap: memberInfoMap,
	}, nil
}
