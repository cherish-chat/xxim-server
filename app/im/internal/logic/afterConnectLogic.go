package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AfterConnectLogic) afterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	connectRecord := &immodel.UserConnectRecord{
		UserId:         in.ConnParam.UserId,
		DeviceId:       in.ConnParam.DeviceId,
		Platform:       in.ConnParam.Platform,
		Ips:            in.ConnParam.Ips,
		IpRegion:       ip2region.Ip2Region(in.ConnParam.Ips),
		NetworkUsed:    in.ConnParam.NetworkUsed,
		Headers:        in.ConnParam.Headers,
		PodIp:          in.ConnParam.PodIp,
		ConnectTime:    utils.AnyToInt64(in.ConnectedAt),
		DisconnectTime: 0,
	}
	_, err := l.svcCtx.Mongo().Collection(connectRecord).InsertOne(l.ctx, connectRecord)
	if err != nil {
		l.Errorf("insert connect record failed, err: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}

func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	var fs []func() error
	fs = append(fs, func() error {
		var err error
		xtrace.StartFuncSpan(l.ctx, "im.afterConnect", func(ctx context.Context) {
			_, err = l.afterConnect(in)
		})
		return err
	})
	fs = append(fs, func() error {
		var err error
		_, err = l.svcCtx.MsgService().AfterConnect(l.ctx, in)
		return err
	})
	err := mr.Finish(fs...)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
