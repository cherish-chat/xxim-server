// Code generated by goctl. DO NOT EDIT.
// Source: msg.proto

package msgservice

import (
	"context"

	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BatchGetConvSeqReq              = pb.BatchGetConvSeqReq
	BatchGetConvSeqResp             = pb.BatchGetConvSeqResp
	BatchGetConvSeqResp_ConvSeq     = pb.BatchGetConvSeqResp_ConvSeq
	BatchGetMsgListByConvIdReq      = pb.BatchGetMsgListByConvIdReq
	BatchGetMsgListByConvIdReq_Item = pb.BatchGetMsgListByConvIdReq_Item
	BatchSetMinSeqReq               = pb.BatchSetMinSeqReq
	BatchSetMinSeqResp              = pb.BatchSetMinSeqResp
	EditMsgReq                      = pb.EditMsgReq
	EditMsgResp                     = pb.EditMsgResp
	FlushShieldWordTireTreeReq      = pb.FlushShieldWordTireTreeReq
	FlushShieldWordTireTreeResp     = pb.FlushShieldWordTireTreeResp
	FlushUsersSubConvReq            = pb.FlushUsersSubConvReq
	GetAllMsgListReq                = pb.GetAllMsgListReq
	GetAllMsgListResp               = pb.GetAllMsgListResp
	GetConvOnlineCountReq           = pb.GetConvOnlineCountReq
	GetConvOnlineCountResp          = pb.GetConvOnlineCountResp
	GetConvSubscribersReq           = pb.GetConvSubscribersReq
	GetConvSubscribersResp          = pb.GetConvSubscribersResp
	GetMsgByIdReq                   = pb.GetMsgByIdReq
	GetMsgByIdResp                  = pb.GetMsgByIdResp
	GetMsgListResp                  = pb.GetMsgListResp
	GetRedPacketDetailReq           = pb.GetRedPacketDetailReq
	GetRedPacketDetailResp          = pb.GetRedPacketDetailResp
	MsgData                         = pb.MsgData
	MsgDataList                     = pb.MsgDataList
	MsgData_OfflinePush             = pb.MsgData_OfflinePush
	MsgData_Options                 = pb.MsgData_Options
	MsgMQBody                       = pb.MsgMQBody
	OfflinePushMsgReq               = pb.OfflinePushMsgReq
	OfflinePushMsgResp              = pb.OfflinePushMsgResp
	PushMsgListReq                  = pb.PushMsgListReq
	ReadMsgReq                      = pb.ReadMsgReq
	ReadMsgResp                     = pb.ReadMsgResp
	ReceiveRedPacketReq             = pb.ReceiveRedPacketReq
	ReceiveRedPacketResp            = pb.ReceiveRedPacketResp
	RedPacket                       = pb.RedPacket
	RedPacket_Receiver              = pb.RedPacket_Receiver
	SendMsgListReq                  = pb.SendMsgListReq
	SendMsgListResp                 = pb.SendMsgListResp
	SendRedPacketReq                = pb.SendRedPacketReq
	SendRedPacketResp               = pb.SendRedPacketResp

	MsgService interface {
		InsertMsgDataList(ctx context.Context, in *MsgDataList, opts ...grpc.CallOption) (*MsgDataList, error)
		SendMsgListSync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
		SendMsgListAsync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
		PushMsgList(ctx context.Context, in *PushMsgListReq, opts ...grpc.CallOption) (*CommonResp, error)
		// BatchGetMsgListByConvId 通过seq拉取一个会话的消息
		BatchGetMsgListByConvId(ctx context.Context, in *BatchGetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error)
		// GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
		GetMsgById(ctx context.Context, in *GetMsgByIdReq, opts ...grpc.CallOption) (*GetMsgByIdResp, error)
		// BatchSetMinSeq 批量设置用户某会话的minseq
		BatchSetMinSeq(ctx context.Context, in *BatchSetMinSeqReq, opts ...grpc.CallOption) (*BatchSetMinSeqResp, error)
		// BatchGetConvSeq 批量获取会话的seq
		BatchGetConvSeq(ctx context.Context, in *BatchGetConvSeqReq, opts ...grpc.CallOption) (*BatchGetConvSeqResp, error)
		// conn hook
		AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error)
		KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*KeepAliveResp, error)
		// GetConvSubscribers 获取一个会话里所有的消息订阅者
		GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error)
		// OfflinePushMsg 离线推送消息
		OfflinePushMsg(ctx context.Context, in *OfflinePushMsgReq, opts ...grpc.CallOption) (*OfflinePushMsgResp, error)
		// GetConvOnlineCount 获取一个会话里所有的在线用户
		GetConvOnlineCount(ctx context.Context, in *GetConvOnlineCountReq, opts ...grpc.CallOption) (*GetConvOnlineCountResp, error)
		// FlushUsersSubConv 刷新用户订阅的会话
		FlushUsersSubConv(ctx context.Context, in *FlushUsersSubConvReq, opts ...grpc.CallOption) (*CommonResp, error)
		// GetAllMsgList 获取所有消息
		GetAllMsgList(ctx context.Context, in *GetAllMsgListReq, opts ...grpc.CallOption) (*GetAllMsgListResp, error)
		// ReadMsg 设置会话已读
		ReadMsg(ctx context.Context, in *ReadMsgReq, opts ...grpc.CallOption) (*ReadMsgResp, error)
		// EditMsg 编辑消息
		EditMsg(ctx context.Context, in *EditMsgReq, opts ...grpc.CallOption) (*EditMsgResp, error)
		// FlushShieldWordTireTree 刷新屏蔽词
		FlushShieldWordTireTree(ctx context.Context, in *FlushShieldWordTireTreeReq, opts ...grpc.CallOption) (*FlushShieldWordTireTreeResp, error)
		// SendRedPacket 发红包
		SendRedPacket(ctx context.Context, in *SendRedPacketReq, opts ...grpc.CallOption) (*SendRedPacketResp, error)
		// ReceiveRedPacket 领取红包
		ReceiveRedPacket(ctx context.Context, in *ReceiveRedPacketReq, opts ...grpc.CallOption) (*ReceiveRedPacketResp, error)
		// GetRedPacketDetail 获取红包详情
		GetRedPacketDetail(ctx context.Context, in *GetRedPacketDetailReq, opts ...grpc.CallOption) (*GetRedPacketDetailResp, error)
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

