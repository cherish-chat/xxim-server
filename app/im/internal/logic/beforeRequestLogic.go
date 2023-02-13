package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// allIpBlackList 为所有ip黑名单
var allIpBlackList [][3]string // startIp, endIp, userId

func InitAllIpBlackList(svcCtx *svc.ServiceContext) {
	syncAllIpBlackList(svcCtx)
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-ticker.C:
				syncAllIpBlackList(svcCtx)
			}
		}
	}()
}
func syncAllIpBlackList(svcCtx *svc.ServiceContext) {
	var ipBlackList []*usermodel.IpBlackList
	m := &usermodel.IpBlackList{}
	maxCreateTime := int64(0)
	// 1000条一次
	var tmpIpBlackList [][3]string
	for {
		err := svcCtx.Mysql().Model(m).
			Where("isEnable = ?", true).
			Where("createTime > ?", maxCreateTime).
			Order("createTime asc").
			Limit(1000).Find(&ipBlackList).Error
		if err != nil {
			panic(err)
		}
		if len(ipBlackList) == 0 {
			break
		}
		for _, v := range ipBlackList {
			tmpIpBlackList = append(tmpIpBlackList, [3]string{v.StartIp, v.EndIp, v.UserId})
			if v.CreateTime > maxCreateTime {
				maxCreateTime = v.CreateTime
			}
		}
	}
	allIpBlackList = tmpIpBlackList
}

type BeforeRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBeforeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BeforeRequestLogic {
	return &BeforeRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BeforeRequestLogic) BeforeRequest(in *pb.BeforeRequestReq) (*pb.BeforeRequestResp, error) {
	// 检查ip
	if in.CommonReq.UserId != "" {
		ip := in.CommonReq.Ip
		for _, v := range allIpBlackList {
			if v[2] == in.CommonReq.UserId || v[2] == "" {
				if strings.Compare(ip, v[0]) >= 0 && strings.Compare(ip, v[1]) <= 0 {
					return &pb.BeforeRequestResp{
						CommonResp: pb.NewAuthErrorResp("ip被封禁"),
					}, status.Error(codes.Unauthenticated, "ip被封禁")
				}
			}
		}
	}
	if strings.Contains(in.Method, "/white/") {
		return &pb.BeforeRequestResp{
			CommonResp: pb.NewSuccessResp(),
		}, nil
	}

	// 验证token
	code, msg := xjwt.VerifyToken(l.ctx, l.svcCtx.Redis(), in.CommonReq.UserId, in.CommonReq.Token, xjwt.WithPlatform(in.CommonReq.Platform), xjwt.WithDeviceId(in.CommonReq.DeviceId))
	switch code {
	case xjwt.VerifyTokenCodeInternalError, xjwt.VerifyTokenCodeError, xjwt.VerifyTokenCodeExpire, xjwt.VerifyTokenCodeBaned, xjwt.VerifyTokenCodeReplace:
		return &pb.BeforeRequestResp{CommonResp: pb.NewAuthErrorResp(msg)}, nil
	default:
		return &pb.BeforeRequestResp{}, nil
	}
}
