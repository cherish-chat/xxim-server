package xmgo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Str2ObjectID(str string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(str)
	return id
}

func DuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(mongo.WriteException); ok {
		for _, writeError := range e.WriteErrors {
			return writeError.Code == 11000
		}
	}
	if e, ok := err.(mongo.BulkWriteException); ok {
		for _, writeError := range e.WriteErrors {
			if writeError.Code == 11000 {
				return true
			}
		}
	}
	return false
}
