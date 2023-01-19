package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncGroupMemberCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncGroupMemberCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncGroupMemberCountLogic {
	return &SyncGroupMemberCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SyncGroupMemberCount 同步群成员数量
func (l *SyncGroupMemberCountLogic) SyncGroupMemberCount(in *pb.SyncGroupMemberCountReq) (*pb.SyncGroupMemberCountResp, error) {
	// 查询群成员数量
	var count int64
	var err error
	err = l.svcCtx.Mysql().Model(&groupmodel.GroupMember{}).Where("groupId = ?", in.GroupId).Count(&count).Error
	if err != nil {
		l.Errorf("查询群成员数量失败, err: %v", err)
		return &pb.SyncGroupMemberCountResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新群成员数量
	err = l.svcCtx.Mysql().Model(&groupmodel.Group{}).Where("id = ?", in.GroupId).Update("memberCount", count).Error
	if err != nil {
		l.Errorf("更新群成员数量失败, err: %v", err)
		return &pb.SyncGroupMemberCountResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.SyncGroupMemberCountResp{}, nil
}
