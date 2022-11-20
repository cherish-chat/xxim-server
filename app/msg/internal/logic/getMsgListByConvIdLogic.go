package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"sort"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMsgListByConvIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMsgListByConvIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMsgListByConvIdLogic {
	return &GetMsgListByConvIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMsgListByConvIdLogic) fromRedis(ids []string) (msgList []*msgmodel.Msg, err error) {
	// 从redis中获取
	redisMsgList, err := l.svcCtx.Redis().MgetCtx(l.ctx, utils.UpdateSlice(ids, func(v string) string {
		return rediskey.MsgKey(v)
	})...)
	if err != nil {
		l.Errorf("redis MGet error: %v", err)
		return nil, err
	}
	for i, redisMsg := range redisMsgList {
		msg := &msgmodel.Msg{}
		if redisMsg == xredis.NotFound || redisMsg == "" {
			id := ids[i]
			msg.NotFound(id)
		} else {
			err = json.Unmarshal([]byte(redisMsg), msg)
			if err != nil {
				l.Errorf("msg Unmarshal error: %v redisMsg: %s", err, redisMsg)
				continue
			}
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}

func (l *GetMsgListByConvIdLogic) proxyGetMsgListByIds(ids []string) (msgList []*msgmodel.Msg, err error) {
	msgs, err := l.fromRedis(ids)
	if err != nil {
		return msgmodel.MsgFromMysql(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), ids)
	}
	// 判断是否有缺失
	msgMap := make(map[string]*msgmodel.Msg)
	for _, msg := range msgs {
		msgMap[msg.ServerMsgId] = msg
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := msgMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		// 从 mysql 中获取
		mysqlMsgs, err := msgmodel.MsgFromMysql(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), notFoundIds)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, mysqlMsgs...)
	}
	return msgs, nil
}

// GetMsgListByConvId 通过seq拉取一个会话的消息
func (l *GetMsgListByConvIdLogic) GetMsgListByConvId(in *pb.GetMsgListByConvIdReq) (*pb.GetMsgListResp, error) {
	if len(in.SeqList) == 0 {
		return &pb.GetMsgListResp{}, nil
	}
	// 会话id
	convId := in.ConvId
	// 组成想要查询的 id 列表
	expectIds := make([]string, 0)
	for _, seq := range in.SeqList {
		expectIds = append(expectIds, msgmodel.ServerMsgId(convId, utils.AnyToInt64(seq)))
	}
	// 查询
	var msgList []*msgmodel.Msg
	var err error
	xtrace.StartFuncSpan(l.ctx, "FindMsgByIds", func(ctx context.Context) {
		msgList, err = l.proxyGetMsgListByIds(expectIds)
	})
	if err != nil {
		l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
		return &pb.GetMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var msgMap = make(map[string]*msgmodel.Msg)
	for _, msg := range msgList {
		msgMap[msg.ServerMsgId] = msg
	}
	var notFoundIds []string
	for _, id := range expectIds {
		if m, ok := msgMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		} else if m.IsNotFound() {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		xtrace.StartFuncSpan(l.ctx, "FindMsgByIdsFromBatchMsg", func(ctx context.Context) {
			var kvs []xorm.HashKv
			for _, id := range notFoundIds {
				kvs = append(kvs, xorm.HashKv{
					Key: rediskey.ConvMsgIdMapping(convId),
					HK:  id,
					V:   "",
				})
			}
			var results []*xorm.HashKv
			results, err = xorm.MHGet(l.svcCtx.Mysql(), kvs...)
			if err != nil {
				l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
				return
			}
			batchMsgIds := make([]string, 0)
			batchIdMsgIdMap := make(map[string]string)
			for _, result := range results {
				batchMsgIds = append(batchMsgIds, utils.AnyToString(result.V))
				batchIdMsgIdMap[utils.AnyToString(result.V)] = result.HK
			}
			if len(batchMsgIds) > 0 {
				var batchMsgList []*msgmodel.BatchMsg
				err = l.svcCtx.Mysql().Model(&msgmodel.BatchMsg{}).Where("id in (?)", batchMsgIds).Find(&batchMsgList).Error
				if err != nil {
					l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
					return
				}
				if msgmodel.ConvIsGroup(convId) {
					for _, batchMsg := range batchMsgList {
						msg := batchMsg.Msg
						msg.ConvId = convId
						msg.Receiver.GroupId = convId
						msg.ServerMsgId = batchIdMsgIdMap[batchMsg.Id]
						msg.ClientMsgId = msgmodel.BatchMsgClientMsgId(msg.ClientMsgId, convId)
						_, msg.Seq = msgmodel.ParseGroupServerMsgId(msg.ServerMsgId)
						msgList = append(msgList, msg)
						msgMap[msg.ServerMsgId] = msg
					}
				} else {
					for _, batchMsg := range batchMsgList {
						msg := batchMsg.Msg
						msg.ConvId = convId
						if msg.Sender == in.Requester.Id {
							id1, id2 := msgmodel.ParseSingleConvId(convId)
							msg.Receiver.UserId = utils.If(id1 == in.Requester.Id, id2, id1)
						} else {
							msg.Receiver.UserId = in.Requester.Id
						}
						msg.ServerMsgId = batchIdMsgIdMap[batchMsg.Id]
						msg.ClientMsgId = msgmodel.BatchMsgClientMsgId(msg.ClientMsgId, convId)
						_, msg.Seq = msgmodel.ParseSingleServerMsgId(msg.ServerMsgId)
						msgList = append(msgList, msg)
						msgMap[msg.ServerMsgId] = msg
					}
				}
			}
		})
		if err != nil {
			l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
			return &pb.GetMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	for _, id := range expectIds {
		if _, ok := msgMap[id]; !ok {
			_, seq := msgmodel.ParseSingleServerMsgId(id)
			nullMsg := &msgmodel.Msg{}
			nullMsg.NotFound(msgmodel.ServerMsgId(convId, seq))
			msgList = append(msgList, nullMsg)
			msgMap[id] = nullMsg
		}
	}
	// seq正序排序
	xtrace.StartFuncSpan(l.ctx, "SortMsgList", func(ctx context.Context) {
		sort.Slice(msgList, func(i, j int) bool {
			return msgList[i].Seq < msgList[j].Seq
		})
	})
	var resp []*pb.MsgData
	for _, msg := range msgList {
		resp = append(resp, msg.ToMsgData())
	}
	if !in.Push {
		return &pb.GetMsgListResp{MsgDataList: resp}, nil
	} else {
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "PushMsgList", func(ctx context.Context) {
			msgDataListBytes, _ := proto.Marshal(&pb.MsgDataList{MsgDataList: resp})
			_, _ = l.svcCtx.ImService().SendMsg(ctx, &pb.SendMsgReq{
				GetUserConnReq: &pb.GetUserConnReq{
					UserIds: []string{in.Requester.Id},
					Devices: []string{in.Requester.DeviceId},
				},
				Event: pb.PushEvent_PushMsgDataList,
				Data:  msgDataListBytes,
			})
		}, nil)
		return &pb.GetMsgListResp{}, nil
	}
}
