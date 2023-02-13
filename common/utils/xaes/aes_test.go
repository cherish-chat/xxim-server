package xaes

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	encryptBytes := Encrypt([]byte("iv"), []byte("key"), []byte("data"))
	t.Log(encryptBytes)
	bytes := Decrypt([]byte("iv"), []byte("key"), encryptBytes)
	t.Log(string(bytes))
}
