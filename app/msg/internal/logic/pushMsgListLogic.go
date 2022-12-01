package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
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
		// 查询会话所有接受离线消息的成员 过滤在线用户 并推送
		xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "findAndPushOfflineMsgList", func(ctx context.Context) {
			go l.batchFindAndPushOfflineMsgList(ctx, convIdMsgListMap)
		}, propagation.MapCarrier{})
	}
	return &pb.CommonResp{}, nil
}

func (l *PushMsgListLogic) batchFindAndPushMsgList(listMap map[string]*pb.MsgDataList) {
	convIds := make([]string, 0)
	senders := make([]string, 0)
	for convId, msgDataList := range listMap {
		convIds = append(convIds, convId)
		for _, data := range msgDataList.MsgDataList {
			senders = append(senders, data.SenderId)
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
			} else {
				convSubscribers[convId] = &pb.GetConvSubscribersResp{}
				convUserIds[convId] = []string{}
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
			if len(userIds) == 0 {
				continue
			}
			msgDataListBytes, _ := proto.Marshal(msgDataList)
			resp, err := l.svcCtx.ImService().SendMsg(l.ctx, &pb.SendMsgReq{
				GetUserConnReq: &pb.GetUserConnReq{
					UserIds: userIds,
				},
				Event: pb.PushEvent_PushMsgDataList,
				Data:  msgDataListBytes,
			})
			if err != nil {
				l.Logger.Errorf("SendMsg err: %v", err)
				l.offlinePushMsgList(msgDataList, userIds)
				continue
			}
			successUserIdMap := make(map[string]bool)
			failedUserIds := make([]string, 0)
			for _, param := range resp.SuccessConnParams {
				successUserIdMap[param.UserId] = true
			}
			for _, param := range resp.FailedConnParams {
				if _, ok := successUserIdMap[param.UserId]; !ok {
					failedUserIds = append(failedUserIds, param.UserId)
				}
			}
			failedUserIds = utils.Set(failedUserIds)
			if len(failedUserIds) > 0 {
				l.offlinePushMsgList(msgDataList, failedUserIds)
			}
		}
	}
}

func (l *PushMsgListLogic) offlinePushMsgList(list *pb.MsgDataList, userIds []string) {
	if len(userIds) == 0 {
		return
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "OfflinePushMsgList", func(ctx context.Context) {
		for _, data := range list.MsgDataList {
			// 判断是否需要离线推送
			if !data.Options.OfflinePush {
				continue
			}
			convId := data.ConvId
			// 查询用户在此会话的离线推送设置
			if data.IsSingleConv() {
				l.offlinePushUser(ctx, data, convId, data.ReceiverUid())
			}
		}
	}, propagation.MapCarrier{})
}

func (l *PushMsgListLogic) batchFindAndPushOfflineMsgList(ctx context.Context, listMap map[string]*pb.MsgDataList) {
	for convId, msgDataList := range listMap {
		for _, data := range msgDataList.MsgDataList {
			if data.IsSingleConv() {
				// 单聊
				receiver := data.ReceiverUid()
				// 用户是否在线
				resp, err := l.svcCtx.ImService().GetUserConn(ctx, &pb.GetUserConnReq{
					UserIds: []string{receiver},
				})
				if err != nil {
					l.Errorf("GetUserConn err: %v", err)
					continue
				}
				if len(resp.ConnParams) > 0 {
					// 在线
					continue
				}
				l.offlinePushUser(ctx, data, convId, receiver)
			} else if data.IsGroupConv() {
				// 群聊
				receiver := data.ReceiverGid()
				// 查询群成员
				memberList, err := l.svcCtx.GroupService().GetGroupMemberList(ctx, &pb.GetGroupMemberListReq{
					GroupId: receiver,
					Page: &pb.Page{
						Page: 1,
						Size: 99999,
					},
					Filter: &pb.GetGroupMemberListReq_GetGroupMemberListFilter{NoDisturb: utils.AnyPtr(false)},
					Opt: &pb.GetGroupMemberListReq_GetGroupMemberListOpt{
						OnlyId:       utils.AnyPtr(true),
						GetNotifyOpt: utils.AnyPtr(true),
					},
				})
				if err != nil {
					l.Errorf("GetGroupMemberList err: %v", err)
					continue
				}
				if len(memberList.GroupMemberList) == 0 {
					continue
				}
				l.offlinePushGroup(ctx, data, convId, memberList.GroupMemberList...)
			}
		}
	}
}

