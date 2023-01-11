package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/common/pb"
	"strconv"
)

func (s *ConnServer) registerGateway() {
	{
		route := conngateway.Route[*pb.SendMsgListReq, *pb.SendMsgListResp]{
			NewRequest: func() *pb.SendMsgListReq {
				return &pb.SendMsgListReq{}
			},
			Do: s.svcCtx.MsgService().SendMsgListAsync,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SendMsgList.Number())), route)
	}
}
