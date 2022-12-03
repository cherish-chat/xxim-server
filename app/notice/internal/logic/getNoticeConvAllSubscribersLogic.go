package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeConvAllSubscribersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeConvAllSubscribersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeConvAllSubscribersLogic {
	return &GetNoticeConvAllSubscribersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetNoticeConvAllSubscribers 获取通知号所有的订阅者
func (l *GetNoticeConvAllSubscribersLogic) GetNoticeConvAllSubscribers(in *pb.GetNoticeConvAllSubscribersReq) (*pb.GetNoticeConvAllSubscribersResp, error) {
	// ZRANGEBYSCORE conv:subscribers:group:1 min +inf
	min := time.Now().AddDate(0, 0, -1).UnixMilli()
	//if in.LastActiveTime != nil {
	//	min = *in.LastActiveTime
	//}
	val, err := l.svcCtx.Redis().ZrangebyscoreWithScoresCtx(l.ctx, rediskey.NoticeConvMembersSubscribed(in.ConvId), min, time.Now().UnixMilli()+1000*60*60)
	if err != nil {
		if err == redis.Nil {
			return &pb.GetNoticeConvAllSubscribersResp{}, nil
		}
		l.Errorf("get conv subscribers error: %v", err)
		return &pb.GetNoticeConvAllSubscribersResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	userIds := make([]string, 0)
	for _, pair := range val {
		userId, _ := rediskey.ConvMembersSubscribedSplit(pair.Key)
		userIds = append(userIds, userId)
	}
	return &pb.GetNoticeConvAllSubscribersResp{
		UserIds: utils.Set(userIds),
	}, nil
}
