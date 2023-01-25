package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/propagation"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AfterConnect conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	return l.afterConnect(in)
}

func (l *AfterConnectLogic) afterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	// 查询用户所有订阅的会话 检测是否有未消费的消息 进行推送
	var convIds []string
	var err error
	xtrace.StartFuncSpan(l.ctx, "getAllConv", func(ctx context.Context) {
		convIds, err = l.getAllConv(in)
	})
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	// 并发查询未消费的消息
	var fs []func() error
	for _, convId := range convIds {
		convId := convId
		fs = append(fs, func() error {
			var err error
			xtrace.StartFuncSpan(l.ctx, "getUserNoticeData", func(ctx context.Context) {
				_, err = NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
					CommonReq: &pb.CommonReq{
						UserId: in.ConnParam.UserId,
					},
					UserId: in.ConnParam.UserId,
					ConvId: convId,
				})
			}, xtrace.StartFuncSpanWithCarrier(propagation.MapCarrier{
				"conv_id": convId,
			}))
			if err != nil {
				l.Errorf("get user notice data error: %v", err)
			}
			return err
		})
	}
	err = mr.Finish(fs...)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	return pb.NewSuccessResp(), nil
}

func (l *AfterConnectLogic) getAllConv(in *pb.AfterConnectReq) ([]string, error) {
	var userId = in.ConnParam.UserId
	var friendIds []string
	var groupIds []string
	var convIds []string
	// 获取用户订阅的好友列表
	{
		getFriendList, err := l.svcCtx.RelationService().GetFriendList(l.ctx, &pb.GetFriendListReq{
			CommonReq: &pb.CommonReq{
				UserId: userId,
			},
			Page: &pb.Page{
				Page: 1,
				Size: 0,
			},
			Opt: pb.GetFriendListReq_OnlyId,
		})
		if err != nil {
			l.Errorf("get friend list error: %v", err)
			return convIds, err
		}
		friendIds = getFriendList.Ids
		for _, id := range friendIds {
			convIds = append(convIds, pb.SingleConvId(userId, id))
		}
	}
	// 获取用户订阅的群组列表
	{
		getMyGroupList, err := l.svcCtx.GroupService().GetMyGroupList(l.ctx, &pb.GetMyGroupListReq{
			CommonReq: &pb.CommonReq{
				UserId: userId,
			},
			Page: &pb.Page{Page: 1},
			Filter: &pb.GetMyGroupListReq_Filter{
				FilterFold:   true,
				FilterShield: true,
			},
			Opt: pb.GetMyGroupListReq_ONLY_ID,
		})
		if err != nil {
			l.Errorf("get group list error: %v", err)
			return convIds, err
		}
		groupIds = getMyGroupList.Ids
		for _, id := range groupIds {
			convIds = append(convIds, pb.GroupConvId(id))
		}
	}
	return convIds, nil
}
