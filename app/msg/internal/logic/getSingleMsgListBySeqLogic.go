package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.mongodb.org/mongo-driver/bson"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSingleMsgListBySeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSingleMsgListBySeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSingleMsgListBySeqLogic {
	return &GetSingleMsgListBySeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetSingleMsgListBySeq 通过seq拉取一个单聊会话的消息
func (l *GetSingleMsgListBySeqLogic) GetSingleMsgListBySeq(in *pb.GetSingleMsgListBySeqReq) (*pb.GetSingleMsgListBySeqResp, error) {
	if len(in.SeqList) == 0 {
		return &pb.GetSingleMsgListBySeqResp{}, nil
	}
	// 会话id
	convId := msgmodel.SingleConvId(in.Requester.Id, in.UserId)
	// 组成想要查询的 id 列表
	expectIds := make([]string, 0)
	for _, seq := range in.SeqList {
		expectIds = append(expectIds, msgmodel.ServerMsgId(convId, seq))
	}
	// 查询
	var msgList []*msgmodel.Msg
	var err error
	xtrace.StartFuncSpan(l.ctx, "FindMsgByIds", func(ctx context.Context) {
		err = l.svcCtx.Mongo().Collection(&msgmodel.Msg{}).Find(l.ctx, bson.M{
			"_id": bson.M{"$in": expectIds},
		}).All(&msgList)
	})
	if err != nil {
		l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
		return &pb.GetSingleMsgListBySeqResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var msgMap = make(map[string]*msgmodel.Msg)
	for _, msg := range msgList {
		msgMap[msg.ServerMsgId] = msg
	}
	var notFoundIds []string
	for _, id := range expectIds {
		if _, ok := msgMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		xtrace.StartFuncSpan(l.ctx, "FindMsgByIdsFromBatchMsg", func(ctx context.Context) {
			var kvs []xmgo.MHSetKv
			for _, id := range notFoundIds {
				kvs = append(kvs, xmgo.MHSetKv{
					Key: rediskey.ConvMsgIdMapping(convId),
					HK:  id,
					V:   nil,
				})
			}
			var results []*xmgo.MHSetKv
			results, err = xmgo.MHGet(l.svcCtx.Mongo().Collection(&xmgo.MHSetKv{}), l.ctx, kvs...)
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
				err = l.svcCtx.Mongo().Collection(&msgmodel.BatchMsg{}).Find(l.ctx, bson.M{
					"_id": bson.M{"$in": batchMsgIds},
				}).All(&batchMsgList)
				if err != nil {
					l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
					return
				}
				for _, batchMsg := range batchMsgList {
					msg := batchMsg.Msg
					msg.ConvId = convId
					if msg.Sender == in.Requester.Id {
						msg.Receiver.UserId = in.UserId
					} else {
						msg.Receiver.UserId = in.Requester.Id
					}
					msg.ServerMsgId = batchIdMsgIdMap[batchMsg.Id]
					_, msg.Seq = msgmodel.ParseSingleServerMsgId(msg.ServerMsgId)
					msgList = append(msgList, msg)
					msgMap[msg.ServerMsgId] = msg
				}
			}
		})
		if err != nil {
			l.Errorf("GetSingleMsgListBySeq failed, err: %v", err)
			return &pb.GetSingleMsgListBySeqResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	for _, id := range expectIds {
		if _, ok := msgMap[id]; !ok {
			_, seq := msgmodel.ParseSingleServerMsgId(id)
			nullMsg := msgmodel.NewNullMsg(convId, seq)
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
	return &pb.GetSingleMsgListBySeqResp{MsgDataList: resp}, nil
}
