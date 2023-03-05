package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFriendLogic) DeleteFriend(in *pb.DeleteFriendReq) (*pb.DeleteFriendResp, error) {
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := tx.Model(&relationmodel.Friend{}).Where("userId = ? and friendId = ?", in.CommonReq.UserId, in.UserId).Delete(&relationmodel.Friend{}).Error
		if err != nil {
			l.Errorf("delete friend failed, err: %v", err)
			return err
		}
		err = tx.Model(&relationmodel.Friend{}).Where("userId = ? and friendId = ?", in.UserId, in.CommonReq.UserId).Delete(&relationmodel.Friend{}).Error
		if err != nil {
			l.Errorf("delete friend failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		for _, userId := range []string{in.CommonReq.UserId, in.UserId} {
			notice := &noticemodel.Notice{
				ConvId: pb.HiddenConvIdCommand(),
				UserId: userId,
				Options: noticemodel.NoticeOption{
					StorageForClient: false,
					UpdateConvNotice: false,
				},
				ContentType: pb.NoticeContentType_SyncFriendList,
				Content: utils.AnyToBytes(pb.NoticeContent_SyncFriendList{
					Comment: "deleteFriend",
				}),
				UniqueId: "syncFriendList",
				Title:    "",
				Ext:      nil,
			}
			err := notice.Insert(l.ctx, tx, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.Errorf("DeleteFriend failed, err: %v", err)
		return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	{
		// 删除缓存
		err := relationmodel.FlushFriendList(l.ctx, l.svcCtx.Redis(), in.UserId, in.CommonReq.UserId)
		if err != nil {
			l.Errorf("FlushFriendList failed, err: %v", err)
			return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		// 预热缓存
		xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.UserId)
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.CommonReq.UserId)
		}, nil)
		// 刷新订阅
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: []string{
				in.UserId, in.CommonReq.UserId,
			}})
			if err != nil {
				l.Errorf("FlushUsersSubConv failed, err: %v", err)
				return err
			}
			for _, userId := range []string{in.UserId, in.CommonReq.UserId} {
				_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
					UserId: userId,
					ConvId: pb.HiddenConvIdCommand(),
				})
				if err != nil {
					l.Errorf("SendNoticeData failed, err: %v", err)
					return err
				}
			}
			return err
		})
	}
	if in.Block {
		xtrace.StartFuncSpan(l.ctx, "BlockUser", func(ctx context.Context) {
			_, err = NewBlockUserLogic(ctx, l.svcCtx).BlockUser(&pb.BlockUserReq{
				UserId:    in.UserId,
				CommonReq: in.CommonReq,
			})
		})
		if err != nil {
			l.Errorf("BlockUser failed, err: %v", err)
			return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.DeleteFriendResp{}, nil
}
