package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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

func (l *AfterDisconnectLogic) AfterDisconnect(in *pb.AfterDisconnectReq) (*pb.CommonResp, error) {
	_, err := l.svcCtx.Mongo().Collection(&immodel.UserConnectRecord{}).UpdateAll(l.ctx, bson.M{
		"userId":         in.ConnParam.UserId,
		"deviceId":       in.ConnParam.DeviceId,
		"disconnectTime": 0,
	}, bson.M{
		"$set": bson.M{
			"disconnectTime": utils.AnyToInt64(in.DisconnectedAt),
		},
	})
	if err != nil {
		l.Errorf("update connect record failed, err: %v", err)
		return pb.NewRetryErrorResp(), err
	}
	return &pb.CommonResp{}, nil
}
