package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMsgListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgListLogic {
	return &PushMsgListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMsgListLogic) PushMsgList(in *pb.PushMsgListReq) (*pb.CommonResp, error) {
	var convIdMsgListMap = make(map[string]*pb.MsgDataList)
	for _, msgData := range in.MsgDataList {
		if msgData.ConvId == "" {
			if msgData.Receiver.UserId != nil {
				msgData.ConvId = msgmodel.SingleConvId(msgData.Sender, *msgData.Receiver.UserId)
			} else {
				msgData.ConvId = *msgData.Receiver.GroupId
			}
		}
		if utils.AnyToInt64(msgData.ServerTime) == 0 {
			msgData.ServerTime = utils.AnyToString(time.Now().UnixMilli())
		}
		if msgDataList, ok := convIdMsgListMap[msgData.ConvId]; ok {
			msgDataList.MsgDataList = append(msgDataList.MsgDataList, msgData)
		} else {
			convIdMsgListMap[msgData.ConvId] = &pb.MsgDataList{MsgDataList: []*pb.MsgData{msgData}}
		}
	}
	if len(convIdMsgListMap) > 0 {
		// 查询会话的订阅者并推送
		l.batchFindAndPushMsgList(convIdMsgListMap)
	}
	return &pb.CommonResp{}, nil
}

func (l *PushMsgListLogic) batchFindAndPushMsgList(listMap map[string]*pb.MsgDataList) {
	convIds := make([]string, 0)
	senders := make([]string, 0)
	for convId, msgDataList := range listMap {
		convIds = append(convIds, convId)
		for _, data := range msgDataList.MsgDataList {
			senders = append(senders, data.Sender)
		}
	}
	var convSubscribers = make(map[string]*pb.GetConvSubscribersResp)
	convUserIds := make(map[string][]string)
	xtrace.StartFuncSpan(l.ctx, "BatchGetConvSubscribers", func(ctx context.Context) {
		for _, convId := range convIds {
			subscribers, err := NewGetConvSubscribersLogic(ctx, l.svcCtx).GetConvSubscribers(&pb.GetConvSubscribersReq{
				ConvId:         convId,
				LastActiveTime: utils.AnyPtr(time.Now().Add(-time.Minute * 15).UnixMilli()),
			})
			if err != nil {
				l.Logger.Errorf("BatchGetConvSubscribers err: %v", err)
				continue
			}
			if len(subscribers.UserIdList) > 0 {
				convSubscribers[convId] = subscribers
				convUserIds[convId] = subscribers.UserIdList
			}
		}
		for convId, subscribers := range convSubscribers {
			for _, subscriber := range subscribers.UserIdList {
				if _, ok := convUserIds[convId]; !ok {
					convUserIds[convId] = make([]string, 0)
				}
				convUserIds[convId] = append(convUserIds[convId], subscriber)
			}
		}
	})
	for convId, msgDataList := range listMap {
		if userIds, ok := convUserIds[convId]; ok {
			userIds = utils.Set(append(userIds, senders...))
			msgDataListBytes, _ := proto.Marshal(msgDataList)
			_, _ = l.svcCtx.ImService().SendMsg(l.ctx, &pb.SendMsgReq{
				GetUserConnReq: &pb.GetUserConnReq{
					UserIds: userIds,
				},
				Event: pb.PushEvent_PushMsgDataList,
				Data:  msgDataListBytes,
			})
		}
	}
}
