package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/cherish-chat/xxim-server/common/utils"
)

var DecryptError = errors.New("decrypt error")

// Encrypt aes加密
func Encrypt(iv []byte, key []byte, data []byte) []byte {
	iv = []byte(utils.Md516(string(iv)))
	key = []byte(utils.Md5(string(key)))
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	data = pkcs7Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return encrypted
}

// Decrypt aes解密
func Decrypt(iv []byte, key []byte, data []byte) (decrypted []byte, err error) {
	iv = []byte(utils.Md516(string(iv)))
	key = []byte(utils.Md5(string(key)))
	block, _ := aes.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(data))
	defer func() {
		if r := recover(); r != nil {
			err = DecryptError
		}
	}()
	blockMode.CryptBlocks(decrypted, data)
	decrypted = pkcs7UnPadding(decrypted)
	return
}

// 使用PKCS7进行填充，IOS也是7
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
