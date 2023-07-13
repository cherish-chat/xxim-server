// Code generated by goctl. DO NOT EDIT.
// Source: conversation.peer.proto

package server

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/logic/sessionservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"
)

type SessionServiceServer struct {
	svcCtx *svc.ServiceContext
	peerpb.UnimplementedSessionServiceServer
}

func NewSessionServiceServer(svcCtx *svc.ServiceContext) *SessionServiceServer {
	return &SessionServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *SessionServiceServer) ListJoinedConversations(ctx context.Context, in *peerpb.ListJoinedConversationsReq) (*peerpb.ListJoinedConversationsResp, error) {
	l := sessionservicelogic.NewListJoinedConversationsLogic(ctx, s.svcCtx)
	return l.ListJoinedConversations(in)
}
