package utils

import (
	"crypto/rand"
	"crypto/rsa"
)

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRsa(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *XRsa {
	return &XRsa{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (r *XRsa) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
}

func (r *XRsa) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, data)
}
