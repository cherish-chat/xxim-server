package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xaes"
	"github.com/cherish-chat/xxim-server/common/utils/xstorage"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppLineConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppLineConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppLineConfigLogic {
	return &UpdateAppLineConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppLineConfigLogic) UpdateAppLineConfig(in *pb.UpdateAppLineConfigReq) (*pb.UpdateAppLineConfigResp, error) {
	// 校验参数
	if in.AppLineConfig == nil {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	valid := json.Valid([]byte(in.AppLineConfig.Config))
	if !valid {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	// 验证对象存储
	if in.AppLineConfig.Storage == nil {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	if in.AppLineConfig.Storage.Type == "" {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	if in.AppLineConfig.Storage.ObjectId == "" {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	var err error
	var storage xstorage.Storage
	// 验证对象存储配置
	switch in.AppLineConfig.Storage.Type {
	case "cos":
		if in.AppLineConfig.Storage.Cos == nil {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.AppId == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.SecretId == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.SecretKey == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.Region == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.BucketName == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Cos.BucketUrl == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		storage, err = xstorage.NewCosStorage(in.AppLineConfig.Storage.Cos)
	case "oss":
		if in.AppLineConfig.Storage.Oss == nil {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Oss.AccessKeyId == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Oss.AccessKeySecret == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Oss.Endpoint == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Oss.BucketName == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Oss.BucketUrl == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		storage, err = xstorage.NewOssStorage(in.AppLineConfig.Storage.Oss)
	case "minio":
		if in.AppLineConfig.Storage.Minio == nil {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Minio.Endpoint == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Minio.AccessKeyId == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Minio.SecretAccessKey == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Minio.BucketName == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Minio.BucketUrl == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		storage, err = xstorage.NewMinioStorage(in.AppLineConfig.Storage.Minio)
	case "kodo":
		if in.AppLineConfig.Storage.Kodo == nil {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Kodo.AccessKey == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Kodo.SecretKey == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Kodo.BucketName == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Kodo.Region == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		if in.AppLineConfig.Storage.Kodo.BucketUrl == "" {
			return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
		}
		storage, err = xstorage.NewKodoStorage(in.AppLineConfig.Storage.Kodo)
	default:
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, errors.New("参数错误")
	}
	if err != nil {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 加密 config
	encryptConfig := xaes.Encrypt([]byte(in.AppLineConfig.AesIv), []byte(in.AppLineConfig.AesKey), []byte(in.AppLineConfig.Config))
	// 保存到redis
	err = l.svcCtx.Redis().Set(rediskey.AppLineConfigKey(), utils.AnyToString(in.AppLineConfig))
	if err != nil {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	cfg := xorm.M{
		"config": utils.Base64.EncodeToString(encryptConfig),
	}
	bytes, _ := json.Marshal(cfg)
	// 上传 config
	_, err = storage.PutObject(l.ctx, in.AppLineConfig.Storage.ObjectId, bytes)
	if err != nil {
		return &pb.UpdateAppLineConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.UpdateAppLineConfigResp{}, nil
}
