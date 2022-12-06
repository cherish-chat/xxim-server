package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"

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
	var headers = xorm.M{}
	for k, v := range in.ConnParam.Headers {
		headers[k] = v
	}
	connectRecord := &immodel.UserConnectRecord{
		UserId:         in.ConnParam.UserId,
		DeviceId:       in.ConnParam.DeviceId,
		Platform:       in.ConnParam.Platform,
		Ips:            in.ConnParam.Ips,
		IpRegion:       ip2region.Ip2Region(in.ConnParam.Ips),
		NetworkUsed:    in.ConnParam.NetworkUsed,
		Headers:        headers,
		PodIp:          in.ConnParam.PodIp,
		ConnectTime:    utils.AnyToInt64(in.ConnectedAt),
		DisconnectTime: 0,
	}
	// 判断是否存在 通过 ConnectTime
	err := l.svcCtx.Mysql().Model(&immodel.UserConnectRecord{}).Where("user_id = ? and device_id = ? and connect_time = ?", in.ConnParam.UserId, in.ConnParam.DeviceId, utils.AnyToInt64(in.ConnectedAt)).First(connectRecord).Error
	if xorm.RecordNotFound(err) {
		err := xorm.InsertOne(l.svcCtx.Mysql(), connectRecord)
		if err != nil {
			l.Errorf("insert connect record failed, err: %v", err)
			return pb.NewRetryErrorResp(), err
		}
	}
	// 写入redis latest connect record
	err = l.svcCtx.Redis().SetexCtx(l.ctx, rediskey.LatestConnectRecord(in.ConnParam.UserId), utils.AnyToString(connectRecord), rediskey.LatestConnectRecordExpire())
	if err != nil {
		l.Errorf("set latest connect record failed, err: %v", err)
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
	fs = append(fs, func() error {
		var err error
		_, err = l.svcCtx.NoticeService().AfterConnect(l.ctx, in)
		return err
	})
	err := mr.Finish(fs...)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
