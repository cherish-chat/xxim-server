package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type xAes struct {
}

var Aes = &xAes{}

var DecryptError = errors.New("aes decrypt error")

func (x *xAes) Decrypt(key, iv string, data []byte) (decrypted []byte, err error) {
	aesKey := []byte(Md5(key))
	aesIv := []byte(Md516(iv))
	block, _ := aes.NewCipher(aesKey)
	blockMode := cipher.NewCBCDecrypter(block, aesIv)
	decrypted = make([]byte, len(data))
	defer func() {
		if r := recover(); r != nil {
			err = DecryptError
		}
	}()
	blockMode.CryptBlocks(decrypted, data)
	decrypted = x.pkcs7UnPadding(decrypted)
	return
}

// Encrypt aes加密
func (x *xAes) Encrypt(key string, iv string, data []byte) []byte {
	aesKey := []byte(Md5(key))
	aesIv := []byte(Md516(iv))
	block, _ := aes.NewCipher(aesKey)
	blockSize := block.BlockSize()
	data = x.pkcs7Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, aesIv)
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return encrypted
}

// 使用PKCS7进行填充，IOS也是7
func (x *xAes) pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	repeat := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, repeat...)
}

func (x *xAes) pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
