package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"

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
	_, err := l.svcCtx.Mongo().Collection(&relationmodel.Friend{}).RemoveAll(l.ctx, bson.M{
		"$or": []bson.M{{
			"userId":   in.Requester.Id,
			"friendId": in.UserId,
		}, {
			"userId":   in.UserId,
			"friendId": in.Requester.Id,
		}},
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
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&relationmodel.Friend{}), in.UserId)
			_, _ = relationmodel.GetMyFriendList(ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&relationmodel.Friend{}), in.Requester.Id)
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
