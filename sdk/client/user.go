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

// UserAccessToken 获取用户访问令牌
func (c *HttpClient) UserAccessToken(req *pb.UserAccessTokenReq) (resp *pb.UserAccessTokenResp, err error) {
	resp = &pb.UserAccessTokenResp{}
	err = c.Request("/v1/user/userAccessToken", req, resp)
	return
}

// UserAccessToken 获取用户访问令牌
func (c *WsClient) UserAccessToken(req *pb.UserAccessTokenReq) (resp *pb.UserAccessTokenResp, err error) {
	resp = &pb.UserAccessTokenResp{}
	err = c.Request("/v1/user/userAccessToken", req, resp)
	return
}
