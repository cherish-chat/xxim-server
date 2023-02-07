// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: msg.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MsgServiceClient is the client API for MsgService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgServiceClient interface {
	InsertMsgDataList(ctx context.Context, in *MsgDataList, opts ...grpc.CallOption) (*MsgDataList, error)
	SendMsgListSync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
	SendMsgListAsync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error)
	PushMsgList(ctx context.Context, in *PushMsgListReq, opts ...grpc.CallOption) (*CommonResp, error)
	//BatchGetMsgListByConvId 通过seq拉取一个会话的消息
	BatchGetMsgListByConvId(ctx context.Context, in *BatchGetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error)
	//GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
	GetMsgById(ctx context.Context, in *GetMsgByIdReq, opts ...grpc.CallOption) (*GetMsgByIdResp, error)
	//BatchSetMinSeq 批量设置用户某会话的minseq
	BatchSetMinSeq(ctx context.Context, in *BatchSetMinSeqReq, opts ...grpc.CallOption) (*BatchSetMinSeqResp, error)
	//BatchGetConvSeq 批量获取会话的seq
	BatchGetConvSeq(ctx context.Context, in *BatchGetConvSeqReq, opts ...grpc.CallOption) (*BatchGetConvSeqResp, error)
	// conn hook
	AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error)
	AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error)
	KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*KeepAliveResp, error)
	//GetConvSubscribers 获取一个会话里所有的消息订阅者
	GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error)
	//OfflinePushMsg 离线推送消息
	OfflinePushMsg(ctx context.Context, in *OfflinePushMsgReq, opts ...grpc.CallOption) (*OfflinePushMsgResp, error)
	//GetConvOnlineCount 获取一个会话里所有的在线用户
	GetConvOnlineCount(ctx context.Context, in *GetConvOnlineCountReq, opts ...grpc.CallOption) (*GetConvOnlineCountResp, error)
	//FlushUsersSubConv 刷新用户订阅的会话
	FlushUsersSubConv(ctx context.Context, in *FlushUsersSubConvReq, opts ...grpc.CallOption) (*CommonResp, error)
	//GetAllMsgList 获取所有消息
	GetAllMsgList(ctx context.Context, in *GetAllMsgListReq, opts ...grpc.CallOption) (*GetAllMsgListResp, error)
	//ReadMsg 设置会话已读
	ReadMsg(ctx context.Context, in *ReadMsgReq, opts ...grpc.CallOption) (*ReadMsgResp, error)
	//EditMsg 编辑消息
	EditMsg(ctx context.Context, in *EditMsgReq, opts ...grpc.CallOption) (*EditMsgResp, error)
}

type msgServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgServiceClient(cc grpc.ClientConnInterface) MsgServiceClient {
	return &msgServiceClient{cc}
}

