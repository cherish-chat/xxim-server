package conngateway

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils/xerr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type IReq interface {
	proto.Message
	GetCommonReq() *pb.CommonReq
	SetCommonReq(*pb.CommonReq)
	Validate() error
	Path() string
}

type IResp interface {
	proto.Message
	GetCommonResp() *pb.CommonResp
}

type Route[REQ IReq, RESP IResp] struct {
	NewRequest func() REQ
	Do         func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
}

var routeMap = map[string]func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error){}

func AddRoute[REQ IReq, RESP IResp](method string, route Route[REQ, RESP]) {
	routeMap[method] = func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error) {
		return OnReceiveCustom(ctx, c, body, route.NewRequest(), route.Do)
	}
}

func OnReceive(method string, ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error) {
	if fn, ok := routeMap[method]; ok {
		return fn(ctx, c, body)
	}
	return nil, xerr.InvalidParamError
}
