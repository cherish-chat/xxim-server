package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InsertMsgDataListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInsertMsgDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertMsgDataListLogic {
	return &InsertMsgDataListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InsertMsgDataListLogic) InsertMsgDataList(in *pb.MsgDataList) (*pb.CommonResp, error) {
	var models []*msgmodel.Msg
	var err error
	xtrace.StartFuncSpan(l.ctx, "InsertMsgDataList.GenModels", func(ctx context.Context) {
		for _, msgData := range in.MsgDataList {
			model := msgmodel.NewMsgFromPb(msgData)
			model.AutoConvId()
			convId := model.ConvId
			// 给会话生成一个新的seq
			k := rediskey.ConvKv(convId)
			var seq int
			seq, err = l.svcCtx.Redis().HincrbyCtx(l.ctx, k, rediskey.HKConvMaxSeq(), 1)
			if err != nil {
				return
			}
			model.SetSeq(int64(seq)).Check()
			models = append(models, model)
		}
	})
	if err != nil {
		l.Errorf("InsertMsgDataList.GenModels err:%v", err)
		return pb.NewRetryErrorResp(), err
	}
	_, err = l.svcCtx.Mongo().Collection(&msgmodel.Msg{}).InsertMany(l.ctx, models, opts.InsertManyOptions{
		InsertManyOptions: options.InsertMany().SetOrdered(true),
	})
	if err != nil {
		l.Errorf("mongo InsertMany error: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