func (l *PushMsgListLogic) offlinePushUser(ctx context.Context, data *pb.MsgData, convId string, userId string) {
	// 查询用户在此会话的离线推送设置
	convSetting, err := l.svcCtx.RelationService().GetSingleConvSetting(ctx, &pb.GetSingleConvSettingReq{
		ConvId: convId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("GetSingleMsgNotifyOpt err: %v", err)
		return
	}
	if convSetting.Setting.GetIsDisturb() {
		// 免打扰
		return
	}
	alert, content := data.OfflinePush.Title, data.OfflinePush.Content
	if !convSetting.Setting.GetNotifyPreview() {
		alert, content = l.svcCtx.SystemConfigMgr.GetOrDefault("offline_push_title", "惺惺线路"), l.svcCtx.SystemConfigMgr.GetOrDefault("offline_push_content", "您有一条新消息")
	}
	// 推送
	xtrace.StartFuncSpan(ctx, "PushOfflineMsg", func(ctx context.Context) {
		_, err := NewOfflinePushMsgLogic(ctx, l.svcCtx).OfflinePushMsg(&pb.OfflinePushMsgReq{
			UserIds:  []string{userId},
			Title:    alert,
			Content:  content,
			Payload:  "",
			UniqueId: fmt.Sprintf("%s:%s", convId, userId),
		})
		if err != nil {
			l.Errorf("OfflinePushMsg err: %v", err)
			return
		}
	})
}

func (l *PushMsgListLogic) offlinePushGroup(ctx context.Context, data *pb.MsgData, convId string, members ...*pb.GetGroupMemberListResp_GroupMember) {
	previewUids := make([]string, 0)
	noPreviewUids := make([]string, 0)
	for _, member := range members {
		if member.NotifyOpt.NoDisturb {
			continue
		}
		if member.NotifyOpt.Preview {
			previewUids = append(previewUids, member.UserId)
		} else {
			noPreviewUids = append(noPreviewUids, member.UserId)
		}
	}
	alert, content := data.OfflinePush.Title, data.OfflinePush.Content
	noPreviewAlert, noPreviewContent := l.svcCtx.SystemConfigMgr.GetOrDefault("offline_push_title", "惺惺线路"), l.svcCtx.SystemConfigMgr.GetOrDefault("offline_push_content", "您有一条新消息")
	// 推送
	xtrace.StartFuncSpan(ctx, "PushOfflineMsg", func(ctx context.Context) {
		if len(previewUids) > 0 {
			_, err := NewOfflinePushMsgLogic(ctx, l.svcCtx).OfflinePushMsg(&pb.OfflinePushMsgReq{
				UserIds:  previewUids,
				Title:    alert,
				Content:  content,
				Payload:  "",
				UniqueId: fmt.Sprintf("%s:%s", convId, data.ServerMsgId),
			})
			if err != nil {
				l.Errorf("OfflinePushMsg err: %v", err)
				return
			}
		}
		if len(noPreviewUids) > 0 {
			_, err := NewOfflinePushMsgLogic(ctx, l.svcCtx).OfflinePushMsg(&pb.OfflinePushMsgReq{
				UserIds:  noPreviewUids,
				Title:    noPreviewAlert,
				Content:  noPreviewContent,
				Payload:  "",
				UniqueId: fmt.Sprintf("%s:%s", convId, data.ServerMsgId),
			})
			if err != nil {
				l.Errorf("OfflinePushMsg err: %v", err)
				return
			}
		}
	})
}
