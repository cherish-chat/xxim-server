package client

import "github.com/cherish-chat/xxim-server/sdk/types"

func (c *HttpClient) GatewayGetUserConnection(req *types.GatewayGetUserConnectionReq) (resp *types.GatewayGetUserConnectionResp, err error) {
	resp = &types.GatewayGetUserConnectionResp{}
	err = c.Request("/api/v1/gateway/getUserConnection", req, resp)
	return
}