func (c *msgServiceClient) InsertMsgDataList(ctx context.Context, in *MsgDataList, opts ...grpc.CallOption) (*MsgDataList, error) {
	out := new(MsgDataList)
	err := c.cc.Invoke(ctx, "/pb.msgService/InsertMsgDataList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) SendMsgListSync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error) {
	out := new(SendMsgListResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/SendMsgListSync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) SendMsgListAsync(ctx context.Context, in *SendMsgListReq, opts ...grpc.CallOption) (*SendMsgListResp, error) {
	out := new(SendMsgListResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/SendMsgListAsync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) PushMsgList(ctx context.Context, in *PushMsgListReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/PushMsgList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) BatchGetMsgListByConvId(ctx context.Context, in *BatchGetMsgListByConvIdReq, opts ...grpc.CallOption) (*GetMsgListResp, error) {
	out := new(GetMsgListResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/BatchGetMsgListByConvId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) GetMsgById(ctx context.Context, in *GetMsgByIdReq, opts ...grpc.CallOption) (*GetMsgByIdResp, error) {
	out := new(GetMsgByIdResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/GetMsgById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) BatchSetMinSeq(ctx context.Context, in *BatchSetMinSeqReq, opts ...grpc.CallOption) (*BatchSetMinSeqResp, error) {
	out := new(BatchSetMinSeqResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/BatchSetMinSeq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) BatchGetConvSeq(ctx context.Context, in *BatchGetConvSeqReq, opts ...grpc.CallOption) (*BatchGetConvSeqResp, error) {
	out := new(BatchGetConvSeqResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/BatchGetConvSeq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/AfterConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/AfterDisconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*KeepAliveResp, error) {
	out := new(KeepAliveResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) GetConvSubscribers(ctx context.Context, in *GetConvSubscribersReq, opts ...grpc.CallOption) (*GetConvSubscribersResp, error) {
	out := new(GetConvSubscribersResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/GetConvSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) OfflinePushMsg(ctx context.Context, in *OfflinePushMsgReq, opts ...grpc.CallOption) (*OfflinePushMsgResp, error) {
	out := new(OfflinePushMsgResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/OfflinePushMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) GetConvOnlineCount(ctx context.Context, in *GetConvOnlineCountReq, opts ...grpc.CallOption) (*GetConvOnlineCountResp, error) {
	out := new(GetConvOnlineCountResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/GetConvOnlineCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) FlushUsersSubConv(ctx context.Context, in *FlushUsersSubConvReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/FlushUsersSubConv", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) GetAllMsgList(ctx context.Context, in *GetAllMsgListReq, opts ...grpc.CallOption) (*GetAllMsgListResp, error) {
	out := new(GetAllMsgListResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/GetAllMsgList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) ReadMsg(ctx context.Context, in *ReadMsgReq, opts ...grpc.CallOption) (*ReadMsgResp, error) {
	out := new(ReadMsgResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/ReadMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgServiceClient) EditMsg(ctx context.Context, in *EditMsgReq, opts ...grpc.CallOption) (*EditMsgResp, error) {
	out := new(EditMsgResp)
	err := c.cc.Invoke(ctx, "/pb.msgService/EditMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServiceServer is the server API for MsgService service.
// All implementations must embed UnimplementedMsgServiceServer
// for forward compatibility
type MsgServiceServer interface {
	InsertMsgDataList(context.Context, *MsgDataList) (*MsgDataList, error)
	SendMsgListSync(context.Context, *SendMsgListReq) (*SendMsgListResp, error)
	SendMsgListAsync(context.Context, *SendMsgListReq) (*SendMsgListResp, error)
	PushMsgList(context.Context, *PushMsgListReq) (*CommonResp, error)
	//BatchGetMsgListByConvId 通过seq拉取一个会话的消息
	BatchGetMsgListByConvId(context.Context, *BatchGetMsgListByConvIdReq) (*GetMsgListResp, error)
	//GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
	GetMsgById(context.Context, *GetMsgByIdReq) (*GetMsgByIdResp, error)
	//BatchSetMinSeq 批量设置用户某会话的minseq
	BatchSetMinSeq(context.Context, *BatchSetMinSeqReq) (*BatchSetMinSeqResp, error)
	//BatchGetConvSeq 批量获取会话的seq
	BatchGetConvSeq(context.Context, *BatchGetConvSeqReq) (*BatchGetConvSeqResp, error)
	// conn hook
	AfterConnect(context.Context, *AfterConnectReq) (*CommonResp, error)
	AfterDisconnect(context.Context, *AfterDisconnectReq) (*CommonResp, error)
	KeepAlive(context.Context, *KeepAliveReq) (*KeepAliveResp, error)
	//GetConvSubscribers 获取一个会话里所有的消息订阅者
	GetConvSubscribers(context.Context, *GetConvSubscribersReq) (*GetConvSubscribersResp, error)
	//OfflinePushMsg 离线推送消息
	OfflinePushMsg(context.Context, *OfflinePushMsgReq) (*OfflinePushMsgResp, error)
	//GetConvOnlineCount 获取一个会话里所有的在线用户
	GetConvOnlineCount(context.Context, *GetConvOnlineCountReq) (*GetConvOnlineCountResp, error)
	//FlushUsersSubConv 刷新用户订阅的会话
	FlushUsersSubConv(context.Context, *FlushUsersSubConvReq) (*CommonResp, error)
	//GetAllMsgList 获取所有消息
	GetAllMsgList(context.Context, *GetAllMsgListReq) (*GetAllMsgListResp, error)
	//ReadMsg 设置会话已读
	ReadMsg(context.Context, *ReadMsgReq) (*ReadMsgResp, error)
	//EditMsg 编辑消息
	EditMsg(context.Context, *EditMsgReq) (*EditMsgResp, error)
	mustEmbedUnimplementedMsgServiceServer()
}

// UnimplementedMsgServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServiceServer struct {
}

func (UnimplementedMsgServiceServer) InsertMsgDataList(context.Context, *MsgDataList) (*MsgDataList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertMsgDataList not implemented")
}
func (UnimplementedMsgServiceServer) SendMsgListSync(context.Context, *SendMsgListReq) (*SendMsgListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsgListSync not implemented")
}
func (UnimplementedMsgServiceServer) SendMsgListAsync(context.Context, *SendMsgListReq) (*SendMsgListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsgListAsync not implemented")
}
func (UnimplementedMsgServiceServer) PushMsgList(context.Context, *PushMsgListReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMsgList not implemented")
}
func (UnimplementedMsgServiceServer) BatchGetMsgListByConvId(context.Context, *BatchGetMsgListByConvIdReq) (*GetMsgListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetMsgListByConvId not implemented")
}
func (UnimplementedMsgServiceServer) GetMsgById(context.Context, *GetMsgByIdReq) (*GetMsgByIdResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMsgById not implemented")
}
func (UnimplementedMsgServiceServer) BatchSetMinSeq(context.Context, *BatchSetMinSeqReq) (*BatchSetMinSeqResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchSetMinSeq not implemented")
}
func (UnimplementedMsgServiceServer) BatchGetConvSeq(context.Context, *BatchGetConvSeqReq) (*BatchGetConvSeqResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetConvSeq not implemented")
}
func (UnimplementedMsgServiceServer) AfterConnect(context.Context, *AfterConnectReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AfterConnect not implemented")
}
func (UnimplementedMsgServiceServer) AfterDisconnect(context.Context, *AfterDisconnectReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AfterDisconnect not implemented")
}
func (UnimplementedMsgServiceServer) KeepAlive(context.Context, *KeepAliveReq) (*KeepAliveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (UnimplementedMsgServiceServer) GetConvSubscribers(context.Context, *GetConvSubscribersReq) (*GetConvSubscribersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConvSubscribers not implemented")
}
func (UnimplementedMsgServiceServer) OfflinePushMsg(context.Context, *OfflinePushMsgReq) (*OfflinePushMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OfflinePushMsg not implemented")
}
func (UnimplementedMsgServiceServer) GetConvOnlineCount(context.Context, *GetConvOnlineCountReq) (*GetConvOnlineCountResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConvOnlineCount not implemented")
}
func (UnimplementedMsgServiceServer) FlushUsersSubConv(context.Context, *FlushUsersSubConvReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FlushUsersSubConv not implemented")
}
func (UnimplementedMsgServiceServer) GetAllMsgList(context.Context, *GetAllMsgListReq) (*GetAllMsgListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMsgList not implemented")
}
func (UnimplementedMsgServiceServer) ReadMsg(context.Context, *ReadMsgReq) (*ReadMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadMsg not implemented")
}
func (UnimplementedMsgServiceServer) EditMsg(context.Context, *EditMsgReq) (*EditMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditMsg not implemented")
}
func (UnimplementedMsgServiceServer) mustEmbedUnimplementedMsgServiceServer() {}

// UnsafeMsgServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServiceServer will
// result in compilation errors.
type UnsafeMsgServiceServer interface {
	mustEmbedUnimplementedMsgServiceServer()
}

func RegisterMsgServiceServer(s grpc.ServiceRegistrar, srv MsgServiceServer) {
	s.RegisterService(&MsgService_ServiceDesc, srv)
}

func _MsgService_InsertMsgDataList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDataList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).InsertMsgDataList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/InsertMsgDataList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).InsertMsgDataList(ctx, req.(*MsgDataList))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_SendMsgListSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).SendMsgListSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/SendMsgListSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).SendMsgListSync(ctx, req.(*SendMsgListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_SendMsgListAsync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).SendMsgListAsync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/SendMsgListAsync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).SendMsgListAsync(ctx, req.(*SendMsgListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_PushMsgList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMsgListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).PushMsgList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/PushMsgList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).PushMsgList(ctx, req.(*PushMsgListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_BatchGetMsgListByConvId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetMsgListByConvIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).BatchGetMsgListByConvId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/BatchGetMsgListByConvId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).BatchGetMsgListByConvId(ctx, req.(*BatchGetMsgListByConvIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_GetMsgById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMsgByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).GetMsgById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/GetMsgById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).GetMsgById(ctx, req.(*GetMsgByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_BatchSetMinSeq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchSetMinSeqReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).BatchSetMinSeq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/BatchSetMinSeq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).BatchSetMinSeq(ctx, req.(*BatchSetMinSeqReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_BatchGetConvSeq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetConvSeqReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).BatchGetConvSeq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/BatchGetConvSeq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).BatchGetConvSeq(ctx, req.(*BatchGetConvSeqReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_AfterConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AfterConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).AfterConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/AfterConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).AfterConnect(ctx, req.(*AfterConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_AfterDisconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AfterDisconnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).AfterDisconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/AfterDisconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).AfterDisconnect(ctx, req.(*AfterDisconnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeepAliveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).KeepAlive(ctx, req.(*KeepAliveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_GetConvSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConvSubscribersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).GetConvSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/GetConvSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).GetConvSubscribers(ctx, req.(*GetConvSubscribersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_OfflinePushMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OfflinePushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).OfflinePushMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/OfflinePushMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).OfflinePushMsg(ctx, req.(*OfflinePushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_GetConvOnlineCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConvOnlineCountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).GetConvOnlineCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/GetConvOnlineCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).GetConvOnlineCount(ctx, req.(*GetConvOnlineCountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_FlushUsersSubConv_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlushUsersSubConvReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).FlushUsersSubConv(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/FlushUsersSubConv",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).FlushUsersSubConv(ctx, req.(*FlushUsersSubConvReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_GetAllMsgList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllMsgListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).GetAllMsgList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/GetAllMsgList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).GetAllMsgList(ctx, req.(*GetAllMsgListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_ReadMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).ReadMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/ReadMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).ReadMsg(ctx, req.(*ReadMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MsgService_EditMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServiceServer).EditMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.msgService/EditMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServiceServer).EditMsg(ctx, req.(*EditMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

// MsgService_ServiceDesc is the grpc.ServiceDesc for MsgService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MsgService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.msgService",
	HandlerType: (*MsgServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InsertMsgDataList",
			Handler:    _MsgService_InsertMsgDataList_Handler,
		},
		{
			MethodName: "SendMsgListSync",
			Handler:    _MsgService_SendMsgListSync_Handler,
		},
		{
			MethodName: "SendMsgListAsync",
			Handler:    _MsgService_SendMsgListAsync_Handler,
		},
		{
			MethodName: "PushMsgList",
			Handler:    _MsgService_PushMsgList_Handler,
		},
		{
			MethodName: "BatchGetMsgListByConvId",
			Handler:    _MsgService_BatchGetMsgListByConvId_Handler,
		},
		{
			MethodName: "GetMsgById",
			Handler:    _MsgService_GetMsgById_Handler,
		},
		{
			MethodName: "BatchSetMinSeq",
			Handler:    _MsgService_BatchSetMinSeq_Handler,
		},
		{
			MethodName: "BatchGetConvSeq",
			Handler:    _MsgService_BatchGetConvSeq_Handler,
		},
		{
			MethodName: "AfterConnect",
			Handler:    _MsgService_AfterConnect_Handler,
		},
		{
			MethodName: "AfterDisconnect",
			Handler:    _MsgService_AfterDisconnect_Handler,
		},
		{
			MethodName: "KeepAlive",
			Handler:    _MsgService_KeepAlive_Handler,
		},
		{
			MethodName: "GetConvSubscribers",
			Handler:    _MsgService_GetConvSubscribers_Handler,
		},
		{
			MethodName: "OfflinePushMsg",
			Handler:    _MsgService_OfflinePushMsg_Handler,
		},
		{
			MethodName: "GetConvOnlineCount",
			Handler:    _MsgService_GetConvOnlineCount_Handler,
		},
		{
			MethodName: "FlushUsersSubConv",
			Handler:    _MsgService_FlushUsersSubConv_Handler,
		},
		{
			MethodName: "GetAllMsgList",
			Handler:    _MsgService_GetAllMsgList_Handler,
		},
		{
			MethodName: "ReadMsg",
			Handler:    _MsgService_ReadMsg_Handler,
		},
		{
			MethodName: "EditMsg",
			Handler:    _MsgService_EditMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "msg.proto",
}
