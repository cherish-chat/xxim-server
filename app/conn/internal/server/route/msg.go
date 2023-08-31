package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterMsg(svcCtx *svc.ServiceContext) {
	// SendMsgListReq SendMsgListResp
	{
		route := conngateway.Route[*pb.SendMsgListReq, *pb.SendMsgListResp]{
			NewRequest: func() *pb.SendMsgListReq {
				return &pb.SendMsgListReq{}
			},
			Do: svcCtx.MsgService().SendMsgListAsync,
		}
		conngateway.AddRoute("/v1/msg/sendMsgList", route)
	}
	// BatchGetConvSeqReq BatchGetConvSeqResp
	{
		route := conngateway.Route[*pb.BatchGetConvSeqReq, *pb.BatchGetConvSeqResp]{
			NewRequest: func() *pb.BatchGetConvSeqReq {
				return &pb.BatchGetConvSeqReq{}
			},
			Do: svcCtx.MsgService().BatchGetConvSeq,
		}
		conngateway.AddRoute("/v1/msg/batchGetConvSeq", route)
	}
	// BatchGetMsgListByConvIdReq GetMsgListResp
	{
		route := conngateway.Route[*pb.BatchGetMsgListByConvIdReq, *pb.GetMsgListResp]{
			NewRequest: func() *pb.BatchGetMsgListByConvIdReq {
				return &pb.BatchGetMsgListByConvIdReq{}
			},
			Do: svcCtx.MsgService().BatchGetMsgListByConvId,
		}
		conngateway.AddRoute("/v1/msg/batchGetMsgListByConvId", route)
	}
	// GetMsgByIdReq GetMsgByIdResp
	{
		route := conngateway.Route[*pb.GetMsgByIdReq, *pb.GetMsgByIdResp]{
			NewRequest: func() *pb.GetMsgByIdReq {
				return &pb.GetMsgByIdReq{}
			},
			Do: svcCtx.MsgService().GetMsgById,
		}
		conngateway.AddRoute("/v1/msg/getMsgById", route)
	}
	// ReadMsgReq ReadMsgResp
	{
		route := conngateway.Route[*pb.ReadMsgReq, *pb.ReadMsgResp]{
			NewRequest: func() *pb.ReadMsgReq {
				return &pb.ReadMsgReq{}
			},
			Do: svcCtx.MsgService().ReadMsg,
		}
		conngateway.AddRoute("/v1/msg/sendReadMsg", route)
	}
	// EditMsgReq EditMsgResp
	{
		route := conngateway.Route[*pb.EditMsgReq, *pb.EditMsgResp]{
			NewRequest: func() *pb.EditMsgReq {
				return &pb.EditMsgReq{}
			},
			Do: svcCtx.MsgService().EditMsg,
		}
		conngateway.AddRoute("/v1/msg/sendEditMsg", route)
	}
	//SendRedPacketReq SendRedPacketResp
	{
		route := conngateway.Route[*pb.SendRedPacketReq, *pb.SendRedPacketResp]{
			NewRequest: func() *pb.SendRedPacketReq {
				return &pb.SendRedPacketReq{}
			},
			Do: svcCtx.MsgService().SendRedPacket,
		}
		conngateway.AddRoute("/v1/msg/sendRedPacket", route)
	}
	//ReceiveRedPacketReq ReceiveRedPacketResp
	{
		route := conngateway.Route[*pb.ReceiveRedPacketReq, *pb.ReceiveRedPacketResp]{
			NewRequest: func() *pb.ReceiveRedPacketReq {
				return &pb.ReceiveRedPacketReq{}
			},
			Do: svcCtx.MsgService().ReceiveRedPacket,
		}
		conngateway.AddRoute("/v1/msg/receiveRedPacket", route)
	}
	//GetRedPacketDetailReq GetRedPacketDetailResp
	{
		route := conngateway.Route[*pb.GetRedPacketDetailReq, *pb.GetRedPacketDetailResp]{
			NewRequest: func() *pb.GetRedPacketDetailReq {
				return &pb.GetRedPacketDetailReq{}
			},
			Do: svcCtx.MsgService().GetRedPacketDetail,
		}
		conngateway.AddRoute("/v1/msg/getRedPacketDetail", route)
	}
}
