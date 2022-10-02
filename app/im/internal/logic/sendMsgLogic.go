package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendMsg 发送消息到 pulsar
func (l *SendMsgLogic) SendMsg(in *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	err := l.beforeSendMsg(in)
	if err != nil {
		return &pb.SendMsgResp{FailedReason: err.Error()}, nil
	}
	msg := &pb.MsgToMQData{MsgDataList: in.MsgDataList}
	bytes, _ := proto.Marshal(msg)
	var options = []xmq.ProducerOptFunc{xmq.ProduceWithProperties(map[string]string{
		"selfId":      in.SelfId,
		"platform":    in.Platform,
		"appVersion":  in.AppVersion,
		"deviceModel": in.DeviceModel,
		"ips":         in.Ips,
	})}
	if in.SendAt != nil {
		if *in.SendAt > utils.GetNowMilli()+60000 {
			sendAt := time.UnixMilli(*in.SendAt)
			options = append(options, xmq.ProduceWithDeliverAt(sendAt))
		}
	}
	_, err = l.svcCtx.StorageProducer().Produce(l.ctx, "default", bytes, options...)
	if err != nil {
		l.Errorf("send msg to pulsar failed, err: %v", err)
		return nil, err
	}
	l.afterSendMsg(in)
	return &pb.SendMsgResp{}, nil
}

func (l *SendMsgLogic) beforeSendMsg(in *pb.SendMsgReq) error {
	// todo: 验证消息、判断消息是否允许发送
	return nil
}

func (l *SendMsgLogic) afterSendMsg(in *pb.SendMsgReq) {
	// todo: 消息发送成功后的处理
}
