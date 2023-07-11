package types

import (
	"github.com/cherish-chat/xxim-proto/peerpb"
	"google.golang.org/protobuf/proto"
)

func MarshalResponse(data proto.Message) []byte {
	protobuf, _ := proto.Marshal(data)
	return protobuf
}

func MarshalWriteData(data *peerpb.GatewayApiResponse) []byte {
	writeData := &peerpb.GatewayWriteDataContent{
		DataType: peerpb.GatewayWriteDataType_Response,
		Response: data,
		Message:  nil,
		Notice:   nil,
	}
	protobuf, _ := proto.Marshal(writeData)
	return protobuf
}
