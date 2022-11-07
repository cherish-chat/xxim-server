package xjwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"strings"
)

func tokenClaimsString(token string) string {
	split := strings.Split(token, ".")
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func rsaKeyGenToKey(bits int) *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
	return privateKey
}

func pubKeyToBytes(key *rsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(key))
}

func bytesToPubKey(base64Ket string) (interface{}, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Ket)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PublicKey(keyBytes)
}
