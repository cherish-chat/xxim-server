package conngateway

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"strings"
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

var svcCtx *svc.ServiceContext

func Init(sc *svc.ServiceContext) {
	svcCtx = sc
}

type Route[REQ IReq, RESP IResp] struct {
	NewRequest func() REQ
	Do         func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
	Callback   func(ctx context.Context, resp RESP, c *types.UserConn)
}

var routeMap = map[string]func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error){}

func AddRoute[REQ IReq, RESP IResp](method string, route Route[REQ, RESP]) {
	routeMap[method] = func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error) {
		return OnReceiveCustom(ctx, method, c, body, route.NewRequest(), route.Do, route.Callback)
	}
}

type Handler func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error)

func Add(method string, handler Handler) {
	routeMap[method] = handler
}

func PrintRoutes() {
	strBuilder := strings.Builder{}
	for method := range routeMap {
		strBuilder.WriteString(method)
		strBuilder.WriteString("\n")
	}
	fmt.Printf("路由列表:\n%s", strBuilder.String())
}

func OnReceive(method string, ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error) {
	if c.ConnParam.UserId == "" || c.ConnParam.Token == "" {
		// 未登录
		logx.WithContext(ctx).Debugf("OnReceive: %s, user not login, conn的内存地址: %p", method, c)
		if !strings.Contains(method, "/white/") {
			// 不能访问
			return &pb.ResponseBody{
				Method: method,
				ReqId:  body.GetReqId(),
				Code:   pb.ResponseBody_AuthError,
				Data:   nil,
			}, nil
		}
	}
	if fn, ok := routeMap[method]; ok {
		return fn(ctx, c, body)
	}
	logx.Infof("OnReceive: %s, 404 not found", method)
	return nil, xerr.InvalidParamError
}
