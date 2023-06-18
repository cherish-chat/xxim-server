package utils

import (
	"github.com/avast/retry-go"
	_ "github.com/avast/retry-go"
)

type xRetry struct {
}

var Retry = xRetry{}

func (x *xRetry) Do(f func() error) {
	_ = retry.Do(f, retry.Attempts(999))
}
