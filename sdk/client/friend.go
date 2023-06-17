package client

import "github.com/cherish-chat/xxim-server/common/pb"

// FriendApply 申请好友
func (c *HttpClient) FriendApply(req *pb.FriendApplyReq) (resp *pb.FriendApplyResp, err error) {
	resp = &pb.FriendApplyResp{}
	err = c.Request("/v1/friend/friendApply", req, resp)
	return
}

// FriendApply 申请好友
func (c *WsClient) FriendApply(req *pb.FriendApplyReq) (resp *pb.FriendApplyResp, err error) {
	resp = &pb.FriendApplyResp{}
	err = c.Request("/v1/friend/friendApply", req, resp)
	return
}
