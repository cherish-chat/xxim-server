package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgListAsyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgListAsyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgListAsyncLogic {
	return &SendMsgListAsyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgListAsyncLogic) check(in *pb.SendMsgListReq) (*pb.SendMsgListResp, error) {
	for _, msgData := range in.MsgDataList {
		// 检查消息类型
		switch pb.ContentType(msgData.ContentType) {
		case pb.ContentType_TEXT:
			resp, err := l.checkText(msgData)
			if err != nil {
				return resp, err
			} else if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
				return resp, nil
			}
		}
	}
	return &pb.SendMsgListResp{}, nil
}

func (l *SendMsgListAsyncLogic) checkText(data *pb.MsgData) (*pb.SendMsgListResp, error) {
	// 是否开启了敏感词过滤
	if l.svcCtx.ConfigMgr.MessageShieldWordCheck(l.ctx, data.SenderId) {
		text := string(data.Content)
		sentence, found := ShieldWordTrieTreeInstance.Check(text)
		if found {
			// 是否不允许发
			if !l.svcCtx.ConfigMgr.MessageShieldWordAllow(l.ctx, data.SenderId) {
				// 直接报错返回
				return &pb.SendMsgListResp{CommonResp: pb.NewAlertErrorResp("发送失败", "内容包含违规词")}, nil
			}
			// 是否需要替换
			if l.svcCtx.ConfigMgr.MessageShieldWordAllowReplace(l.ctx, data.SenderId) {
				data.Content = []byte(sentence)
				// 检查offlinePush
				if data.OfflinePush != nil {
					data.OfflinePush.Title = l.svcCtx.ConfigMgr.OfflinePushTitle(l.ctx, data.SenderId)
					data.OfflinePush.Content = l.svcCtx.ConfigMgr.OfflinePushContent(l.ctx, data.SenderId)
				}
			}
		}
	}
	return &pb.SendMsgListResp{}, nil
}

func (l *SendMsgListAsyncLogic) SendMsgListAsync(in *pb.SendMsgListReq) (*pb.SendMsgListResp, error) {
	// check
	if len(in.MsgDataList) == 0 {
		return nil, nil
	}
	resp, err := l.check(in)
	if err != nil {
		return resp, err
	} else if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
		return resp, nil
	}
	if l.svcCtx.Config.TDMQ.Enabled && !l.svcCtx.SyncSendMsgLimiter.AllowCtx(l.ctx) {
		// 发送到消息队列
		var options []xtdmq.ProducerOptFunc
		if in.DeliverAfter != nil {
			options = append(options, xtdmq.ProduceWithDeliverAfter(time.Second*time.Duration(*in.DeliverAfter)))
		}
		_, err := l.svcCtx.MsgProducer().Produce(l.ctx, "msg", utils.ProtoToBytes(&pb.MsgMQBody{
			Event: pb.MsgMQBody_SendMsgListSync,
			Data:  utils.ProtoToBytes(in),
		}), options...)
		if err != nil {
			l.Errorf("MsgProducer.Produce error: %v", err)
			return &pb.SendMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		return &pb.SendMsgListResp{}, nil
	} else {
		var resp *pb.SendMsgListResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "SendMsgListSync", func(ctx context.Context) {
			resp, err = NewSendMsgListSyncLogic(ctx, l.svcCtx).SendMsgListSync(in)
		})
		return resp, err
	}
}
