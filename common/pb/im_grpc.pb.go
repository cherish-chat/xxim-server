// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: im.proto

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

// ImServiceClient is the client API for ImService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImServiceClient interface {
	BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error)
	AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error)
	AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error)
	KickUserConn(ctx context.Context, in *KickUserConnReq, opts ...grpc.CallOption) (*KickUserConnResp, error)
	GetUserConn(ctx context.Context, in *GetUserConnReq, opts ...grpc.CallOption) (*GetUserConnResp, error)
	GetUserLatestConn(ctx context.Context, in *GetUserLatestConnReq, opts ...grpc.CallOption) (*GetUserLatestConnResp, error)
	BatchGetUserLatestConn(ctx context.Context, in *BatchGetUserLatestConnReq, opts ...grpc.CallOption) (*BatchGetUserLatestConnResp, error)
	SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error)
}

type imServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewImServiceClient(cc grpc.ClientConnInterface) ImServiceClient {
	return &imServiceClient{cc}
}

func (c *imServiceClient) BeforeConnect(ctx context.Context, in *BeforeConnectReq, opts ...grpc.CallOption) (*BeforeConnectResp, error) {
	out := new(BeforeConnectResp)
	err := c.cc.Invoke(ctx, "/pb.imService/BeforeConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) AfterConnect(ctx context.Context, in *AfterConnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.imService/AfterConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) AfterDisconnect(ctx context.Context, in *AfterDisconnectReq, opts ...grpc.CallOption) (*CommonResp, error) {
	out := new(CommonResp)
	err := c.cc.Invoke(ctx, "/pb.imService/AfterDisconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) KickUserConn(ctx context.Context, in *KickUserConnReq, opts ...grpc.CallOption) (*KickUserConnResp, error) {
	out := new(KickUserConnResp)
	err := c.cc.Invoke(ctx, "/pb.imService/KickUserConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) GetUserConn(ctx context.Context, in *GetUserConnReq, opts ...grpc.CallOption) (*GetUserConnResp, error) {
	out := new(GetUserConnResp)
	err := c.cc.Invoke(ctx, "/pb.imService/GetUserConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) GetUserLatestConn(ctx context.Context, in *GetUserLatestConnReq, opts ...grpc.CallOption) (*GetUserLatestConnResp, error) {
	out := new(GetUserLatestConnResp)
	err := c.cc.Invoke(ctx, "/pb.imService/GetUserLatestConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) BatchGetUserLatestConn(ctx context.Context, in *BatchGetUserLatestConnReq, opts ...grpc.CallOption) (*BatchGetUserLatestConnResp, error) {
	out := new(BatchGetUserLatestConnResp)
	err := c.cc.Invoke(ctx, "/pb.imService/BatchGetUserLatestConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imServiceClient) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error) {
	out := new(SendMsgResp)
	err := c.cc.Invoke(ctx, "/pb.imService/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImServiceServer is the server API for ImService service.
// All implementations must embed UnimplementedImServiceServer
// for forward compatibility
type ImServiceServer interface {
	BeforeConnect(context.Context, *BeforeConnectReq) (*BeforeConnectResp, error)
	AfterConnect(context.Context, *AfterConnectReq) (*CommonResp, error)
	AfterDisconnect(context.Context, *AfterDisconnectReq) (*CommonResp, error)
	KickUserConn(context.Context, *KickUserConnReq) (*KickUserConnResp, error)
	GetUserConn(context.Context, *GetUserConnReq) (*GetUserConnResp, error)
	GetUserLatestConn(context.Context, *GetUserLatestConnReq) (*GetUserLatestConnResp, error)
	BatchGetUserLatestConn(context.Context, *BatchGetUserLatestConnReq) (*BatchGetUserLatestConnResp, error)
	SendMsg(context.Context, *SendMsgReq) (*SendMsgResp, error)
	mustEmbedUnimplementedImServiceServer()
}

// UnimplementedImServiceServer must be embedded to have forward compatible implementations.
type UnimplementedImServiceServer struct {
}

func (UnimplementedImServiceServer) BeforeConnect(context.Context, *BeforeConnectReq) (*BeforeConnectResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeforeConnect not implemented")
}
func (UnimplementedImServiceServer) AfterConnect(context.Context, *AfterConnectReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AfterConnect not implemented")
}
func (UnimplementedImServiceServer) AfterDisconnect(context.Context, *AfterDisconnectReq) (*CommonResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AfterDisconnect not implemented")
}
func (UnimplementedImServiceServer) KickUserConn(context.Context, *KickUserConnReq) (*KickUserConnResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KickUserConn not implemented")
}
func (UnimplementedImServiceServer) GetUserConn(context.Context, *GetUserConnReq) (*GetUserConnResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserConn not implemented")
}
func (UnimplementedImServiceServer) GetUserLatestConn(context.Context, *GetUserLatestConnReq) (*GetUserLatestConnResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserLatestConn not implemented")
}
func (UnimplementedImServiceServer) BatchGetUserLatestConn(context.Context, *BatchGetUserLatestConnReq) (*BatchGetUserLatestConnResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetUserLatestConn not implemented")
}
func (UnimplementedImServiceServer) SendMsg(context.Context, *SendMsgReq) (*SendMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (UnimplementedImServiceServer) mustEmbedUnimplementedImServiceServer() {}

// UnsafeImServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImServiceServer will
// result in compilation errors.
type UnsafeImServiceServer interface {
	mustEmbedUnimplementedImServiceServer()
}

func RegisterImServiceServer(s grpc.ServiceRegistrar, srv ImServiceServer) {
	s.RegisterService(&ImService_ServiceDesc, srv)
}

func _ImService_BeforeConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeforeConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).BeforeConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/BeforeConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).BeforeConnect(ctx, req.(*BeforeConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_AfterConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AfterConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).AfterConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/AfterConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).AfterConnect(ctx, req.(*AfterConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_AfterDisconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AfterDisconnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).AfterDisconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/AfterDisconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).AfterDisconnect(ctx, req.(*AfterDisconnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_KickUserConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KickUserConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).KickUserConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/KickUserConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).KickUserConn(ctx, req.(*KickUserConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_GetUserConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).GetUserConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/GetUserConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).GetUserConn(ctx, req.(*GetUserConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_GetUserLatestConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserLatestConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).GetUserLatestConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/GetUserLatestConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).GetUserLatestConn(ctx, req.(*GetUserLatestConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_BatchGetUserLatestConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetUserLatestConnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).BatchGetUserLatestConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/BatchGetUserLatestConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).BatchGetUserLatestConn(ctx, req.(*BatchGetUserLatestConnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImService_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImServiceServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.imService/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImServiceServer).SendMsg(ctx, req.(*SendMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ImService_ServiceDesc is the grpc.ServiceDesc for ImService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.imService",
	HandlerType: (*ImServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BeforeConnect",
			Handler:    _ImService_BeforeConnect_Handler,
		},
		{
			MethodName: "AfterConnect",
			Handler:    _ImService_AfterConnect_Handler,
		},
		{
			MethodName: "AfterDisconnect",
			Handler:    _ImService_AfterDisconnect_Handler,
		},
		{
			MethodName: "KickUserConn",
			Handler:    _ImService_KickUserConn_Handler,
		},
		{
			MethodName: "GetUserConn",
			Handler:    _ImService_GetUserConn_Handler,
		},
		{
			MethodName: "GetUserLatestConn",
			Handler:    _ImService_GetUserLatestConn_Handler,
		},
		{
			MethodName: "BatchGetUserLatestConn",
			Handler:    _ImService_BatchGetUserLatestConn_Handler,
		},
		{
			MethodName: "SendMsg",
			Handler:    _ImService_SendMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "im.proto",
}
