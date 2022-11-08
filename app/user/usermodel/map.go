package usermodel

import "go.mongodb.org/mongo-driver/bson"

type M bson.M

func (m M) Get(key string, defaultValue string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return defaultValue
}
