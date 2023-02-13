package utils

import "encoding/base64"

type Base64Util struct {
}

var Base64 = Base64Util{}

func (b Base64Util) EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
