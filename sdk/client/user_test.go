package client

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"testing"
	"time"
)

// UserRegister 注册用户
func TestHttpClient_UserRegister(t *testing.T) {
	userId := "9"
	phone := "13600000009"
	client := getHttpClient(t, nil)
	time.Sleep(1 * time.Second)
	userRegisterResp, err := client.(*HttpClient).UserRegister(&pb.UserRegisterReq{
		Header:       &pb.RequestHeader{},
		UserId:       userId,
		RegisterTime: nil,
		//Nickname:     utils.AnyPtr("用户" + userId),
		//Avatar:       utils.AnyPtr("https://www.baidu.com"),
		AccountMap: map[string]string{
			pb.AccountTypeUsername:     "user" + userId,
			pb.AccountTypePassword:     utils.Md5("123456"),
			pb.AccountTypePasswordSalt: utils.Snowflake.String(),
			pb.AccountTypePhone:        phone,
			pb.AccountTypePhoneCode:    "86",
			pb.AccountTypeEmail:        phone + "@xxim.com",
		},
		ProfileMap: map[string]string{
			"birthday": "2020-01-01",
			"sex":      "1",
		},
		ExtraMap: map[string]string{
			"extra1": "extra1",
		},
		VerifyMap: map[string]string{
			pb.AccountVerifyTypeSmsCode:     "123456",
			pb.AccountVerifyTypeEmailCode:   "123456",
			pb.AccountVerifyTypeCaptchaId:   "123456",
			pb.AccountVerifyTypeCaptchaCode: "123456",
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(userRegisterResp))
}

// UserAccessToken 获取用户访问令牌
func TestHttpClient_UserAccessToken(t *testing.T) {
	client := getHttpClient(t, nil)
	time.Sleep(1 * time.Second)
	userAccessTokenResp, err := client.(*HttpClient).UserAccessToken(&pb.UserAccessTokenReq{
		AccountMap: map[string]string{
			pb.AccountTypeUsername: "user4",
			pb.AccountTypePassword: utils.Md5("123456"),
		},
		VerifyMap:  nil,
		ExpireTime: nil,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(userAccessTokenResp))
}

// CreateRobot 创建机器人
func TestHttpClient_CreateRobot(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	createRobotResp, err := client.CreateRobot(&pb.CreateRobotReq{
		RobotId:  "robot11",
		Nickname: utils.AnyPtr("机器人11"),
		Avatar:   utils.AnyPtr("https://www.baidu.com"),
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(createRobotResp))
}

// RefreshUserAccessToken 刷新用户访问令牌
func TestHttpClient_RefreshUserAccessToken(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	refreshUserAccessTokenResp, err := client.RefreshUserAccessToken(&pb.RefreshUserAccessTokenReq{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(refreshUserAccessTokenResp))
}

// RevokeUserAccessToken 注销用户token
func TestHttpClient_RevokeUserAccessToken(t *testing.T) {
	client := getHttpClient(t, nil)
	//client := getWsClient(t, nil)
	time.Sleep(1 * time.Second)
	revokeUserAccessTokenResp, err := client.RevokeUserAccessToken(&pb.RevokeUserAccessTokenReq{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%s", utils.Json.MarshalToString(revokeUserAccessTokenResp))
}
