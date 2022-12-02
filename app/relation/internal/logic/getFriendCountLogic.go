package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendCountLogic {
	return &GetFriendCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendCountLogic) GetFriendCount(in *pb.GetFriendCountReq) (*pb.GetFriendCountResp, error) {
	hLen, err := xredis.HLen(l.svcCtx.Redis(), l.ctx, rediskey.FriendList(in.CommonReq.UserId))
	if err != nil {
		friendList, err := relationmodel.GetMyFriendList(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.CommonReq.UserId)
		if err != nil {
			l.Errorf("GetFriendCount failed, err: %v", err)
			return &pb.GetFriendCountResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		return &pb.GetFriendCountResp{Count: int32(len(friendList))}, nil
	}
	if hLen == 1 {
		// 是否是Not Found
		exist, err := l.svcCtx.Redis().HexistsCtx(l.ctx, rediskey.FriendList(in.CommonReq.UserId), xredis.NotFound)
		if err != nil {
			l.Errorf("GetFriendCount failed, err: %v", err)
			return &pb.GetFriendCountResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if exist {
			// 0
			return &pb.GetFriendCountResp{Count: 0}, nil
		}
	}
	return &pb.GetFriendCountResp{Count: int32(hLen)}, nil
}
