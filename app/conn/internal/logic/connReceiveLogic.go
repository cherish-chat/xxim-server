package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

func (l *ConnLogic) OnReceive(c *types.UserConn, typ int, msg []byte) {
	switch websocket.MessageType(typ) {
	case websocket.MessageBinary:
		// 接收到消息
		body := &pb.RequestBody{}
		var bodyData []byte
		err := proto.Unmarshal(msg, body)
		if err == nil {
			var respBody *pb.ResponseBody
			respBody, err = l.onReceiveBody(c, body)
			bodyData, _ = proto.Marshal(respBody)
		}
		if err != nil {
			logx.WithContext(c.Ctx).Errorf("OnReceiveBody error: %s", err.Error())
		}
		data, _ := proto.Marshal(&pb.PushBody{
			Event: pb.PushEvent_PushResponseBody,
			Data:  bodyData,
		})
		err = l.SendMsgToConn(c, data)
		if err != nil {
			logx.WithContext(c.Ctx).Errorf("SendMsgToConn error: %s", err.Error())
		}
	default:
		// 无效的消息类型
		l.Errorf("invalid message type: %d, msg: %s", typ, string(msg))
	}
}

func (l *ConnLogic) onReceiveBody(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	switch body.Event {
	case pb.ActiveEvent_SendMsgList:
		return l.onReceiveSendMsgList(c, body)
	case pb.ActiveEvent_SyncConvSeq:
		return l.onReceiveSyncConvSeq(c, body)
	case pb.ActiveEvent_SyncMsgList:
		return l.onReceiveSyncMsgList(c, body)
	case pb.ActiveEvent_AckNotice:
		return l.onReceiveAckNotice(c, body)
	case pb.ActiveEvent_GetMsgById:
		return l.onReceiveGetMsgById(c, body)
	default:
		return nil, errors.New("invalid event")
	}
}

func (l *ConnLogic) onReceiveSendMsgList(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	req := &pb.SendMsgListReq{}
	err := proto.Unmarshal(body.Data, req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("SendMsgListReq unmarshal error: %s", err.Error())
		return nil, err
	}
	var resp *pb.SendMsgListResp
	xtrace.RunWithTrace("", "onReceiveSendMsgList", func(ctx context.Context) {
		resp, err = l.svcCtx.MsgService().SendMsgListAsync(ctx, req)
	}, propagation.MapCarrier{
		"req-id":    body.ReqId,
		"event":     body.Event.String(),
		"user-id":   c.ConnParam.UserId,
		"platform":  c.ConnParam.Platform,
		"device-id": c.ConnParam.DeviceId,
		"ip":        c.ConnParam.Ips,
		"net-type":  c.ConnParam.NetworkUsed,
	})
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("SendMsgList error: %s", err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	return &pb.ResponseBody{
		Event: body.Event,
		ReqId: body.ReqId,
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}

func (l *ConnLogic) onReceiveSyncConvSeq(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	req := &pb.BatchGetConvSeqReq{}
	err := proto.Unmarshal(body.Data, req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("BatchGetConvSeqReq unmarshal error: %s", err.Error())
		return nil, err
	}
	var resp *pb.BatchGetConvSeqResp
	xtrace.RunWithTrace("", "onReceiveSyncConvSeq", func(ctx context.Context) {
		resp, err = l.svcCtx.MsgService().BatchGetConvSeq(ctx, req)
	}, propagation.MapCarrier{
		"req-id":    body.ReqId,
		"event":     body.Event.String(),
		"user-id":   c.ConnParam.UserId,
		"platform":  c.ConnParam.Platform,
		"device-id": c.ConnParam.DeviceId,
		"ip":        c.ConnParam.Ips,
		"net-type":  c.ConnParam.NetworkUsed,
	})
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("BatchGetConvSeq error: %s", err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	return &pb.ResponseBody{
		Event: body.Event,
		ReqId: body.ReqId,
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}

func (l *ConnLogic) onReceiveSyncMsgList(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	req := &pb.BatchGetMsgListByConvIdReq{}
	err := proto.Unmarshal(body.Data, req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("SyncMsgListReq unmarshal error: %s", err.Error())
		return nil, err
	}
	var resp *pb.GetMsgListResp
	xtrace.RunWithTrace("", "onReceiveSyncMsgList", func(ctx context.Context) {
		resp, err = l.svcCtx.MsgService().BatchGetMsgListByConvId(ctx, req)
	}, propagation.MapCarrier{
		"req-id":    body.ReqId,
		"event":     body.Event.String(),
		"user-id":   c.ConnParam.UserId,
		"platform":  c.ConnParam.Platform,
		"device-id": c.ConnParam.DeviceId,
		"ip":        c.ConnParam.Ips,
		"net-type":  c.ConnParam.NetworkUsed,
	})
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("BatchGetMsgListByConvId error: %s", err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	return &pb.ResponseBody{
		Event: body.Event,
		ReqId: body.ReqId,
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}

func (l *ConnLogic) onReceiveAckNotice(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	req := &pb.AckNoticeDataReq{}
	err := proto.Unmarshal(body.Data, req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("AckNoticeReq unmarshal error: %s", err.Error())
		return nil, err
	}
	var resp *pb.AckNoticeDataResp
	xtrace.RunWithTrace("", "onReceiveAckNotice", func(ctx context.Context) {
		resp, err = l.svcCtx.NoticeService().AckNoticeData(ctx, req)
	}, propagation.MapCarrier{
		"req-id":    body.ReqId,
		"event":     body.Event.String(),
		"user-id":   c.ConnParam.UserId,
		"platform":  c.ConnParam.Platform,
		"device-id": c.ConnParam.DeviceId,
		"ip":        c.ConnParam.Ips,
		"net-type":  c.ConnParam.NetworkUsed,
	})
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("AckNoticeData error: %s", err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	return &pb.ResponseBody{
		Event: body.Event,
		ReqId: body.ReqId,
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}

func (l *ConnLogic) onReceiveGetMsgById(c *types.UserConn, body *pb.RequestBody) (*pb.ResponseBody, error) {
	req := &pb.GetMsgByIdReq{}
	err := proto.Unmarshal(body.Data, req)
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("GetMsgByIdReq unmarshal error: %s", err.Error())
		return nil, err
	}
	var resp *pb.GetMsgByIdResp
	xtrace.RunWithTrace("", "onReceiveGetMsgById", func(ctx context.Context) {
		resp, err = l.svcCtx.MsgService().GetMsgById(ctx, req)
	}, propagation.MapCarrier{
		"req-id":    body.ReqId,
		"event":     body.Event.String(),
		"user-id":   c.ConnParam.UserId,
		"platform":  c.ConnParam.Platform,
		"device-id": c.ConnParam.DeviceId,
		"ip":        c.ConnParam.Ips,
		"net-type":  c.ConnParam.NetworkUsed,
	})
	if err != nil {
		logx.WithContext(c.Ctx).Errorf("GetMsgById error: %s", err.Error())
	}
	respBuff, _ := proto.Marshal(resp)
	return &pb.ResponseBody{
		Event: body.Event,
		ReqId: body.ReqId,
		Code:  pb.ResponseBody_Code(resp.GetCommonResp().GetCode()),
		Data:  respBuff,
	}, err
}
