package usermodel

import (
	"github.com/qiniu/qmgo"
)

type LoginRecord struct {
	Id        string `bson:"_id" json:"id"`
	UserId    string `bson:"userId" json:"userId"`
	LoginInfo `bson:",inline"`
}

func (m *LoginRecord) CollectionName() string {
	return "login_record"
}

func (m *LoginRecord) Indexes(c *qmgo.Collection) error {
	//TODO implement me
	return nil
}
