package types

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"google.golang.org/protobuf/proto"
)

func MarshalResponse(data proto.Message) []byte {
	protobuf, _ := proto.Marshal(data)
	return protobuf
}

func MarshalWriteData(data *pb.GatewayApiResponse) []byte {
	writeData := &pb.GatewayWriteDataContent{
		DataType: pb.GatewayWriteDataType_Response,
		Response: data,
		Message:  nil,
		Notice:   nil,
	}
	protobuf, _ := proto.Marshal(writeData)
	return protobuf
}
