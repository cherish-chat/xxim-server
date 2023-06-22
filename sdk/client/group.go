package client

import "github.com/cherish-chat/xxim-server/common/pb"

// GroupCreate 创建群
func (c *HttpClient) GroupCreate(req *pb.GroupCreateReq) (resp *pb.GroupCreateResp, err error) {
	resp = &pb.GroupCreateResp{}
	err = c.Request("/v1/group/groupCreate", req, resp)
	return
}

// GroupCreate 创建群
func (c *WsClient) GroupCreate(req *pb.GroupCreateReq) (resp *pb.GroupCreateResp, err error) {
	resp = &pb.GroupCreateResp{}
	err = c.Request("/v1/group/groupCreate", req, resp)
	return
}
