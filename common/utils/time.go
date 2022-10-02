package utils

import "time"

func GetNowMilli() int64 {
	return time.Now().UnixMilli()
}
