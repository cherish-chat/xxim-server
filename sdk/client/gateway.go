package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
)

// GatewayGetUserConnection 获取用户连接
func (c *HttpClient) GatewayGetUserConnection(req *pb.GatewayGetUserConnectionReq) (resp *pb.GatewayGetUserConnectionResp, err error) {
	resp = &pb.GatewayGetUserConnectionResp{}
	err = c.Request("/v1/gateway/getUserConnection", req, resp)
	return
}

// GatewayGetUserConnection 获取用户连接
func (c *WsClient) GatewayGetUserConnection(req *pb.GatewayGetUserConnectionReq) (resp *pb.GatewayGetUserConnectionResp, err error) {
	resp = &pb.GatewayGetUserConnectionResp{}
	err = c.Request("/v1/gateway/getUserConnection", req, resp)
	return
}

// GatewayBatchGetUserConnection 批量获取用户连接
func (c *HttpClient) GatewayBatchGetUserConnection(req *pb.GatewayBatchGetUserConnectionReq) (resp *pb.GatewayBatchGetUserConnectionResp, err error) {
	resp = &pb.GatewayBatchGetUserConnectionResp{}
	err = c.Request("/v1/gateway/batchGetUserConnection", req, resp)
	return
}

// GatewayBatchGetUserConnection 批量获取用户连接
func (c *WsClient) GatewayBatchGetUserConnection(req *pb.GatewayBatchGetUserConnectionReq) (resp *pb.GatewayBatchGetUserConnectionResp, err error) {
	resp = &pb.GatewayBatchGetUserConnectionResp{}
	err = c.Request("/v1/gateway/batchGetUserConnection", req, resp)
	return
}

// GatewayGetConnectionByFilter 根据条件获取连接
func (c *HttpClient) GatewayGetConnectionByFilter(req *pb.GatewayGetConnectionByFilterReq) (resp *pb.GatewayGetConnectionByFilterResp, err error) {
	resp = &pb.GatewayGetConnectionByFilterResp{}
	err = c.Request("/v1/gateway/getConnectionByFilter", req, resp)
	return
}

// GatewayGetConnectionByFilter 根据条件获取连接
func (c *WsClient) GatewayGetConnectionByFilter(req *pb.GatewayGetConnectionByFilterReq) (resp *pb.GatewayGetConnectionByFilterResp, err error) {
	resp = &pb.GatewayGetConnectionByFilterResp{}
	err = c.Request("/v1/gateway/getConnectionByFilter", req, resp)
	return
}

// GatewayWriteDataToWs 写入数据到ws
func (c *HttpClient) GatewayWriteDataToWs(req *pb.GatewayWriteDataToWsReq) (resp *pb.GatewayWriteDataToWsResp, err error) {
	resp = &pb.GatewayWriteDataToWsResp{}
	err = c.Request("/v1/gateway/writeDataToWs", req, resp)
	return
}

// GatewayWriteDataToWs 写入数据到ws
func (c *WsClient) GatewayWriteDataToWs(req *pb.GatewayWriteDataToWsReq) (resp *pb.GatewayWriteDataToWsResp, err error) {
	resp = &pb.GatewayWriteDataToWsResp{}
	err = c.Request("/v1/gateway/writeDataToWs", req, resp)
	return
}
