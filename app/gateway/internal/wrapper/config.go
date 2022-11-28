package wrapper

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type IReq interface {
	proto.Message
	GetCommonReq() *pb.CommonReq
	SetCommonReq(*pb.CommonReq)
}

type IResp interface {
	proto.Message
	GetCommonResp() *pb.CommonResp
}

type Config[REQ IReq, RESP IResp] struct {
	NewRequest func() REQ
	Do         func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
}
