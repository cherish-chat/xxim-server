package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// rsa 加解密 publicKey: spki格式的公钥 privateKey: pkcs8格式的私钥  使用PKCS1v15填充

// Decrypt 解密
func Decrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	// 解密pem格式的私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("private key error")
	}

	// 解析PKCS8格式的私钥
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言
	pri := priv.(*rsa.PrivateKey)

	return rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext)
}

// Encrypt 加密
func Encrypt(plaintext []byte, publicKey []byte) ([]byte, error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, fmt.Errorf("public key error!")
	}

	// 解析PKIX格式的公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}
