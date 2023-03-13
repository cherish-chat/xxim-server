package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearZombieMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearZombieMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearZombieMemberLogic {
	return &ClearZombieMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearZombieMember 清除僵尸用户
func (l *ClearZombieMemberLogic) ClearZombieMember(in *pb.ClearZombieMemberReq) (*pb.ClearZombieMemberResp, error) {
	// 查询 group_member 表中  所有 user.role = usermodel.RoleZombie 的成员
	logic := NewBatchKickGroupMemberLogic(context.Background(), l.svcCtx)
	for {
		var memberIds []string
		var groupMemberModel = &groupmodel.GroupMember{}
		var userModel = &usermodel.User{}
		err := l.svcCtx.Mysql().Model(groupMemberModel).
			Select("users.id AS memberId").
			Joins(fmt.Sprintf("INNER JOIN %s AS users ON %s.userId=users.id", userModel.TableName(), groupMemberModel.TableName())).
			Where(fmt.Sprintf("%s.groupId = ? AND users.role = ?", groupMemberModel.TableName()), in.GroupId, usermodel.RoleZombie).
			Limit(1000).
			Pluck("memberId", &memberIds).Error
		if err != nil {
			l.Errorf("mysql join find error: %v", err)
			return &pb.ClearZombieMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if len(memberIds) == 0 {
			return &pb.ClearZombieMemberResp{}, nil
		}
		_, err = logic.kickGroupMember(&pb.BatchKickGroupMemberReq{
			CommonReq: in.CommonReq,
			GroupId:   in.GroupId,
			MemberIds: memberIds,
		})
		if err != nil {
			l.Errorf("kickGroupMember error: %v", err)
			return &pb.ClearZombieMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
}
