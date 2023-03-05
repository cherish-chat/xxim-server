package xstorage

import "context"

// Storage 抽象存储接口
type Storage interface {
	// ExistObject 对象是否存在
	ExistObject(ctx context.Context, key string) (exists bool, err error)
	// PutObject 上传对象
	PutObject(ctx context.Context, key string, data []byte) (url string, err error)
	GetObjectUrl(key string) string
}
