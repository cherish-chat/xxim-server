package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppLineConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppLineConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppLineConfigLogic {
	return &GetAppLineConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppLineConfigLogic) GetAppLineConfig(in *pb.GetAppLineConfigReq) (*pb.GetAppLineConfigResp, error) {
	val, err := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.AppLineConfigKey())
	if err == redis.Nil || val == "" {
		return &pb.GetAppLineConfigResp{CommonResp: pb.NewSuccessResp(), AppLineConfig: &pb.AppLineConfig{
			Config: "",
			AesIv:  "",
			AesKey: "",
			Storage: &pb.AppLineConfig_Storage{
				Type:     "",
				ObjectId: "",
				Cos:      &pb.AppLineConfig_Storage_Cos{},
				Oss:      &pb.AppLineConfig_Storage_Oss{},
				Minio:    &pb.AppLineConfig_Storage_Minio{},
				Kodo:     &pb.AppLineConfig_Storage_Kodo{},
			},
		}}, nil
	}
	// json unmarshal
	var appLineConfig = &pb.AppLineConfig{}
	err = json.Unmarshal([]byte(val), appLineConfig)
	if err != nil {
		return &pb.GetAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppLineConfigResp{
		CommonResp:    pb.NewSuccessResp(),
		AppLineConfig: appLineConfig,
	}, nil
}
