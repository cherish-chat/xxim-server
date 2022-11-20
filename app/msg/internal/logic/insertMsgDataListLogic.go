package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

func (l *InsertMsgDataListLogic) InsertMsgDataList(in *pb.MsgDataList) (*pb.MsgDataList, error) {
	var models []*msgmodel.Msg
	var err error
	var updateCacheIds []string
	var clientIds []string
	var clientIdModelMap = make(map[string]*msgmodel.Msg)
	var respMsgDataList = &pb.MsgDataList{MsgDataList: make([]*pb.MsgData, 0)}
	xtrace.StartFuncSpan(l.ctx, "GenModels", func(ctx context.Context) {
		for _, msgData := range in.MsgDataList {
			model := msgmodel.NewMsgFromPb(msgData)
			model.AutoConvId()
			if model.Options.StorageForServer {
				convId := model.ConvId
				// 给会话生成一个新的seq
				var seq int
				seq, err = IncrConvMaxSeq(l.svcCtx.Redis(), ctx, convId)
				if err != nil {
					return
				}
				model.SetSeq(int64(seq)).Check()
				models = append(models, model)
				updateCacheIds = append(updateCacheIds, model.ServerMsgId)
				clientIds = append(clientIds, model.ClientMsgId)
				clientIdModelMap[model.ClientMsgId] = model
			}
			respMsgDataList.MsgDataList = append(respMsgDataList.MsgDataList, model.ToMsgData())
		}
	})
	// 删除这些消息的缓存
	err = msgmodel.FlushMsgCache(l.ctx, l.svcCtx.Redis(), updateCacheIds)
	if err != nil {
		l.Errorf("InsertMsgDataList.DeleteCache err:%v", err)
		return respMsgDataList, err
	}
	// 只能单条修改消息 多条不支持
	var fs []func() error
	if len(clientIds) == 1 {
		fs = append(fs, func() error {
			// 判断是否有重复的clientMsgId
			var existServerIds []string
			var existModels []*msgmodel.Msg
			var err error
			xtrace.StartFuncSpan(l.ctx, "CheckClientMsgIdExist", func(ctx context.Context) {
				err = l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Where("clientMsgId = ?", clientIds[0]).Find(&existModels).Error
				if err != nil {
					l.Errorf("check clientMsgId exist failed, err: %v", err)
					return
				}
				for _, model := range existModels {
					existServerIds = append(existServerIds, model.ServerMsgId)
					updateCacheIds = append(updateCacheIds, model.ServerMsgId)
				}
				if len(existServerIds) > 0 {
					model := clientIdModelMap[clientIds[0]]
					// 更新已存在的clientMsgId的contentType 和 content 和 offlinePush 和 ext
					err = xorm.Update(l.svcCtx.Mysql(), &msgmodel.Msg{}, map[string]interface{}{
						"contentType": model.ContentType,
						"content":     model.Content,
						"offlinePush": model.OfflinePush,
						"ext":         model.Ext,
					}, xorm.Where("id in (?)", existServerIds))
				}
			})
			return err
		})
	}
	if len(models) > 0 {
		fs = append(fs, func() error {
			var err error
			xtrace.StartFuncSpan(l.ctx, "InsertManyMsg", func(ctx context.Context) {
				err = xorm.InsertMany(l.svcCtx.Mysql(), &msgmodel.Msg{}, models)
			})
			return err
		})
	}
	if len(fs) > 1 {
		err = mr.Finish(fs...)
		if err != nil {
			l.Errorf("mr.Finish error: %v", err)
			return respMsgDataList, err
		}
	} else if len(fs) == 1 {
		err = fs[0]()
		if err != nil {
			l.Errorf("mr.Finish error: %v", err)
			return respMsgDataList, err
		}
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		msgmodel.MsgFromMysql(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), updateCacheIds)
	}, nil)
	return respMsgDataList, nil
}
