package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
)

func (c *HttpClient) GatewayGetUserConnection(req *pb.GatewayGetUserConnectionReq) (resp *pb.GatewayGetUserConnectionResp, err error) {
	resp = &pb.GatewayGetUserConnectionResp{}
	err = c.Request("/v1/gateway/getUserConnection", req, resp)
	return
}

func (c *WsClient) GatewayGetUserConnection(req *pb.GatewayGetUserConnectionReq) (resp *pb.GatewayGetUserConnectionResp, err error) {
	resp = &pb.GatewayGetUserConnectionResp{}
	err = c.Request("/v1/gateway/getUserConnection", req, resp)
	return
}
