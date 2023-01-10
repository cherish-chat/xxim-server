package xerr

import (
	"errors"
	"strings"
)

func IsCanceled(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "context canceled") {
		return true
	}
	return false
}

var (
	InvalidParamError = errors.New("invalid param")
)
