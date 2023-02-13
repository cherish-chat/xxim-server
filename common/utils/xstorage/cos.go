package xstorage

import (
	"bytes"
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

// CosStorage 腾讯云COS存储 实现Storage接口
type CosStorage struct {
	Config *pb.AppLineConfig_Storage_Cos
	bucket *cos.Client
}

var singletonCosStorage *CosStorage

// NewCosStorage 创建CosStorage
func NewCosStorage(config *pb.AppLineConfig_Storage_Cos) (*CosStorage, error) {
	if singletonCosStorage != nil {
		return singletonCosStorage, nil
	}
	bucketUrl := config.BucketUrl
	// 转换成 url.Url
	parse, err := url.Parse(bucketUrl)
	if err != nil {
		return nil, err
	}
	bucket := cos.NewClient(&cos.BaseURL{BucketURL: parse}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretId,
			SecretKey: config.SecretKey,
		},
	})
	o := &CosStorage{
		Config: config,
		bucket: bucket,
	}
	singletonCosStorage = o
	return singletonCosStorage, nil
}

// PutObject 上传文件
func (o *CosStorage) PutObject(ctx context.Context, objectName string, data []byte) (url string, err error) {
	ioReader := bytes.NewReader(data)
	_, err = o.bucket.Object.Put(ctx, objectName, ioReader, nil)
	if err != nil {
		return "", err
	}
	return o.Config.BucketUrl + "/" + objectName, nil
}
