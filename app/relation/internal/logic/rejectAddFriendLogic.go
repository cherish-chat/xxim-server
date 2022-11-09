package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectAddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectAddFriendLogic {
	return &RejectAddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectAddFriendLogic) RejectAddFriend(in *pb.RejectAddFriendReq) (*pb.RejectAddFriendResp, error) {
	apply := &relationmodel.RequestAddFriend{Id: in.RequestId}
	err := l.svcCtx.Mongo().Collection(apply).Find(l.ctx, bson.M{
		"_id": apply.Id,
	}).One(apply)
	if err != nil {
		l.Errorf("FindOne failed, err: %v", err)
		return &pb.RejectAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	apply.Status = pb.RequestAddFriendStatus_Refused
	apply.UpdateTime = time.Now().UnixMilli()
	err = l.svcCtx.Mongo().Collection(apply).UpdateOne(l.ctx, bson.M{
		"$or": []bson.M{{
			"fromUserId": in.ApplyUserId,
			"toUserId":   in.Requester.Id,
		}, {
			"fromUserId": in.Requester.Id,
			"toUserId":   in.ApplyUserId,
		}},
		"status": pb.RequestAddFriendStatus_Unhandled,
	}, bson.M{
		"$set": bson.M{
			"status":      apply.Status,
			"update_time": apply.UpdateTime,
		},
	})
	if err != nil {
		l.Errorf("Upsert failed, err: %v", err)
		return &pb.RejectAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 是否拉黑
	if in.Block {
		xtrace.StartFuncSpan(l.ctx, "BlockUser", func(ctx context.Context) {
			_, err = NewBlockUserLogic(ctx, l.svcCtx).BlockUser(&pb.BlockUserReq{
				Requester: in.Requester,
				UserId:    in.ApplyUserId,
			})
		})
		if err != nil {
			l.Errorf("BlockUser failed, err: %v", err)
			return &pb.RejectAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.RejectAddFriendResp{}, nil
}
