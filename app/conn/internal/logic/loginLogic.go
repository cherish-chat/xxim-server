package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type LoginLogic struct {
	svcCtx *svc.ServiceContext
}

var singletonLoginLogic *LoginLogic

func NewLoginLogic(svcCtx *svc.ServiceContext) *LoginLogic {
	if singletonLoginLogic == nil {
		singletonLoginLogic = &LoginLogic{
			svcCtx: svcCtx,
		}
	}
	return singletonLoginLogic
}

func (l *LoginLogic) Callback(ctx context.Context, resp *pb.LoginResp, c *types.UserConn) {
	c.SetConnParams(&pb.ConnParam{
		UserId:      resp.UserId,
		Token:       resp.Token,
		DeviceId:    c.ConnParam.DeviceId,
		Platform:    c.ConnParam.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: c.ConnParam.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		PodIp:       utils.GetPodIp(),
		DeviceModel: c.ConnParam.DeviceModel,
		OsVersion:   c.ConnParam.OsVersion,
		AppVersion:  c.ConnParam.AppVersion,
		Language:    c.ConnParam.Language,
	})
	if resp.Token != "" {
		// 登录成功
		GetConnLogic().AddSubscriber(c)
	}
}
