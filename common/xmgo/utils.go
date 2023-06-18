package xmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrIsDup 判断mongo error是不是唯一索引冲突
func ErrIsDup(err error) bool {
	if writeException, ok := err.(mongo.WriteException); ok {
		for _, writeError := range writeException.WriteErrors {
			if writeError.Code == 11000 {
				return true
			}
		}
	} else {
		if bulkWriteException, ok := err.(*mongo.BulkWriteException); ok {
			for _, writeError := range bulkWriteException.WriteErrors {
				if writeError.Code == 11000 {
					return true
				}
			}
		}
	}
	return false
}

// BatchInsertMany 批量插入 防止一次插入太多数据 mongodb限制一次插入的数据量
func BatchInsertMany[T any](
	//集合
	collection *qmgo.QmgoClient,
	//上下文
	ctx context.Context,
	//数据
	models []T,
	//每次插入的数量
	batchSize int,
	//insertMany选项
	opts ...opts.InsertManyOptions,
) error {
	if batchSize <= 0 {
		batchSize = 1000
	}
	var (
		startIndex int
		endIndex   int
	)
	for {
		startIndex = endIndex
		endIndex = startIndex + batchSize
		if endIndex > len(models) {
			endIndex = len(models)
		}
		if startIndex == endIndex {
			break
		}
		_, err := collection.InsertMany(ctx, models[startIndex:endIndex], opts...)
		if err != nil {
			return err
		}
	}
	return nil
}
