package client

import "github.com/cherish-chat/xxim-server/common/pb"

// UserRegister 注册用户
func (c *HttpClient) UserRegister(req *pb.UserRegisterReq) (resp *pb.UserRegisterResp, err error) {
	resp = &pb.UserRegisterResp{}
	err = c.Request("/v1/user/white/userRegister", req, resp)
	return
}

// UserAccessToken 获取用户访问令牌
func (c *HttpClient) UserAccessToken(req *pb.UserAccessTokenReq) (resp *pb.UserAccessTokenResp, err error) {
	resp = &pb.UserAccessTokenResp{}
	err = c.Request("/v1/user/white/userAccessToken", req, resp)
	return
}

// CreateRobot 创建机器人
func (c *HttpClient) CreateRobot(req *pb.CreateRobotReq) (resp *pb.CreateRobotResp, err error) {
	resp = &pb.CreateRobotResp{}
	err = c.Request("/v1/user/createRobot", req, resp)
	return
}

// CreateRobot 创建机器人
func (c *WsClient) CreateRobot(req *pb.CreateRobotReq) (resp *pb.CreateRobotResp, err error) {
	resp = &pb.CreateRobotResp{}
	err = c.Request("/v1/user/createRobot", req, resp)
	return
}

// RefreshUserAccessToken 刷新用户访问令牌
func (c *HttpClient) RefreshUserAccessToken(req *pb.RefreshUserAccessTokenReq) (resp *pb.RefreshUserAccessTokenResp, err error) {
	resp = &pb.RefreshUserAccessTokenResp{}
	err = c.Request("/v1/user/refreshUserAccessToken", req, resp)
	return
}

// RefreshUserAccessToken 刷新用户访问令牌
func (c *WsClient) RefreshUserAccessToken(req *pb.RefreshUserAccessTokenReq) (resp *pb.RefreshUserAccessTokenResp, err error) {
	resp = &pb.RefreshUserAccessTokenResp{}
	err = c.Request("/v1/user/refreshUserAccessToken", req, resp)
	return
}
