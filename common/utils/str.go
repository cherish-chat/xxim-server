package utils

import "encoding/json"

func Any2Bytes(i any) []byte {
	switch v := i.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		buf, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		return buf
	}
}
