package utils

import "time"

type xTime struct {
}

var Time = &xTime{}

func (x *xTime) Now13() int64 {
	now := time.Now()
	return now.UnixMilli()
}

func (x *xTime) Int64ToString(i int64) string {
	// 判断是不是10位
	if i <= 9999999999 && i >= 1000000000 {
		i = i * 1000
	}
	// 判断是不是19位
	if i >= 1000000000000000000 {
		i = i / 1000000
	}
	if i == 0 {
		return ""
	}
	return time.UnixMilli(i).Format("2006-01-02 15:04:05")
}
