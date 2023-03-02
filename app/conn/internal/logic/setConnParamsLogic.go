package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conn/internal/types"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xrsa"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

type SetConnParamsLogic struct {
	svcCtx *svc.ServiceContext
}

var singletonSetConnParamsLogic *SetConnParamsLogic

func NewSetConnParamsLogic(svcCtx *svc.ServiceContext) *SetConnParamsLogic {
	if singletonSetConnParamsLogic == nil {
		singletonSetConnParamsLogic = &SetConnParamsLogic{
			svcCtx: svcCtx,
		}
	}
	return singletonSetConnParamsLogic
}

func (l *SetConnParamsLogic) SetConnParams(ctx context.Context, req *pb.SetCxnParamsReq, opts ...grpc.CallOption) (*pb.SetCxnParamsResp, error) {
	return &pb.SetCxnParamsResp{
		Platform:    req.GetPlatform(),
		DeviceId:    req.GetPackageId(),
		DeviceModel: req.GetDeviceModel(),
		OsVersion:   req.GetOsVersion(),
		AppVersion:  req.GetAppVersion(),
		Language:    req.GetLanguage(),
		NetworkUsed: req.GetNetworkUsed(),
		Ext:         req.GetExt(),
		AesKey:      req.GetAesKey(),
		AesIv:       req.GetAesIv(),
	}, nil
}

func (l *SetConnParamsLogic) Callback(ctx context.Context, resp *pb.SetCxnParamsResp, c *types.UserConn) {
	// 验证
	_, ok := pb.PlatformMap[resp.Platform]
	if !ok {
		c.Conn.Close(types.WebsocketStatusCodePlatformFailed(), "platform failed")
		return
	}
	// rsa加密后的 aesKey
	aesKeyEncrypted := resp.GetAesKey()
	var aesKey *string
	// 是否不为空
	if len(aesKeyEncrypted) > 0 {
		// 解密
		decrypt, err := xrsa.Decrypt(aesKeyEncrypted, []byte(l.svcCtx.Config.RsaPrivateKey))
		if err != nil {
			// 断开连接
			c.Conn.Close(types.WebsocketStatusCodeRsaFailed(), "rsa decrypt failed")
			return
		}
		// 设置 aesKey
		aesKey = utils.AnyPtr(string(decrypt))
	}
	// rsa加密后的 aesIv
	aesIvEncrypted := resp.GetAesIv()
	var aesIv *string
	// 是否不为空
	if len(aesIvEncrypted) > 0 {
		decrypt, err := xrsa.Decrypt(aesIvEncrypted, []byte(l.svcCtx.Config.RsaPrivateKey))
		if err != nil {
			// 断开连接
			c.Conn.Close(types.WebsocketStatusCodeRsaFailed(), "rsa decrypt failed")
			return
		}
		// 设置 aesIv
		aesIv = utils.AnyPtr(string(decrypt))
	}
	if aesIv != nil && aesKey != nil {
		logx.Infof("aesKey: %s, aesIv: %s", *aesKey, *aesIv)
	}
	c.SetConnParams(&pb.ConnParam{
		UserId:      c.ConnParam.UserId,
		Token:       c.ConnParam.Token,
		DeviceId:    resp.DeviceId,
		Platform:    resp.Platform,
		Ips:         c.ConnParam.Ips,
		NetworkUsed: resp.NetworkUsed,
		Headers:     c.ConnParam.Headers,
		PodIp:       utils.GetPodIp(),
		DeviceModel: resp.DeviceModel,
		OsVersion:   resp.OsVersion,
		AppVersion:  resp.AppVersion,
		Language:    resp.Language,
		AesKey:      aesKey,
		AesIv:       aesIv,
	})
}
