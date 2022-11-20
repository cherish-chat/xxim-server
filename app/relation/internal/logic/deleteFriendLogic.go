package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"

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
		err := tx.Model(&relationmodel.Friend{}).Where("userId = ? and friendId = ?", in.Requester.Id, in.UserId).Delete(&relationmodel.Friend{}).Error
		if err != nil {
			l.Errorf("delete friend failed, err: %v", err)
			return err
		}
		err = tx.Model(&relationmodel.Friend{}).Where("userId = ? and friendId = ?", in.UserId, in.Requester.Id).Delete(&relationmodel.Friend{}).Error
		if err != nil {
			l.Errorf("delete friend failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorf("DeleteFriend failed, err: %v", err)
		return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	{
		// 删除缓存
		err := relationmodel.FlushFriendList(l.ctx, l.svcCtx.Redis(), in.UserId, in.Requester.Id)
		if err != nil {
			l.Errorf("FlushFriendList failed, err: %v", err)
			return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		// 预热缓存
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.UserId)
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.Requester.Id)
		}, nil)
	}
	if in.Block {
		xtrace.StartFuncSpan(l.ctx, "BlockUser", func(ctx context.Context) {
			_, err = NewBlockUserLogic(ctx, l.svcCtx).BlockUser(&pb.BlockUserReq{
				UserId:    in.UserId,
				Requester: in.Requester,
			})
		})
		if err != nil {
			l.Errorf("BlockUser failed, err: %v", err)
			return &pb.DeleteFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.DeleteFriendResp{}, nil
}
