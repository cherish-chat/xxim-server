package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddGroupMember 添加群成员
func (l *AddGroupMemberLogic) AddGroupMember(in *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		member := &groupmodel.GroupMember{
			GroupId:    in.GroupId,
			UserId:     in.UserId,
			CreateTime: time.Now().UnixMilli(),
			Role:       groupmodel.RoleType_MEMBER,
		}
		err := xorm.InsertOne(tx, member)
		if err != nil {
			// 判断是不是唯一索引冲突
			if !xorm.DuplicateError(err) {
				l.Errorf("InviteFriendToGroup InsertMany error: %v", err)
				return err
			}
		}
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvIdGroup(in.GroupId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_NewGroupMember,
			Content: utils.AnyToBytes(pb.NoticeContent_NewGroupMember{
				GroupId:  in.GroupId,
				MemberId: in.UserId,
			}),
			UniqueId: "member",
			Title:    "",
			Ext:      nil,
		}
		err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		return groupmodel.FlushGroupMemberListCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
	})
	if err != nil {
		l.Errorf("AddGroupMember error: %v", err)
		return &pb.AddGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
		// 删除缓存
		{
			err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), in.UserId)
			if err != nil {
				l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
				return err
			}
		}
		// 预热缓存
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarmUp", func(ctx context.Context) {
			_, err = groupmodel.ListGroupsByUserIdFromMysql(ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.UserId)
		}, propagation.MapCarrier{
			"group_id": in.GroupId,
		})
		_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: []string{in.UserId}})
		if err != nil {
			l.Errorf("FlushUsersSubConv failed, err: %v", err)
			return err
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
	return &pb.AddGroupMemberResp{}, nil
}
