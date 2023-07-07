package utils

import "testing"

func TestAes(t *testing.T) {
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	encrypted := Aes.Encrypt(key, iv, data)
	t.Logf("encrypted: %v", encrypted)
	decrypted, err := Aes.Decrypt(key, iv, encrypted)
	if err != nil {
		t.Errorf("Aes.Decrypt() error = %v", err)
		return
	}
	t.Logf("decrypted: %v", decrypted)
}
