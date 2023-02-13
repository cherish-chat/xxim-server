package xstorage

import (
	"bytes"
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// KodoStorage kodo存储 实现Storage接口
type KodoStorage struct {
	Config    *pb.AppLineConfig_Storage_Kodo
	mac       *qbox.Mac
	putPolicy *storage.PutPolicy
}

var singletonKodoStorage *KodoStorage

// NewKodoStorage 创建kodo存储
func NewKodoStorage(config *pb.AppLineConfig_Storage_Kodo) (*KodoStorage, error) {
	if singletonKodoStorage != nil {
		return singletonKodoStorage, nil
	}
	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	putPolicy := storage.PutPolicy{
		Scope: config.BucketName,
	}
	s := &KodoStorage{
		Config:    config,
		mac:       mac,
		putPolicy: &putPolicy,
	}
	singletonKodoStorage = s
	return s, nil
}

// PutObject 上传文件
func (s *KodoStorage) PutObject(ctx context.Context, objectName string, data []byte) (url string, err error) {
	upToken := s.putPolicy.UploadToken(s.mac)
	cfg := storage.Config{}
	cfg.UseHTTPS = s.Config.UseHTTPS
	cfg.Zone = &storage.ZoneHuadong
	switch s.Config.Region {
	case string(storage.RIDHuadong):
		cfg.Zone = &storage.ZoneHuadong
	case string(storage.RIDHuabei):
		cfg.Zone = &storage.ZoneHuabei
	case string(storage.RIDHuanan):
		cfg.Zone = &storage.ZoneHuanan
	case string(storage.RIDApNortheast1):
		cfg.Zone = &storage.ZoneShouEr1
	case string(storage.RIDSingapore):
		cfg.Zone = &storage.ZoneXinjiapo
	case string(storage.RIDNorthAmerica):
		cfg.Zone = &storage.ZoneBeimei
	default:
		cfg.Zone = nil
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err = formUploader.Put(ctx, &ret, upToken, objectName, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		return "", err
	}
	return s.Config.BucketUrl + "/" + objectName, nil
}
