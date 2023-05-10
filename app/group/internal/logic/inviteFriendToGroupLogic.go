package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteFriendToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	now time.Time
}

func NewInviteFriendToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteFriendToGroupLogic {
	return &InviteFriendToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		now:    time.Now(),
	}
}

// InviteFriendToGroup 邀请好友加入群聊
func (l *InviteFriendToGroupLogic) InviteFriendToGroup(in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	if len(in.FriendIds) == 0 {
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "请选择好友"))}, nil
	}
	// 验证是否是我的好友
	areFriendsResp, err := l.svcCtx.RelationService().AreFriends(l.ctx, &pb.AreFriendsReq{
		CommonReq: in.CommonReq,
		A:         in.CommonReq.UserId,
		BList:     in.FriendIds,
	})
	if err != nil {
		l.Errorf("InviteFriendToGroup AreFriends error: %v", err)
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, id := range in.FriendIds {
		if is, ok := areFriendsResp.FriendList[id]; !ok || !is {
			l.Errorf("InviteFriendToGroup AreFriends error: %v", err)
			return &pb.InviteFriendToGroupResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "只能邀请好友加入群聊"))}, err
		}
	}
	// 自己是不是管理员或群主
	// 获取群里所有的管理员
	groupManagers, err := getAllGroupManager(l.ctx, l.svcCtx, in.GroupId, true)
	if err != nil {
		l.Errorf("getAllGroupManager error: %v", err)
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
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
		return l.inviteFriendToGroup(in)
	} else {
		// 直接进
		err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			var members []*groupmodel.GroupMember
			for _, user := range in.FriendIds {
				members = append(members, &groupmodel.GroupMember{
					GroupId:    in.GroupId,
					UserId:     user,
					CreateTime: time.Now().UnixMilli(),
					Role:       groupmodel.RoleType_MEMBER,
					Remark:     "",
					UnbanTime:  0,
				})
			}
			// 忽略唯一索引冲突
			err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(members, 100).Error
			if err != nil {
				l.Errorf("RandInsertZombieMember error: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			var notices []*noticemodel.Notice
			for _, user := range in.FriendIds {
				notice := &noticemodel.Notice{
					ConvId: pb.HiddenConvIdGroup(in.GroupId),
					Options: noticemodel.NoticeOption{
						StorageForClient: false,
						UpdateConvNotice: false,
					},
					ContentType: pb.NoticeContentType_NewGroupMember,
					Content: utils.AnyToBytes(pb.NoticeContent_NewGroupMember{
						GroupId:  in.GroupId,
						MemberId: user,
					}),
					UniqueId: utils.GenId(),
					Title:    "",
					Ext:      nil,
				}
				notices = append(notices, notice)
			}
			err := noticemodel.BatchInsert(tx, notices, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("RandInsertZombieMember error: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			return groupmodel.FlushGroupMemberListCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
		}, func(tx *gorm.DB) error {
			return groupmodel.FlushGroupMemberCache(l.ctx, l.svcCtx.Redis(), in.GroupId, in.FriendIds...)
		})
		if err != nil {
			l.Errorf("RandInsertZombieMember error: %v", err)
			return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			// 删除缓存
			{
				err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), in.FriendIds...)
				if err != nil {
					l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
					return err
				}
			}
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: in.FriendIds})
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
			_, err = NewSyncGroupMemberCountLogic(l.ctx, l.svcCtx).SyncGroupMemberCount(&pb.SyncGroupMemberCountReq{
				CommonReq: in.GetCommonReq(),
				GroupId:   in.GroupId,
			})
			if err != nil {
				l.Errorf("SyncGroupMemberCount failed, err: %v", err)
				return err
			}
			return nil
		})
		return &pb.InviteFriendToGroupResp{}, nil
	}
}

func (l *InviteFriendToGroupLogic) inviteFriendToGroup(in *pb.InviteFriendToGroupReq) (*pb.InviteFriendToGroupResp, error) {
	myBaseInfo, err := l.svcCtx.UserService().BatchGetUserBaseInfo(l.ctx, &pb.BatchGetUserBaseInfoReq{
		CommonReq: in.CommonReq,
		Ids:       []string{in.CommonReq.UserId},
	})
	if err != nil {
		l.Errorf("InviteFriendToGroup BatchGetUserBaseInfo error: %v", err)
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(myBaseInfo.UserBaseInfos) == 0 {
		l.Errorf("InviteFriendToGroup BatchGetUserBaseInfo error: %v", err)
		return &pb.InviteFriendToGroupResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "用户不存在"))}, err
	}
	logic := NewApplyToBeGroupMemberLogic(l.ctx, l.svcCtx)
	for _, id := range in.FriendIds {
		resp, err := logic.ApplyToBeGroupMember(&pb.ApplyToBeGroupMemberReq{
			CommonReq: &pb.CommonReq{
				UserId: id,
			},
			GroupId: in.GroupId,
			Reason:  fmt.Sprintf("[%s]%s", myBaseInfo.UserBaseInfos[0].Id, myBaseInfo.UserBaseInfos[0].Nickname) + l.svcCtx.T(in.CommonReq.Language, "邀请我加入群聊"),
		})
		if err != nil {
			l.Errorf("InviteFriendToGroup ApplyToBeGroupMember error: %v", err)
			if resp != nil {
				return &pb.InviteFriendToGroupResp{CommonResp: resp.GetCommonResp()}, err
			} else {
				return &pb.InviteFriendToGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
		}
	}
	return &pb.InviteFriendToGroupResp{}, nil
}
