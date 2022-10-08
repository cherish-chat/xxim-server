package tests

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/xx/xxservice"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
)

var ctx = context.Background()

func xxService() xxservice.XxService {
	return xxservice.NewXxService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{
			"127.0.0.1:25000",
		},
	}))
}

// TestRegisterLogic1 注册
func TestRegisterLogic1(t *testing.T) {
	selfId := "TestSendMsgLogic1"
	xxService := xxService()
	resp, err := xxService.Register(ctx, &pb.RegisterReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
		},
		UserData: &pb.UserData{
			Id:           selfId,
			Nickname:     "TestSendMsgLogic1",
			Avatar:       "https://www.baidu.com",
			Xb:           "1",
			Birthday:     "2000-01-01",
			Signature:    "TestSendMsgLogic1",
			Tags:         []string{"TestSendMsgLogic1"},
			Password:     "123456",
			RegisterInfo: nil,
			IsRobot:      false,
			IsGuest:      false,
			IsAdmin:      false,
			IsOfficial:   false,
			UnbanTime:    "",
			AdminRemark:  "",
			Ex:           nil,
		},
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestGetUserLogic1 获取用户信息
func TestGetUserLogic1(t *testing.T) {
	selfId := "TestSendMsgLogic1"
	xxService := xxService()
	resp, err := xxService.GetUser(ctx, &pb.GetUserReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
		},
		UserIdList: []string{selfId},
	})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestLoginLogic1 登录
func TestLoginLogic1(t *testing.T) {
	selfId := "TestLoginLogic1"
	xxService := xxService()
	resp, err := xxService.Login(ctx, &pb.LoginReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
			DeviceId:    "TestLoginLogic1",
		},
		UserId:   "",
		Password: "",
	})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestLoginLogic2 登录
func TestLoginLogic2(t *testing.T) {
	selfId := "TestLoginLogic2"
	xxService := xxService()
	resp, err := xxService.Login(ctx, &pb.LoginReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
			DeviceId:    "TestLoginLogic1",
		},
		UserId:   "TestSendMsgLogic1",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestUpdateUserLogic1 更新用户信息
func TestUpdateUserLogic1(t *testing.T) {
	selfId := "TestSendMsgLogic1"
	xxService := xxService()
	resp, err := xxService.UpdateUser(ctx, &pb.UpdateUserReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
			DeviceId:    "TestLoginLogic1",
		},
		UserData: &pb.UserData{
			Id:           selfId,
			Nickname:     "修改后的昵称",
			Avatar:       "https://www.baidu.com",
			Xb:           "1",
			Birthday:     "2000-01-02",
			Signature:    "修改后的签名",
			Tags:         []string{"修改后的标签"},
			Password:     "1234567",
			RegisterInfo: nil,
			IsRobot:      false,
			IsGuest:      false,
			IsAdmin:      false,
			IsOfficial:   false,
			UnbanTime:    "",
			AdminRemark:  "",
			Ex:           nil,
		},
	})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}

// TestSearchUserLogic1 搜索用户
func TestSearchUserLogic1(t *testing.T) {
	selfId := "TestSendMsgLogic1"
	xxService := xxService()
	resp, err := xxService.SearchUser(ctx, &pb.SearchUserReq{
		Base: &pb.BaseReq{
			SelfId:      selfId,
			Platform:    "Test",
			AppVersion:  "v0.0.1",
			DeviceModel: "Mac",
			Ips:         "123.113.102.114",
			DeviceId:    "TestLoginLogic1",
		},
		Keyword:  "称",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	t.Logf("resp: %+v", resp.String())
}
