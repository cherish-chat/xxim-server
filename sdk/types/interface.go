package types

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"google.golang.org/protobuf/proto"
)

type ReqInterface interface {
	proto.Message
	SetHeader(*pb.RequestHeader)
}
