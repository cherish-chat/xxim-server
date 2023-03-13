package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchKickGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchKickGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchKickGroupMemberLogic {
	return &BatchKickGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchKickGroupMember 批量踢出群成员
func (l *BatchKickGroupMemberLogic) BatchKickGroupMember(in *pb.BatchKickGroupMemberReq) (*pb.BatchKickGroupMemberResp, error) {
	if len(in.MemberIds) == 0 {
		return &pb.BatchKickGroupMemberResp{}, nil
	}
	// 获取群里所有的管理员
	groupManagers, err := getAllGroupManager(l.ctx, l.svcCtx, in.GroupId, true)
	if err != nil {
		l.Errorf("HandleGroupApply getAllGroupManager error: %v", err)
		return &pb.BatchKickGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 判断是否是管理员
	isManager := false
	for _, manager := range groupManagers {
		if manager.UserId == in.CommonReq.UserId {
			isManager = true
			break
		}
	}
	if !isManager {
		return &pb.BatchKickGroupMemberResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "您不是群管理员"),
		)}, nil
	}
	return l.kickGroupMember(in)
}

func (l *BatchKickGroupMemberLogic) kickGroupMember(in *pb.BatchKickGroupMemberReq) (*pb.BatchKickGroupMemberResp, error) {
	var err error
	// 踢出群成员
	xtrace.StartFuncSpan(l.ctx, "KickGroupMember.Transaction", func(ctx context.Context) {
		err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			// groupmember表
			member := &groupmodel.GroupMember{}
			err := tx.Model(member).Where("groupId = ? and userId in (?)", in.GroupId, in.MemberIds).Delete(member).Error
			if err != nil {
				l.Errorf("KickGroupMember groupmember delete error: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			var notices []*noticemodel.Notice
			for _, memberId := range in.MemberIds {
				tip := ""
				if in.CommonReq.UserId != memberId {
					tip = memberId + l.svcCtx.T(in.CommonReq.Language, "被移出群聊")
				} else {
					return errors.New(l.svcCtx.T(in.CommonReq.Language, "不能移出自己"))
				}
				notice := &noticemodel.Notice{
					ConvId: pb.HiddenConvIdGroup(in.GroupId),
					Options: noticemodel.NoticeOption{
						StorageForClient: false,
						UpdateConvNotice: false,
					},
					ContentType: pb.NoticeContentType_GroupMemberLeave,
					Content: utils.AnyToBytes(pb.NoticeContent_GroupMemberLeave{
						GroupId:  in.GroupId,
						Tip:      tip,
						MemberId: memberId,
					}),
					UniqueId: utils.GenId(),
					Title:    "",
					Ext:      nil,
				}
				notices = append(notices, notice)
			}
			err = noticemodel.BatchInsert(tx, notices, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			var notices []*noticemodel.Notice
			for _, memberId := range in.MemberIds {
				notice := &noticemodel.Notice{
					ConvId: pb.HiddenConvIdCommand(),
					Options: noticemodel.NoticeOption{
						StorageForClient: false,
						UpdateConvNotice: false,
					},
					UserId:      memberId,
					ContentType: pb.NoticeContentType_GroupMemberLeave,
					Content: utils.AnyToBytes(pb.NoticeContent_GroupMemberLeave{
						GroupId:  in.GroupId,
						Tip:      memberId + l.svcCtx.T(in.CommonReq.Language, "被移出群聊"),
						MemberId: memberId,
					}),
					UniqueId: utils.GenId(),
					Title:    "",
					Ext:      nil,
				}
				notices = append(notices, notice)
			}
			err = noticemodel.BatchInsert(tx, notices, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			return groupmodel.FlushGroupMemberListCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
		})
	})
	if err != nil {
		l.Errorf("KickGroupMember error: %v", err)
		return &pb.BatchKickGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	{
		// 预热缓存
		// 刷新订阅
		l.ctx = xtrace.NewContext(l.ctx)
		go utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			_, err = NewSyncGroupMemberCountLogic(l.ctx, l.svcCtx).SyncGroupMemberCount(&pb.SyncGroupMemberCountReq{
				CommonReq: in.GetCommonReq(),
				GroupId:   in.GroupId,
			})
			if err != nil {
				l.Errorf("SyncGroupMemberCount failed, err: %v", err)
				return err
			}
			// 删除缓存
			{
				err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), in.MemberIds...)
				if err != nil {
					l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
					return err
				}
			}
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: in.MemberIds})
			if err != nil {
				l.Errorf("FlushUsersSubConv failed, err: %v", err)
				return err
			}
			_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				ConvId:    pb.HiddenConvIdGroup(in.GroupId),
			})
			if err != nil {
				l.Errorf("SendNoticeData failed, err: %v", err)
			}
			for _, id := range in.MemberIds {
				_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
					CommonReq: in.CommonReq,
					ConvId:    pb.HiddenConvIdCommand(),
					UserId:    id,
				})
				if err != nil {
					l.Errorf("SendNoticeData failed, err: %v", err)
				}
			}
			return err
		})
	}
	return &pb.BatchKickGroupMemberResp{}, nil
}
