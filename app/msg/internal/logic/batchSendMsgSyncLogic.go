package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSendMsgSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchSendMsgSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSendMsgSyncLogic {
	return &BatchSendMsgSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchSendMsgSyncLogic) BatchSendMsgSync(in *pb.BatchSendMsgReq) (*pb.CommonResp, error) {
	msg := msgmodel.NewMsgFromPb(in.MsgData)
	model := &msgmodel.BatchMsg{
		Id:          utils.GenId(),
		Msg:         msg.Check(),
		UserIdList:  utils.AnyMakeSlice(in.UserIdList),
		GroupIdList: utils.AnyMakeSlice(in.GroupIdList),
	}
	var err error
	xtrace.StartFuncSpan(l.ctx, "BatchSendMsgSync.InsertOne", func(ctx context.Context) {
		_, err = l.svcCtx.Mongo().Collection(model).InsertOne(l.ctx, model)
	})
	if err != nil {
		l.Errorf("BatchSendMsgSync error: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	xtrace.StartFuncSpan(l.ctx, "BatchSendMsgSync.MHSetLua", func(ctx context.Context) {
		var kvs []xmgo.MHSetKv
		for _, userId := range in.UserIdList {
			convId := msgmodel.SingleConvId(in.MsgData.Sender, userId)
			// 给会话生成一个新的seq
			k := rediskey.ConvKv(convId)
			var seq int
			seq, err = l.svcCtx.Redis().HincrbyCtx(l.ctx, k, rediskey.HKConvMaxSeq(), 1)
			if err != nil {
				return
			}
			msgId := msgmodel.ServerMsgId(convId, int64(seq))
			kvs = append(kvs, xmgo.MHSetKv{
				Key: rediskey.ConvMsgIdMapping(convId),
				HK:  msgId,
				V:   model.Id,
			})
		}
		for _, groupId := range in.GroupIdList {
			convId := groupId
			// 给会话生成一个新的seq
			k := rediskey.ConvKv(convId)
			var seq int
			seq, err = l.svcCtx.Redis().HincrbyCtx(l.ctx, k, rediskey.HKConvMaxSeq(), 1)
			if err != nil {
				l.Errorf("redis Hincrby error: %v", err)
				return
			}
			msgId := msgmodel.ServerMsgId(convId, int64(seq))
			kvs = append(kvs, xmgo.MHSetKv{
				Key: rediskey.ConvMsgIdMapping(convId),
				HK:  msgId,
				V:   model.Id,
			})
		}
		err = xmgo.MHSet(l.svcCtx.Mongo().Collection(&xmgo.MHSetKv{}), l.ctx, kvs...)
	})
	if err != nil {
		l.Errorf("redis MHSetLua error: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	// TODO 推送
	return &pb.CommonResp{}, nil
}
