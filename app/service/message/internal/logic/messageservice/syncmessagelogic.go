package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/message/messagemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"sort"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMessageLogic {
	return &SyncMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SyncMessage 同步消息
func (l *SyncMessageLogic) SyncMessage(in *peerpb.SyncMessageReq) (*peerpb.SyncMessageResp, error) {
	var seqList = in.SeqList
	var ids = make([]string, 0, len(seqList))
	for _, seq := range seqList {
		ids = append(ids, messagemodel.GenerateMessageId(in.ConvId, in.ConvType, seq))
	}
	var results []*messagemodel.Message
	err := l.svcCtx.MessageCollection.Find(context.Background(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}).All(&results)
	if err != nil {
		l.Errorf("SyncMessageLogic.SyncMessage err: %v", err)
		return nil, err
	}
	resultsMap := make(map[uint32]*messagemodel.Message)
	for _, result := range results {
		resultsMap[result.Seq] = result
	}
	for _, seq := range in.SeqList {
		if _, ok := resultsMap[seq]; !ok {
			// 补充消息
			resultsMap[seq] = &messagemodel.Message{
				MessageId:        utils.Snowflake.String(),
				Uuid:             utils.Snowflake.String(),
				ConversationId:   in.ConvId,
				ConversationType: in.ConvType,
				Sender:           messagemodel.MessageSender{},
				Content:          nil,
				ContentType:      0,
				SendTime:         0,
				InsertTime:       0,
				Seq:              seq,
				Option: messagemodel.MessageOptions{
					StorageForServer: true,
					StorageForClient: true,
					NeedDecrypt:      false,
					CountUnread:      false,
				},
				ExtraMap: nil,
			}
		}
	}
	var finalResults []*peerpb.Message
	for _, seq := range in.SeqList {
		finalResults = append(finalResults, resultsMap[seq].ToPb())
	}
	// seq正序
	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Seq < finalResults[j].Seq
	})
	return &peerpb.SyncMessageResp{
		ConvId:   in.ConvId,
		ConvType: in.ConvType,
		Messages: finalResults,
	}, nil
}
