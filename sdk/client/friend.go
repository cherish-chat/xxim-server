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

// ListFriendApply 列出好友申请
func (c *HttpClient) ListFriendApply(req *pb.ListFriendApplyReq) (resp *pb.ListFriendApplyResp, err error) {
	resp = &pb.ListFriendApplyResp{}
	err = c.Request("/v1/friend/listFriendApply", req, resp)
	return
}

// ListFriendApply 列出好友申请
func (c *WsClient) ListFriendApply(req *pb.ListFriendApplyReq) (resp *pb.ListFriendApplyResp, err error) {
	resp = &pb.ListFriendApplyResp{}
	err = c.Request("/v1/friend/listFriendApply", req, resp)
	return
}

// FriendApplyHandle 处理好友申请
func (c *HttpClient) FriendApplyHandle(req *pb.FriendApplyHandleReq) (resp *pb.FriendApplyHandleResp, err error) {
	resp = &pb.FriendApplyHandleResp{}
	err = c.Request("/v1/friend/friendApplyHandle", req, resp)
	return
}

// FriendApplyHandle 处理好友申请
func (c *WsClient) FriendApplyHandle(req *pb.FriendApplyHandleReq) (resp *pb.FriendApplyHandleResp, err error) {
	resp = &pb.FriendApplyHandleResp{}
	err = c.Request("/v1/friend/friendApplyHandle", req, resp)
	return
}
