package utils

import "errors"

type xError struct {
}

var Error = &xError{}

func (x *xError) DeepUnwrap(err error) error {
	e := errors.Unwrap(err)
	if e == nil {
		return err
	} else {
		return x.DeepUnwrap(e)
	}
}
