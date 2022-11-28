// Code generated by goctl. DO NOT EDIT!
// Source: msg.proto

package msgservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BatchGetConvSeqReq            = pb.BatchGetConvSeqReq
	BatchGetConvSeqResp           = pb.BatchGetConvSeqResp
	BatchGetConvSeqResp_ConvSeq   = pb.BatchGetConvSeqResp_ConvSeq
	BatchSetMinSeqReq             = pb.BatchSetMinSeqReq
	BatchSetMinSeqResp            = pb.BatchSetMinSeqResp
	GetConvSubscribersReq         = pb.GetConvSubscribersReq
	GetConvSubscribersResp        = pb.GetConvSubscribersResp
	GetConvSubscribersRespUidList = pb.GetConvSubscribersRespUidList
	GetMsgByIdReq                 = pb.GetMsgByIdReq
	GetMsgByIdResp                = pb.GetMsgByIdResp
	GetMsgListByConvIdReq         = pb.GetMsgListByConvIdReq
	GetMsgListResp                = pb.GetMsgListResp
	MsgData                       = pb.MsgData
	MsgDataList                   = pb.MsgDataList
	MsgData_OfflinePush           = pb.MsgData_OfflinePush
	MsgData_Options               = pb.MsgData_Options
	MsgData_Receiver              = pb.MsgData_Receiver
	MsgMQBody                     = pb.MsgMQBody
	PushMsgListReq                = pb.PushMsgListReq
	SendMsgListReq                = pb.SendMsgListReq
	SendMsgListResp               = pb.SendMsgListResp

	MsgService interface {
		InsertMsgDataList(ctx context.Context, in *MsgDataList, opts ...grpc.CallOption) (*MsgDataList, error)
		SendMsgListSync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
		SendMsgListAsync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
		PushMsgList(ctx context.Context, in *PushMsgListReq, opts ...grpc.CallOption) (*CommonResp, error)
		// GetMsgListByConvId 通过seq拉取一个会话的消息
		GetMsgListByConvId(ctx context.Context, in *GetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error)
		// GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
		GetMsgById(ctx context.Context, in *GetMsgByIdReq, opts ...grpc.CallOption) (*GetMsgByIdResp, error)
		// BatchSetMinSeq 批量设置用户某会话的minseq
		BatchSetMinSeq(ctx context.Context, in *BatchSetMinSeqReq, opts ...grpc.CallOption) (*BatchSetMinSeqResp, error)
		// BatchGetConvSeq 批量获取会话的seq
		BatchGetConvSeq(ctx context.Context, in *BatchGetConvSeqReq, opts ...grpc.CallOption) (*BatchGetConvSeqResp, error)
		//  conn hook
		AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		// GetConvSubscribers 获取一个会话里所有的消息订阅者
		GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error)
	}

	defaultMsgService struct {
		cli zrpc.Client
	}
)

func NewMsgService(cli zrpc.Client) MsgService {
	return &defaultMsgService{
		cli: cli,
	}
}

func (m *defaultMsgService) InsertMsgDataList(ctx context.Context, in *MsgDataList, opts ...grpc.CallOption) (*MsgDataList, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.InsertMsgDataList(ctx, in, opts...)
}

func (m *defaultMsgService) SendMsgListSync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.SendMsgListSync(ctx, in, opts...)
}

func (m *defaultMsgService) SendMsgListAsync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.SendMsgListAsync(ctx, in, opts...)
}

func (m *defaultMsgService) PushMsgList(ctx context.Context, in *PushMsgListReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.PushMsgList(ctx, in, opts...)
}

// GetMsgListByConvId 通过seq拉取一个会话的消息
func (m *defaultMsgService) GetMsgListByConvId(ctx context.Context, in *GetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetMsgListByConvId(ctx, in, opts...)
}

// GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
func (m *defaultMsgService) GetMsgById(ctx context.Context, in *GetMsgByIdReq, opts ...grpc.CallOption) (*GetMsgByIdResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetMsgById(ctx, in, opts...)
}

// BatchSetMinSeq 批量设置用户某会话的minseq
func (m *defaultMsgService) BatchSetMinSeq(ctx context.Context, in *BatchSetMinSeqReq, opts ...grpc.CallOption) (*BatchSetMinSeqResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.BatchSetMinSeq(ctx, in, opts...)
}

// BatchGetConvSeq 批量获取会话的seq
func (m *defaultMsgService) BatchGetConvSeq(ctx context.Context, in *BatchGetConvSeqReq, opts ...grpc.CallOption) (*BatchGetConvSeqResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.BatchGetConvSeq(ctx, in, opts...)
}

//  conn hook
func (m *defaultMsgService) AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.AfterConnect(ctx, in, opts...)
}

func (m *defaultMsgService) AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.AfterDisconnect(ctx, in, opts...)
}

// GetConvSubscribers 获取一个会话里所有的消息订阅者
func (m *defaultMsgService) GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetConvSubscribers(ctx, in, opts...)
}
