package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
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
	var notices []*noticemodel.Notice

	xtrace.StartFuncSpan(l.ctx, "GenModels", func(ctx context.Context) {
		for _, msgData := range in.MsgDataList {
			model := msgmodel.NewMsgFromPb(msgData)
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

			if len(model.AtUsers) > 0 {
				notices = append(notices, &noticemodel.Notice{
					ConvId: pb.HiddenConvId(model.ConvId),
					Options: noticemodel.NoticeOption{
						StorageForClient: false,
						UpdateConvNotice: false,
					},
					ContentType: int32(pb.NoticeType_AT),
					Content: utils.AnyToBytes(map[string]any{
						"atUsers":     model.AtUsers,
						"convId":      model.ConvId,
						"serverMsgId": model.ServerMsgId,
						"clientMsgId": model.ClientMsgId,
						"seq":         model.Seq,
					}),
					UniqueId: model.ServerMsgId,
					Title:    "",
					Ext:      nil,
				})
			}
		}
	})
	// 删除这些消息的缓存
	err = msgmodel.FlushMsgCache(l.ctx, l.svcCtx.Redis(), updateCacheIds)
	if err != nil {
		l.Errorf("InsertMsgDataList.DeleteCache err:%v", err)
		return respMsgDataList, err
	}
	// 只能单条修改消息 多条不支持
	if len(models) > 0 {
		{
			var err error
			xtrace.StartFuncSpan(l.ctx, "InsertManyMsg", func(ctx context.Context) {
				err = l.svcCtx.Mysql().Transaction(func(tx *gorm.DB) error {
					err = msgmodel.InsertManyMsg(ctx, tx, models)
					if err != nil {
						l.Errorf("InsertMsgDataList.InsertManyMsg err:%v", err)
						return err
					}
					if len(notices) > 0 {
						// 插入notices
						for _, notice := range notices {
							err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
							if err != nil {
								l.Errorf("insert notice err: %v", err)
								return err
							}
						}
					}
					return nil
				})

			})
			if err != nil {
				return respMsgDataList, err
			}
		}
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		msgmodel.MsgFromMysql(ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), updateCacheIds)
	}, nil)

	// 通知
	if len(notices) > 0 {
		xtrace.StartFuncSpan(l.ctx, "SendNotice", func(ctx context.Context) {
			utils.RetryProxy(ctx, 12, time.Second, func() error {
				for _, notice := range notices {
					_, err := l.svcCtx.NoticeService().GetUserNoticeData(ctx, &pb.GetUserNoticeDataReq{
						CommonReq: &pb.CommonReq{},
						UserId:    "",
						ConvId:    notice.ConvId,
					})
					if err != nil {
						l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
						return err
					}
				}
				return nil
			})
		})
	}
	return respMsgDataList, nil
}
