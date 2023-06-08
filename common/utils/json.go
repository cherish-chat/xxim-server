package utils

import (
	"encoding/json"
	"fmt"
)

type xJson struct {
}

var Json = &xJson{}

func (x *xJson) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (x *xJson) MarshalToString(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}

func (x *xJson) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (x *xJson) MarshalToBytes(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte(fmt.Sprintf("%v", v))
	}
	return data
}
