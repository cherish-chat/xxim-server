package client

import "github.com/cherish-chat/xxim-server/common/pb"

// UserRegister 注册用户
func (c *HttpClient) UserRegister(req *pb.UserRegisterReq) (resp *pb.UserRegisterResp, err error) {
	resp = &pb.UserRegisterResp{}
	err = c.Request("/v1/user/userRegister", req, resp)
	return
}

// UserRegister 注册用户
func (c *WsClient) UserRegister(req *pb.UserRegisterReq) (resp *pb.UserRegisterResp, err error) {
	resp = &pb.UserRegisterResp{}
	err = c.Request("/v1/user/userRegister", req, resp)
	return
}
