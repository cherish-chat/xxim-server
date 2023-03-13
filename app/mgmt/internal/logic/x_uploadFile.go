package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils/xstorage"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UploadFileLogic) UploadFile(key string, data []byte) (string, error) {
	var configResp *pb.GetAppLineConfigResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetAppLineConfig", func(ctx context.Context) {
		configResp, err = NewGetAppLineConfigLogic(l.ctx, l.svcCtx).GetAppLineConfig(&pb.GetAppLineConfigReq{})
	})
	if err != nil {
		l.Errorf("GetAppLineConfig failed: %v", err)
		return "", err
	}
	appLineConfigClass := &mgmtmodel.AppLineConfigClass{}
	err = json.Unmarshal([]byte(configResp.AppLineConfig.Config), appLineConfigClass)
	if err != nil {
		l.Errorf("json.Unmarshal failed: %v", err)
		return "", err
	}
	var storage xstorage.Storage
	// 验证对象存储配置
	switch appLineConfigClass.ObjectStorage.Type {
	case "cos":
		storage, err = xstorage.NewCosStorage(&pb.AppLineConfig_Storage_Cos{
			AppId:      appLineConfigClass.ObjectStorage.Cos.AppId,
			SecretId:   appLineConfigClass.ObjectStorage.Cos.SecretId,
			SecretKey:  appLineConfigClass.ObjectStorage.Cos.SecretKey,
			BucketName: appLineConfigClass.ObjectStorage.Cos.BucketName,
			Region:     appLineConfigClass.ObjectStorage.Cos.Region,
			BucketUrl:  appLineConfigClass.ObjectStorage.Cos.BucketUrl,
		})
	case "oss":
		storage, err = xstorage.NewOssStorage(&pb.AppLineConfig_Storage_Oss{
			Endpoint:        appLineConfigClass.ObjectStorage.Oss.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Oss.AccessKeyId,
			AccessKeySecret: appLineConfigClass.ObjectStorage.Oss.AccessKeySecret,
			BucketName:      appLineConfigClass.ObjectStorage.Oss.BucketName,
			BucketUrl:       appLineConfigClass.ObjectStorage.Oss.BucketUrl,
		})
	case "minio":
		storage, err = xstorage.NewMinioStorage(&pb.AppLineConfig_Storage_Minio{
			Endpoint:        appLineConfigClass.ObjectStorage.Minio.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Minio.AccessKeyId,
			SecretAccessKey: appLineConfigClass.ObjectStorage.Minio.SecretAccessKey,
			BucketName:      appLineConfigClass.ObjectStorage.Minio.BucketName,
			Ssl:             appLineConfigClass.ObjectStorage.Minio.SSL,
			BucketUrl:       appLineConfigClass.ObjectStorage.Minio.BucketUrl,
		})
	default:
		l.Errorf("请配置app对象存储")
		return "", errors.New("请配置app对象存储")
	}
	if exist, err := storage.ExistObject(l.ctx, key); err == nil && exist {
		return storage.GetObjectUrl(key), nil
	}
	return storage.PutObject(l.ctx, key, data)
}

func (l *UploadFileLogic) MayGetUrl(key string) string {
	if strings.HasPrefix(key, "http") {
		return key
	}
	var configResp *pb.GetAppLineConfigResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "GetAppLineConfig", func(ctx context.Context) {
		configResp, err = NewGetAppLineConfigLogic(l.ctx, l.svcCtx).GetAppLineConfig(&pb.GetAppLineConfigReq{})
	})
	if err != nil {
		l.Errorf("GetAppLineConfig failed: %v", err)
		return key
	}
	appLineConfigClass := &mgmtmodel.AppLineConfigClass{}
	err = json.Unmarshal([]byte(configResp.AppLineConfig.Config), appLineConfigClass)
	if err != nil {
		l.Errorf("json.Unmarshal failed: %v", err)
		return key
	}
	var storage xstorage.Storage
	// 验证对象存储配置
	switch appLineConfigClass.ObjectStorage.Type {
	case "cos":
		storage, err = xstorage.NewCosStorage(&pb.AppLineConfig_Storage_Cos{
			AppId:      appLineConfigClass.ObjectStorage.Cos.AppId,
			SecretId:   appLineConfigClass.ObjectStorage.Cos.SecretId,
			SecretKey:  appLineConfigClass.ObjectStorage.Cos.SecretKey,
			BucketName: appLineConfigClass.ObjectStorage.Cos.BucketName,
			Region:     appLineConfigClass.ObjectStorage.Cos.Region,
			BucketUrl:  appLineConfigClass.ObjectStorage.Cos.BucketUrl,
		})
	case "oss":
		storage, err = xstorage.NewOssStorage(&pb.AppLineConfig_Storage_Oss{
			Endpoint:        appLineConfigClass.ObjectStorage.Oss.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Oss.AccessKeyId,
			AccessKeySecret: appLineConfigClass.ObjectStorage.Oss.AccessKeySecret,
			BucketName:      appLineConfigClass.ObjectStorage.Oss.BucketName,
			BucketUrl:       appLineConfigClass.ObjectStorage.Oss.BucketUrl,
		})
	case "minio":
		storage, err = xstorage.NewMinioStorage(&pb.AppLineConfig_Storage_Minio{
			Endpoint:        appLineConfigClass.ObjectStorage.Minio.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Minio.AccessKeyId,
			SecretAccessKey: appLineConfigClass.ObjectStorage.Minio.SecretAccessKey,
			BucketName:      appLineConfigClass.ObjectStorage.Minio.BucketName,
			Ssl:             appLineConfigClass.ObjectStorage.Minio.SSL,
			BucketUrl:       appLineConfigClass.ObjectStorage.Minio.BucketUrl,
		})
	default:
		l.Errorf("请配置app对象存储")
		return key
	}
	return storage.GetObjectUrl(key)
}

func (l *UploadFileLogic) AlbumAdd(cid uint, filename string, url string, suffix string, size int64, adminId string, typ int) (*mgmtmodel.Album, error) {
	album := &mgmtmodel.Album{
		Cid:  cid,
		Aid:  adminId,
		Type: typ,
		Name: filename,
		Url:  url,
		Ext:  suffix,
		Size: size,
	}
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Album{}).Create(album).Error
	if err != nil {
		l.Errorf("AlbumAdd failed: %v", err)
		return album, err
	}
	return album, nil
}
