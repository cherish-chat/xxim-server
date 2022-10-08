package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"
)

type ConsumerStorage struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumerStorage(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumerStorage {
	return &ConsumerStorage{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConsumerStorage) Consume() {
	go func() {
		err := l.svcCtx.StorageConsumer().Consume(l.ctx, l.consumeMsg)
		if err != nil {
			l.Errorf("consume msg from pulsar failed, err: %v", err)
		}
	}()
}

func (l *ConsumerStorage) consumeMsg(ctx context.Context, _ string, _ string, payload []byte) error {
	msg := &pb.MsgToMQData{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		l.Errorf("unmarshal msg failed, err: %v", err)
		return err
	}
	var msgDataList []*dbmodel.Msg
	var needRewriteMsgIds []string
	nowMilli := utils.GetNowMilli()
	for _, msgData := range msg.MsgDataList {
		for _, conv := range msgData.ConvList {
			if msgData.ServerMsgId == "" {
				msgData.ServerMsgId = utils.GenId()
			}
			if msgData.ClientMsgId == "" {
				msgData.ClientMsgId = utils.GenId()
			}
			if msgData.ServerTime < 1000000000000 {
				msgData.ServerTime = nowMilli
			}
			if msgData.Seq == 0 {
				seq, err := l.svcCtx.Redis().Incr(ctx, rediskey.ConvMaxSeq(conv.Id)).Result()
				if err != nil {
					l.Errorf("incr conv max seq failed, err: %v", err)
					return err
				}
				msgData.Seq = uint32(seq)
			}
			if msgData.OfflinePush == nil {
				msgData.OfflinePush = &pb.MsgData_OfflinePush{}
			}
			if msgData.MsgOptions == nil {
				msgData.MsgOptions = &pb.MsgData_MsgOptions{}
			}
			if len(msgData.ExcludeUIds) == 0 {
				msgData.ExcludeUIds = make([]string, 0)
			}
			model := &dbmodel.Msg{
				Id:          msgData.ServerMsgId,
				ClientMsgId: msgData.ClientMsgId,
				ConvId:      conv.Id,
				ConvInfo:    conv.Info,
				SenderId:    msgData.SenderId,
				SenderInfo:  msgData.SenderInfo,
				ClientTime:  msgData.ClientTime,
				ServerTime:  msgData.ServerTime,
				Seq:         msgData.Seq,
				ContentType: msgData.ContentType,
				OfflinePush: msgData.OfflinePush,
				MsgOptions:  msgData.MsgOptions,
				Ex:          dbmodel.NewMsgEx(msgData.Ex),
				ExcludeUIds: msgData.ExcludeUIds,
				DeletedAt:   0,
			}
			msgDataList = append(msgDataList, model)
			if msgData.MsgOptions.Rewrite {
				needRewriteMsgIds = append(needRewriteMsgIds, model.ClientMsgId)
			}
		}
	}
	if len(msgDataList) == 0 {
		// 没有消息
		return nil
	}
	if len(needRewriteMsgIds) > 0 {
		// 修改原消息
		_, err := l.svcCtx.MsgCollection().UpdateAll(ctx, bson.M{
			"clientMsgId": bson.M{
				"$in": needRewriteMsgIds,
			},
		}, bson.M{
			"$set": bson.M{
				"deletedAt": nowMilli,
			},
		})
		if err != nil {
			l.Errorf("update msg failed, err: %v", err)
			return err
		}
	}
	// 批量插入消息
	_, err = l.svcCtx.MsgCollection().InsertMany(ctx, msgDataList)
	if err != nil {
		l.Errorf("insert msg failed, err: %v", err)
		return err
	}
	go l.afterConsumeMsg(ctx, msg, msgDataList)
	return nil
}

func (l *ConsumerStorage) afterConsumeMsg(ctx context.Context, msg *pb.MsgToMQData, list []*dbmodel.Msg) {

}
