package utils

import (
	"bytes"
	"io"
)

type xBytes struct {
}

var Bytes = &xBytes{}

func (x *xBytes) NewNopCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}
