package xstorage

import (
	"bytes"
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage minio存储 实现Storage接口
type MinioStorage struct {
	Config *pb.AppLineConfig_Storage_Minio
	Client *minio.Client
}

var singletonMinioStorage *MinioStorage

// NewMinioStorage 创建minio存储
func NewMinioStorage(config *pb.AppLineConfig_Storage_Minio) (*MinioStorage, error) {
	if singletonMinioStorage != nil {
		return singletonMinioStorage, nil
	}
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyId, config.SecretAccessKey, ""),
		Secure: config.Ssl,
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}
	s := &MinioStorage{
		Config: config,
		Client: client,
	}
	singletonMinioStorage = s
	return s, nil
}

// PutObject 上传文件
func (s *MinioStorage) PutObject(ctx context.Context, objectName string, data []byte) (url string, err error) {
	ioReader := bytes.NewReader(data)
	objectSize := int64(len(data))
	_, err = s.Client.PutObject(ctx, s.Config.BucketName, objectName, ioReader, objectSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return s.Config.BucketUrl + "/" + objectName, nil
}
