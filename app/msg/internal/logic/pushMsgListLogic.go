package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
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
		if msgData.Receiver.UserId != nil {
			msgData.ConvId = msgmodel.SingleConvId(msgData.Sender, *msgData.Receiver.UserId)
		} else {
			msgData.ConvId = *msgData.Receiver.GroupId
		}
		msgData.ServerTime = time.Now().UnixMilli()
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
	for convId := range listMap {
		convIds = append(convIds, convId)
	}
	var convSubscribers = make(map[string]*MockSubscribers)
	convUserIds := make(map[string][]string)
	xtrace.StartFuncSpan(l.ctx, "BatchGetConvSubscribers", func(ctx context.Context) {
		convSubscribers = mockBatchGetConvSubscribers(convIds)
		for convId, subscribers := range convSubscribers {
			for _, subscriber := range subscribers.Subscribers {
				if _, ok := convUserIds[convId]; !ok {
					convUserIds[convId] = make([]string, 0)
				}
				convUserIds[convId] = append(convUserIds[convId], subscriber.UserId)
			}
		}
	})
	for convId, msgDataList := range listMap {
		if userIds, ok := convUserIds[convId]; ok {
			msgDataListBytes, _ := proto.Marshal(msgDataList)
			for _, pod := range l.svcCtx.ConnPodsMgr.AllConnServices() {
				_, err := pod.SendMsg(l.ctx, &pb.SendMsgReq{
					GetUserConnReq: &pb.GetUserConnReq{
						UserIds: userIds,
					},
					Event: pb.PushEvent_PushMsgDataList,
					Data:  msgDataListBytes,
				})
				if err != nil {
					l.Errorf("SendMsg error: %v", err)
				}
			}
		}
	}
}
