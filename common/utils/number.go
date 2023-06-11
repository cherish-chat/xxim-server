package utils

import "strconv"

type xNumber struct {
}

var Number = &xNumber{}

func (x *xNumber) Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
