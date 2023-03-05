package xaes

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	encryptBytes := Encrypt([]byte("iv"), []byte("key"), []byte("data"))
	t.Log(encryptBytes)
	bytes, _ := Decrypt([]byte("iv"), []byte("key"), encryptBytes)
	t.Log(string(bytes))
	bytes, err := Decrypt([]byte("iv"), []byte("key"), []byte("data"))
	assert.Assert(t, err != nil, "err should not be nil")
	t.Log(string(bytes))
}
