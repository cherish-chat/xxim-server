package types

import (
	"google.golang.org/protobuf/proto"
)

type ReqInterface interface {
	proto.Message
	SetHeader(*RequestHeader)
}
