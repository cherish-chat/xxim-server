// Code generated by goctl. DO NOT EDIT.
// Source: group.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/logic"
	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type GroupServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedGroupServiceServer
}

func NewGroupServiceServer(svcCtx *svc.ServiceContext) *GroupServiceServer {
	return &GroupServiceServer{
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建群聊
func (s *GroupServiceServer) CreateGroup(ctx context.Context, in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	l := logic.NewCreateGroupLogic(ctx, s.svcCtx)
	return l.CreateGroup(in)
}

// GetGroupHome 获取群聊首页
func (s *GroupServiceServer) GetGroupHome(ctx context.Context, in *pb.GetGroupHomeReq) (*pb.GetGroupHomeResp, error) {
	l := logic.NewGetGroupHomeLogic(ctx, s.svcCtx)
	return l.GetGroupHome(in)
}

// InviteFriendToGroup 邀请好友加入群聊
func (s *GroupServiceServer) InviteFriendToGroup(ctx context.Context, in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	l := logic.NewInviteFriendToGroupLogic(ctx, s.svcCtx)
	return l.InviteFriendToGroup(in)
}

// CreateGroupNotice 创建群公告
func (s *GroupServiceServer) CreateGroupNotice(ctx context.Context, in *pb.CreateGroupNoticeReq) (*pb.CreateGroupNoticeResp, error) {
	l := logic.NewCreateGroupNoticeLogic(ctx, s.svcCtx)
	return l.CreateGroupNotice(in)
}

// DeleteGroupNotice 删除群公告
func (s *GroupServiceServer) DeleteGroupNotice(ctx context.Context, in *pb.DeleteGroupNoticeReq) (*pb.DeleteGroupNoticeResp, error) {
	l := logic.NewDeleteGroupNoticeLogic(ctx, s.svcCtx)
	return l.DeleteGroupNotice(in)
}

// GetGroupNoticeList 获取群公告列表
func (s *GroupServiceServer) GetGroupNoticeList(ctx context.Context, in *pb.GetGroupNoticeListReq) (*pb.GetGroupNoticeListResp, error) {
	l := logic.NewGetGroupNoticeListLogic(ctx, s.svcCtx)
	return l.GetGroupNoticeList(in)
}

// SetGroupMemberInfo 设置群成员信息
func (s *GroupServiceServer) SetGroupMemberInfo(ctx context.Context, in *pb.SetGroupMemberInfoReq) (*pb.SetGroupMemberInfoResp, error) {
	l := logic.NewSetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.SetGroupMemberInfo(in)
}

// GetGroupMemberInfo 获取群成员信息
func (s *GroupServiceServer) GetGroupMemberInfo(ctx context.Context, in *pb.GetGroupMemberInfoReq) (*pb.GetGroupMemberInfoResp, error) {
	l := logic.NewGetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.GetGroupMemberInfo(in)
}

// MapGroupMemberInfoByIds 批量获取群成员信息
func (s *GroupServiceServer) MapGroupMemberInfoByIds(ctx context.Context, in *pb.MapGroupMemberInfoByIdsReq) (*pb.MapGroupMemberInfoByIdsResp, error) {
	l := logic.NewMapGroupMemberInfoByIdsLogic(ctx, s.svcCtx)
	return l.MapGroupMemberInfoByIds(in)
}

// MapGroupMemberInfoByGroupIdsReq 批量获取群成员信息
func (s *GroupServiceServer) MapGroupMemberInfoByGroupIds(ctx context.Context, in *pb.MapGroupMemberInfoByGroupIdsReq) (*pb.MapGroupMemberInfoByIdsResp, error) {
	l := logic.NewMapGroupMemberInfoByGroupIdsLogic(ctx, s.svcCtx)
	return l.MapGroupMemberInfoByGroupIds(in)
}

// EditGroupInfo 编辑群信息
func (s *GroupServiceServer) EditGroupInfo(ctx context.Context, in *pb.EditGroupInfoReq) (*pb.EditGroupInfoResp, error) {
	l := logic.NewEditGroupInfoLogic(ctx, s.svcCtx)
	return l.EditGroupInfo(in)
}

// TransferGroupOwner 转让群主
func (s *GroupServiceServer) TransferGroupOwner(ctx context.Context, in *pb.TransferGroupOwnerReq) (*pb.TransferGroupOwnerResp, error) {
	l := logic.NewTransferGroupOwnerLogic(ctx, s.svcCtx)
	return l.TransferGroupOwner(in)
}

// KickGroupMember 踢出群成员
func (s *GroupServiceServer) KickGroupMember(ctx context.Context, in *pb.KickGroupMemberReq) (*pb.KickGroupMemberResp, error) {
	l := logic.NewKickGroupMemberLogic(ctx, s.svcCtx)
	return l.KickGroupMember(in)
}

// BatchKickGroupMember 批量踢出群成员
func (s *GroupServiceServer) BatchKickGroupMember(ctx context.Context, in *pb.BatchKickGroupMemberReq) (*pb.BatchKickGroupMemberResp, error) {
	l := logic.NewBatchKickGroupMemberLogic(ctx, s.svcCtx)
	return l.BatchKickGroupMember(in)
}

// GetGroupMemberList 获取群成员列表
func (s *GroupServiceServer) GetGroupMemberList(ctx context.Context, in *pb.GetGroupMemberListReq) (*pb.GetGroupMemberListResp, error) {
	l := logic.NewGetGroupMemberListLogic(ctx, s.svcCtx)
	return l.GetGroupMemberList(in)
}

// GetMyGroupList 获取我的群聊列表
func (s *GroupServiceServer) GetMyGroupList(ctx context.Context, in *pb.GetMyGroupListReq) (*pb.GetMyGroupListResp, error) {
	l := logic.NewGetMyGroupListLogic(ctx, s.svcCtx)
	return l.GetMyGroupList(in)
}

// MapGroupByIds 获取群聊信息
func (s *GroupServiceServer) MapGroupByIds(ctx context.Context, in *pb.MapGroupByIdsReq) (*pb.MapGroupByIdsResp, error) {
	l := logic.NewMapGroupByIdsLogic(ctx, s.svcCtx)
	return l.MapGroupByIds(in)
}

// SyncGroupMemberCount 同步群成员数量
func (s *GroupServiceServer) SyncGroupMemberCount(ctx context.Context, in *pb.SyncGroupMemberCountReq) (*pb.SyncGroupMemberCountResp, error) {
	l := logic.NewSyncGroupMemberCountLogic(ctx, s.svcCtx)
	return l.SyncGroupMemberCount(in)
}

// ApplyToBeGroupMember 申请加入群聊
func (s *GroupServiceServer) ApplyToBeGroupMember(ctx context.Context, in *pb.ApplyToBeGroupMemberReq) (*pb.ApplyToBeGroupMemberResp, error) {
	l := logic.NewApplyToBeGroupMemberLogic(ctx, s.svcCtx)
	return l.ApplyToBeGroupMember(in)
}

// HandleGroupApply 处理群聊申请
func (s *GroupServiceServer) HandleGroupApply(ctx context.Context, in *pb.HandleGroupApplyReq) (*pb.HandleGroupApplyResp, error) {
	l := logic.NewHandleGroupApplyLogic(ctx, s.svcCtx)
	return l.HandleGroupApply(in)
}

// GetGroupApplyList 获取群聊申请列表
func (s *GroupServiceServer) GetGroupApplyList(ctx context.Context, in *pb.GetGroupApplyListReq) (*pb.GetGroupApplyListResp, error) {
	l := logic.NewGetGroupApplyListLogic(ctx, s.svcCtx)
	return l.GetGroupApplyList(in)
}

// GetGroupListByUserId 分页获取某人的群列表
func (s *GroupServiceServer) GetGroupListByUserId(ctx context.Context, in *pb.GetGroupListByUserIdReq) (*pb.GetGroupListByUserIdResp, error) {
	l := logic.NewGetGroupListByUserIdLogic(ctx, s.svcCtx)
	return l.GetGroupListByUserId(in)
}

// GetAllGroupModel 获取所有群组
func (s *GroupServiceServer) GetAllGroupModel(ctx context.Context, in *pb.GetAllGroupModelReq) (*pb.GetAllGroupModelResp, error) {
	l := logic.NewGetAllGroupModelLogic(ctx, s.svcCtx)
	return l.GetAllGroupModel(in)
}

// GetGroupModelDetail 获取群组详情
func (s *GroupServiceServer) GetGroupModelDetail(ctx context.Context, in *pb.GetGroupModelDetailReq) (*pb.GetGroupModelDetailResp, error) {
	l := logic.NewGetGroupModelDetailLogic(ctx, s.svcCtx)
	return l.GetGroupModelDetail(in)
}

// UpdateGroupModel 更新群组
func (s *GroupServiceServer) UpdateGroupModel(ctx context.Context, in *pb.UpdateGroupModelReq) (*pb.UpdateGroupModelResp, error) {
	l := logic.NewUpdateGroupModelLogic(ctx, s.svcCtx)
	return l.UpdateGroupModel(in)
}

// DismissGroupModel 解散群组
func (s *GroupServiceServer) DismissGroupModel(ctx context.Context, in *pb.DismissGroupModelReq) (*pb.DismissGroupModelResp, error) {
	l := logic.NewDismissGroupModelLogic(ctx, s.svcCtx)
	return l.DismissGroupModel(in)
}

// SearchGroupsByKeyword 搜索群组
func (s *GroupServiceServer) SearchGroupsByKeyword(ctx context.Context, in *pb.SearchGroupsByKeywordReq) (*pb.SearchGroupsByKeywordResp, error) {
	l := logic.NewSearchGroupsByKeywordLogic(ctx, s.svcCtx)
	return l.SearchGroupsByKeyword(in)
}

// AddGroupMember 添加群成员
func (s *GroupServiceServer) AddGroupMember(ctx context.Context, in *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	l := logic.NewAddGroupMemberLogic(ctx, s.svcCtx)
	return l.AddGroupMember(in)
}

// ReportGroup
func (s *GroupServiceServer) ReportGroup(ctx context.Context, in *pb.ReportGroupReq) (*pb.ReportGroupResp, error) {
	l := logic.NewReportGroupLogic(ctx, s.svcCtx)
	return l.ReportGroup(in)
}

// RandInsertZombieMember 随机插入僵尸用户
func (s *GroupServiceServer) RandInsertZombieMember(ctx context.Context, in *pb.RandInsertZombieMemberReq) (*pb.RandInsertZombieMemberResp, error) {
	l := logic.NewRandInsertZombieMemberLogic(ctx, s.svcCtx)
	return l.RandInsertZombieMember(in)
}

// ClearZombieMember 清除僵尸用户
func (s *GroupServiceServer) ClearZombieMember(ctx context.Context, in *pb.ClearZombieMemberReq) (*pb.ClearZombieMemberResp, error) {
	l := logic.NewClearZombieMemberLogic(ctx, s.svcCtx)
	return l.ClearZombieMember(in)
}
