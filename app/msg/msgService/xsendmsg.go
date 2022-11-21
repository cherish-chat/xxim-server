package msgservice

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type SendMsgOption struct {
	DeliverAfter *int32
}
type SendMsgOptionFunc func(*SendMsgOption)

func SendMsgWithDeliverAfter(deliverAfter int32) SendMsgOptionFunc {
	return func(opt *SendMsgOption) {
		opt.DeliverAfter = &deliverAfter
	}
}

func SendMsgSync(
	service MsgService,
	ctx context.Context,
	msgDataList []*pb.MsgData,
	opts ...SendMsgOptionFunc,
) (*SendMsgListResp, error) {
	opt := &SendMsgOption{}
	for _, o := range opts {
		o(opt)
	}
	return service.SendMsgListSync(ctx, &SendMsgListReq{
		CommonReq: &pb.CommonReq{
			Platform: "system",
		},
		MsgDataList:  msgDataList,
		DeliverAfter: opt.DeliverAfter,
	})
}
