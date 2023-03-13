package xstorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/pb"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage minio存储 实现Storage接口
type MinioStorage struct {
	Config *pb.AppLineConfig_Storage_Minio
	Client *minio.Client
}

func (s *MinioStorage) ExistObject(ctx context.Context, key string) (exists bool, err error) {
	_, err = s.Client.StatObject(ctx, s.Config.BucketName, key, minio.StatObjectOptions{})
	if err != nil {
		e, ok := err.(minio.ErrorResponse)
		if ok && e.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
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

func (s *MinioStorage) GetObjectUrl(key string) string {
	return fmt.Sprintf("%s/%s", s.Config.BucketUrl, key)
}

func (s *MinioStorage) setDownloadPolicy() error {
	// mc policy set download xxim
	err := s.Client.SetBucketPolicy(context.Background(), s.Config.BucketName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::`+s.Config.BucketName+`/*"]}]}`)
	if err != nil {
		return err
	}
	return nil
}
