package pb

func (x EncodingProto) ContentType() string {
	switch x {
	case EncodingProto_PROTOBUF:
		return "application/x-protobuf"
	case EncodingProto_JSON:
		return "application/json"
	default:
		return "application/x-protobuf"
	}
}
