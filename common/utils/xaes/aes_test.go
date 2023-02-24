package xaes

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	encryptBytes := Encrypt([]byte("iv"), []byte("key"), []byte("data"))
	t.Log(encryptBytes)
	bytes, _ := Decrypt([]byte("iv"), []byte("key"), encryptBytes)
	t.Log(string(bytes))
	bytes, err := Decrypt([]byte("iv"), []byte("key"), []byte("data"))
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	t.Log(string(bytes))
}