// BatchGetMsgListByConvId 通过seq拉取一个会话的消息
func (m *defaultMsgService) BatchGetMsgListByConvId(ctx context.Context, in *BatchGetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.BatchGetMsgListByConvId(ctx, in, opts...)
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

// conn hook
func (m *defaultMsgService) AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.AfterConnect(ctx, in, opts...)
}

func (m *defaultMsgService) AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.AfterDisconnect(ctx, in, opts...)
}

func (m *defaultMsgService) KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*KeepAliveResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.KeepAlive(ctx, in, opts...)
}

// GetConvSubscribers 获取一个会话里所有的消息订阅者
func (m *defaultMsgService) GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetConvSubscribers(ctx, in, opts...)
}

// OfflinePushMsg 离线推送消息
func (m *defaultMsgService) OfflinePushMsg(ctx context.Context, in *OfflinePushMsgReq, opts ...grpc.CallOption) (*OfflinePushMsgResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.OfflinePushMsg(ctx, in, opts...)
}

// GetConvOnlineCount 获取一个会话里所有的在线用户
func (m *defaultMsgService) GetConvOnlineCount(ctx context.Context, in *GetConvOnlineCountReq, opts ...grpc.CallOption) (*GetConvOnlineCountResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetConvOnlineCount(ctx, in, opts...)
}

// FlushUsersSubConv 刷新用户订阅的会话
func (m *defaultMsgService) FlushUsersSubConv(ctx context.Context, in *FlushUsersSubConvReq, opts ...grpc.CallOption) (*CommonResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.FlushUsersSubConv(ctx, in, opts...)
}

// GetAllMsgList 获取所有消息
func (m *defaultMsgService) GetAllMsgList(ctx context.Context, in *GetAllMsgListReq, opts ...grpc.CallOption) (*GetAllMsgListResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetAllMsgList(ctx, in, opts...)
}

// ReadMsg 设置会话已读
func (m *defaultMsgService) ReadMsg(ctx context.Context, in *ReadMsgReq, opts ...grpc.CallOption) (*ReadMsgResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.ReadMsg(ctx, in, opts...)
}

// EditMsg 编辑消息
func (m *defaultMsgService) EditMsg(ctx context.Context, in *EditMsgReq, opts ...grpc.CallOption) (*EditMsgResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.EditMsg(ctx, in, opts...)
}

// FlushShieldWordTireTree 刷新屏蔽词
func (m *defaultMsgService) FlushShieldWordTireTree(ctx context.Context, in *FlushShieldWordTireTreeReq, opts ...grpc.CallOption) (*FlushShieldWordTireTreeResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.FlushShieldWordTireTree(ctx, in, opts...)
}

// SendRedPacket 发红包
func (m *defaultMsgService) SendRedPacket(ctx context.Context, in *SendRedPacketReq, opts ...grpc.CallOption) (*SendRedPacketResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.SendRedPacket(ctx, in, opts...)
}

// ReceiveRedPacket 领取红包
func (m *defaultMsgService) ReceiveRedPacket(ctx context.Context, in *ReceiveRedPacketReq, opts ...grpc.CallOption) (*ReceiveRedPacketResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.ReceiveRedPacket(ctx, in, opts...)
}

// GetRedPacketDetail 获取红包详情
func (m *defaultMsgService) GetRedPacketDetail(ctx context.Context, in *GetRedPacketDetailReq, opts ...grpc.CallOption) (*GetRedPacketDetailResp, error) {
	client := pb.NewMsgServiceClient(m.cli.Conn())
	return client.GetRedPacketDetail(ctx, in, opts...)
}
