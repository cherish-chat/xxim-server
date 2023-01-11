package server

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/common/pb"
	"strconv"
)

func (s *ConnServer) registerGateway() {
	// SendMsgListReq SendMsgListResp
	{
		route := conngateway.Route[*pb.SendMsgListReq, *pb.SendMsgListResp]{
			NewRequest: func() *pb.SendMsgListReq {
				return &pb.SendMsgListReq{}
			},
			Do: s.svcCtx.MsgService().SendMsgListAsync,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SendMsgList.Number())), route)
	}
	// BatchGetConvSeqReq BatchGetConvSeqResp
	{
		route := conngateway.Route[*pb.BatchGetConvSeqReq, *pb.BatchGetConvSeqResp]{
			NewRequest: func() *pb.BatchGetConvSeqReq {
				return &pb.BatchGetConvSeqReq{}
			},
			Do: s.svcCtx.MsgService().BatchGetConvSeq,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SyncConvSeq.Number())), route)
	}
	// BatchGetMsgListByConvIdReq GetMsgListResp
	{
		route := conngateway.Route[*pb.BatchGetMsgListByConvIdReq, *pb.GetMsgListResp]{
			NewRequest: func() *pb.BatchGetMsgListByConvIdReq {
				return &pb.BatchGetMsgListByConvIdReq{}
			},
			Do: s.svcCtx.MsgService().BatchGetMsgListByConvId,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_SyncMsgList.Number())), route)
	}
	// AckNoticeDataReq AckNoticeDataResp
	{
		route := conngateway.Route[*pb.AckNoticeDataReq, *pb.AckNoticeDataResp]{
			NewRequest: func() *pb.AckNoticeDataReq {
				return &pb.AckNoticeDataReq{}
			},
			Do: s.svcCtx.NoticeService().AckNoticeData,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_AckNotice.Number())), route)
	}
	// GetMsgByIdReq GetMsgByIdResp
	{
		route := conngateway.Route[*pb.GetMsgByIdReq, *pb.GetMsgByIdResp]{
			NewRequest: func() *pb.GetMsgByIdReq {
				return &pb.GetMsgByIdReq{}
			},
			Do: s.svcCtx.MsgService().GetMsgById,
		}
		conngateway.AddRoute(strconv.Itoa(int(pb.ActiveEvent_GetMsgById.Number())), route)
	}
}
