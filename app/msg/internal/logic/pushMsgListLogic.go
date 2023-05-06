package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
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
		// 批量查询会话的订阅者并推送 如果是单聊没推送成功 会直接走第三方离线推送逻辑
		l.batchFindAndPushMsgList(convIdMsgListMap)
		// 查询群所有接受离线消息的成员 过滤在线用户 并推送
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "findAndPushOfflineMsgList", func(ctx context.Context) {
			l.batchFindGroupMemberAndPushOfflineMsgList(ctx, convIdMsgListMap)
		}, propagation.MapCarrier{})
	}
	return &pb.CommonResp{}, nil
}

// batchFindAndPushMsgList 批量查询会话的订阅者并推送 如果是单聊没推送成功 会直接走第三方离线推送逻辑
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
			if pb.IsSingleConv(convId) {
				userIds := pb.ParseSingleConv(convId)
				convSubscribers[convId] = &pb.GetConvSubscribersResp{
					UserIdList: userIds,
				}
				convUserIds[convId] = userIds
			} else {
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
				l.offlinePushMsgListSingle(msgDataList, []string{})
				continue
			}
			successUserIdMap := make(map[string]bool)
			for _, param := range resp.SuccessConnParams {
				successUserIdMap[param.UserId] = true
			}
			successUserIds := make([]string, 0)
			for uid := range successUserIdMap {
				successUserIds = append(successUserIds, uid)
			}
			l.Debugf("SendMsg successUserIds: %v", successUserIds)
			l.offlinePushMsgListSingle(msgDataList, successUserIds)
		}
	}
}

func (l *PushMsgListLogic) offlinePushMsgListSingle(list *pb.MsgDataList, excludeUserIds []string) {
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "OfflinePushMsgList", func(ctx context.Context) {
		for _, data := range list.MsgDataList {
			// 判断是否需要离线推送
			if !data.Options.OfflinePush {
				continue
			}
			convId := data.ConvId
			// 查询用户在此会话的离线推送设置
			if data.IsSingleConv() {
				receiverUid := data.ReceiverUid()
				// 如果 receiverUid 在 excludeUserIds 中 则不推送
				if !utils.InSlice(excludeUserIds, receiverUid) {
					l.offlinePushUser(ctx, data, convId, receiverUid)
				}
			}
		}
	}, propagation.MapCarrier{})
}

func (l *PushMsgListLogic) batchFindGroupMemberAndPushOfflineMsgList(ctx context.Context, listMap map[string]*pb.MsgDataList) {
	for convId, msgDataList := range listMap {
		for _, data := range msgDataList.MsgDataList {
			if !data.GetOptions().GetOfflinePush() {
				continue
			}
			if data.IsSingleConv() {
				continue
			}
			if data.IsGroupConv() {
				// 查询群成员
				var memberList []*immodel.ConvSetting
				var err error
				if l.svcCtx.ConfigMgr.OfflinePushAllowDisturb(ctx, "") {
					memberList, err = immodel.SearchGroupMemberList(l.svcCtx.Mysql(), data.ConvId, 1000, map[string]interface{}{
						"isDisturb": false,
					})
				} else {
					groupMemberList, err := l.svcCtx.GroupService().GetGroupMemberList(ctx, &pb.GetGroupMemberListReq{
						CommonReq: &pb.CommonReq{},
						GroupId:   pb.ParseGroupConv(data.ConvId),
						Page:      &pb.Page{Size: 1000, Page: 1},
						Filter: &pb.GetGroupMemberListReq_GetGroupMemberListFilter{
							NoDisturb:  nil,
							OnlyOwner:  nil,
							OnlyAdmin:  nil,
							OnlyMember: nil,
						},
						Opt: &pb.GetGroupMemberListReq_GetGroupMemberListOpt{OnlyId: utils.AnyPtr(true)},
					})
					if err != nil {
						l.Errorf("GetGroupMemberList err: %v", err)
						continue
					}
					if len(groupMemberList.GroupMemberList) == 0 {
						continue
					}
					for _, info := range groupMemberList.GroupMemberList {
						memberList = append(memberList, &immodel.ConvSetting{
							ConvId:            data.ConvId,
							UserId:            info.MemberId,
							IsTop:             false,
							IsDisturb:         false,
							NotifyPreview:     true,
							NotifySound:       true,
							NotifyCustomSound: "",
							NotifyVibrate:     true,
							IsShield:          false,
							ChatBg:            "",
						})
					}
				}
				if err != nil {
					l.Errorf("GetGroupMemberList err: %v", err)
					continue
				}
				if len(memberList) == 0 {
					continue
				}
				l.offlinePushGroup(ctx, data, convId, memberList...)
			}
		}
	}
}

func (l *PushMsgListLogic) offlinePushUser(ctx context.Context, data *pb.MsgData, convId string, userId string) {
	// 查询用户在此会话的离线推送设置
	convSettings, err := l.svcCtx.ImService().GetConvSetting(ctx, &pb.GetConvSettingReq{
		ConvIds:   []string{convId},
		CommonReq: &pb.CommonReq{UserId: userId},
	})
	if err != nil {
		l.Errorf("GetSingleMsgNotifyOpt err: %v", err)
		return
	}
	if len(convSettings.ConvSettings) == 0 {
		convSettings.ConvSettings = append(convSettings.ConvSettings, &pb.ConvSetting{
			UserId:            userId,
			ConvId:            convId,
			IsTop:             nil,
			IsDisturb:         nil,
			NotifyPreview:     utils.AnyPtr(true),
			NotifySound:       utils.AnyPtr(true),
			NotifyCustomSound: nil,
			NotifyVibrate:     utils.AnyPtr(true),
			IsShield:          nil,
			ChatBg:            nil,
		})
	}
	convSetting := convSettings.ConvSettings[0]
	l.Debugf("offlinePushUser convSetting: %v", utils.AnyToString(convSetting))
	if l.svcCtx.ConfigMgr.OfflinePushAllowDisturb(l.ctx, data.SenderId) && convSetting.GetIsDisturb() {
		// 免打扰
		return
	}
	alert, content := data.OfflinePush.Title, data.OfflinePush.Content
	if !convSetting.GetNotifyPreview() {
		alert, content = l.svcCtx.ConfigMgr.OfflinePushTitle(l.ctx, data.SenderId), l.svcCtx.ConfigMgr.OfflinePushContent(l.ctx, data.SenderId)
	}
	// 推送
	xtrace.StartFuncSpan(ctx, "PushOfflineMsg", func(ctx context.Context) {
		l.Debugf("PushOfflineMsg userId: %s, alert: %s, content: %s", userId, alert, content)
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

func (l *PushMsgListLogic) offlinePushGroup(ctx context.Context, data *pb.MsgData, convId string, members ...*immodel.ConvSetting) {
	previewUids := make([]string, 0)
	noPreviewUids := make([]string, 0)
	for _, member := range members {
		if member.IsDisturb {
			continue
		}
		if member.NotifyPreview {
			previewUids = append(previewUids, member.UserId)
		} else {
			noPreviewUids = append(noPreviewUids, member.UserId)
		}
	}
	alert, content := data.OfflinePush.Title, data.OfflinePush.Content
	noPreviewAlert, noPreviewContent := l.svcCtx.ConfigMgr.OfflinePushTitle(l.ctx, data.SenderId), l.svcCtx.ConfigMgr.OfflinePushContent(l.ctx, data.SenderId)
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
