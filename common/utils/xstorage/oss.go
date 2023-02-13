package xstorage

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cherish-chat/xxim-server/common/pb"
	"golang.org/x/net/context"
)

// OssStorage 阿里云OSS存储 实现Storage接口
type OssStorage struct {
	Config *pb.AppLineConfig_Storage_Oss
	client *oss.Client
	bucket *oss.Bucket
}

var singletonOssStorage *OssStorage

// NewOssStorage 创建OssStorage
func NewOssStorage(config *pb.AppLineConfig_Storage_Oss) (*OssStorage, error) {
	if singletonOssStorage != nil {
		return singletonOssStorage, nil
	}
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, err
	}
	o := &OssStorage{
		Config: config,
		client: client,
		bucket: bucket,
	}
	singletonOssStorage = o
	return singletonOssStorage, nil
}

// PutObject 上传文件
func (o *OssStorage) PutObject(ctx context.Context, objectName string, data []byte) (url string, err error) {
	var ioReader = bytes.NewReader(data)
	err = o.bucket.PutObject(objectName, ioReader)
	if err != nil {
		return "", err
	}
	return o.Config.BucketUrl + "/" + objectName, nil
}
