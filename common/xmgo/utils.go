package xmgo

import "go.mongodb.org/mongo-driver/bson/primitive"

func Str2ObjectID(str string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(str)
	return id
}
