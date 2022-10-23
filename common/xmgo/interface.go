package xmgo

import (
	"github.com/qiniu/qmgo"
)

type ICollection interface {
	CollectionName() string
	Indexes(c *qmgo.Collection) error
}
