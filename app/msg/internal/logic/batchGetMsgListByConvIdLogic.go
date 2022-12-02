package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"sort"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetMsgListByConvIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetMsgListByConvIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetMsgListByConvIdLogic {
	return &BatchGetMsgListByConvIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetMsgListByConvIdLogic) fromRedis(ids []string) (msgList []*msgmodel.Msg, err error) {
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
		if redisMsg == xredis.NotFound {
			id := ids[i]
			msg.NotFound(id)
		} else if redisMsg != "" {
			err = json.Unmarshal([]byte(redisMsg), msg)
			if err != nil {
				l.Errorf("msg Unmarshal error: %v redisMsg: %s", err, redisMsg)
				continue
			}
			msgList = append(msgList, msg)
		}
	}
	return msgList, nil
}

func (l *BatchGetMsgListByConvIdLogic) proxyGetMsgListByIds(ids []string) (msgList []*msgmodel.Msg, err error) {
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

// BatchGetMsgListByConvId 通过seq拉取一个会话的消息
func (l *BatchGetMsgListByConvIdLogic) BatchGetMsgListByConvId(in *pb.BatchGetMsgListByConvIdReq) (*pb.GetMsgListResp, error) {
	expectIds := make([]string, 0)
	for _, item := range in.Items {
		convId := item.ConvId
		for _, seq := range item.SeqList {
			expectIds = append(expectIds, pb.ServerMsgId(convId, utils.AnyToInt64(seq)))
		}
	}
	if len(expectIds) == 0 {
		return &pb.GetMsgListResp{}, nil
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
	for _, id := range expectIds {
		if _, ok := msgMap[id]; !ok {
			convId, seq := pb.ParseConvServerMsgId(id)
			nullMsg := &msgmodel.Msg{}
			nullMsg.NotFound(pb.ServerMsgId(convId, seq))
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
					UserIds: []string{in.CommonReq.UserId},
					Devices: []string{in.CommonReq.DeviceId},
				},
				Event: pb.PushEvent_PushMsgDataList,
				Data:  msgDataListBytes,
			})
		}, nil)
		return &pb.GetMsgListResp{}, nil
	}
}
