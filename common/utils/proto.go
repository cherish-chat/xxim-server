package utils

import (
	"google.golang.org/protobuf/proto"
)

type xProto struct {
}

var Proto = &xProto{}

func (x *xProto) Marshal(m proto.Message) []byte {
	data, _ := proto.Marshal(m)
	return data
}

func (x *xProto) Unmarshal(msg []byte, in proto.Message) error {
	return proto.Unmarshal(msg, in)
}
