package xstorage

import "context"

// Storage 抽象存储接口
type Storage interface {
	// PutObject 上传对象
	PutObject(ctx context.Context, key string, data []byte) (url string, err error)
}
