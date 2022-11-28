package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
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
	err := xorm.DetailByWhere(l.svcCtx.Mysql(), apply, xorm.Where("id", apply.Id))
	if err != nil {
		l.Errorf("FindOne failed, err: %v", err)
		return &pb.RejectAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	apply.Status = pb.RequestAddFriendStatus_Refused
	apply.UpdateTime = time.Now().UnixMilli()
	err = l.svcCtx.Mysql().Model(apply).
		Where("status = ? AND ((fromUserId = ? AND toUserId = ?) OR (fromUserId = ? AND toUserId = ?))",
			pb.RequestAddFriendStatus_Unhandled,
			in.ApplyUserId, in.CommonReq.UserId,
			in.CommonReq.UserId, in.ApplyUserId).
		Updates(map[string]interface{}{
			"status":     apply.Status,
			"updateTime": apply.UpdateTime,
		}).Error
	if err != nil {
		l.Errorf("Upsert failed, err: %v", err)
		return &pb.RejectAddFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 是否拉黑
	if in.Block {
		xtrace.StartFuncSpan(l.ctx, "BlockUser", func(ctx context.Context) {
			_, err = NewBlockUserLogic(ctx, l.svcCtx).BlockUser(&pb.BlockUserReq{
				CommonReq: in.CommonReq,
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
