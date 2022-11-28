package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterDisconnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterDisconnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterDisconnectLogic {
	return &AfterDisconnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AfterDisconnectLogic) afterDisconnect(in *pb.AfterDisconnectReq) (*pb.CommonResp, error) {
	err := xorm.Update(l.svcCtx.Mysql(), &immodel.UserConnectRecord{}, map[string]interface{}{
		"disconnectTime": utils.AnyToInt64(in.DisconnectedAt),
	}, xorm.Where("userId = ? and deviceId = ? and disconnectTime = 0", in.ConnParam.UserId, in.ConnParam.DeviceId))
	if err != nil {
		l.Errorf("update connect record failed, err: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}

func (l *AfterDisconnectLogic) AfterDisconnect(in *pb.AfterDisconnectReq) (*pb.CommonResp, error) {
	var fs []func() error
	fs = append(fs, func() error {
		var err error
		xtrace.StartFuncSpan(l.ctx, "im.afterDisconnect", func(ctx context.Context) {
			_, err = l.afterDisconnect(in)
		})
		return err
	})
	fs = append(fs, func() error {
		var err error
		_, err = l.svcCtx.MsgService().AfterDisconnect(l.ctx, in)
		return err
	})
	err := mr.Finish(fs...)
	if err != nil {
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
