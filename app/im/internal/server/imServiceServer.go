// Code generated by goctl. DO NOT EDIT!
// Source: im.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/im/internal/logic"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type ImServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedImServiceServer
}

func NewImServiceServer(svcCtx *svc.ServiceContext) *ImServiceServer {
	return &ImServiceServer{
		svcCtx: svcCtx,
	}
}

// SendMsg 发送消息到 pulsar
func (s *ImServiceServer) SendMsg(ctx context.Context, in *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	l := logic.NewSendMsgLogic(ctx, s.svcCtx)
	return l.SendMsg(in)
}
