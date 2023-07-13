package utils

import (
	"crypto"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"math/big"
)

// ECDH The main interface for ECDH key exchange.
type ECDH interface {
	// GenerateKey 生成公私钥对
	GenerateKey(io.Reader) (crypto.PrivateKey, crypto.PublicKey, error)
	// Marshal 将公钥转换为字节
	Marshal(crypto.PublicKey) []byte
	// MarshalHex 将公钥转换为hex
	MarshalHex(crypto.PublicKey) string
	// Unmarshal 将字节转换为公钥
	Unmarshal([]byte) (crypto.PublicKey, bool)
	// UnmarshalHex 将hex转换为公钥
	UnmarshalHex(hex string) (crypto.PublicKey, bool)
	// GenerateSharedSecret 生成共享密钥
	GenerateSharedSecret(crypto.PrivateKey, crypto.PublicKey) ([]byte, error)
	// HexEncodePublicKeyToString 将公钥钥转换为hex字符串
	HexEncodePublicKeyToString(crypto.PublicKey) string
	// HexEncodePrivateKeyToString 将私钥转换为hex字符串
	HexEncodePrivateKeyToString(crypto.PrivateKey) string
}

type ecdh struct {
	ECDH
	curve elliptic.Curve
}

// ecdhPublicKey The public key for ecdh.
type ecdhPublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

// ecdhPrivateKey The private key for ecdh.
type ecdhPrivateKey struct {
	D []byte
}

// NewECDH creates a new instance of ECDH with the given elliptic.Curve curve
// to use as the elliptical curve for elliptical curve diffie-hellman.
func NewECDH(curve elliptic.Curve) ECDH {
	return &ecdh{
		curve: curve,
	}
}

// GenerateKey generates a public and private key pair.
func (e *ecdh) GenerateKey(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var d []byte
	var x, y *big.Int
	var privateKey *ecdhPrivateKey
	var publicKey *ecdhPublicKey
	var err error

	d, x, y, err = elliptic.GenerateKey(e.curve, rand)
	if err != nil {
		return nil, nil, err
	}

	privateKey = &ecdhPrivateKey{
		D: d,
	}
	publicKey = &ecdhPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}

	return privateKey, publicKey, nil
}

// Marshal converts a public key to bytes.
func (e *ecdh) Marshal(publicKey crypto.PublicKey) []byte {
	key := publicKey.(*ecdhPublicKey)
	return elliptic.Marshal(e.curve, key.X, key.Y)
}

// MarshalHex converts a public key to bytes.
func (e *ecdh) MarshalHex(publicKey crypto.PublicKey) string {
	return hex.EncodeToString(e.Marshal(publicKey))
}

// Unmarshal converts bytes to a public key.
func (e *ecdh) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var key *ecdhPublicKey
	var x, y *big.Int

	x, y = elliptic.Unmarshal(e.curve, data)
	if x == nil || y == nil {
		return key, false
	}
	key = &ecdhPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}
	return key, true
}

// UnmarshalHex converts hex to a public key.
func (e *ecdh) UnmarshalHex(h string) (crypto.PublicKey, bool) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		logx.Errorf("UnmarshalHex error: %v", err)
		return nil, false
	}
	return e.Unmarshal(bytes)
}

// GenerateSharedSecret takes in a public key and a private key
// and generates a shared secret.
//
// RFC5903 Section 9 states we should only return x.
func (e *ecdh) GenerateSharedSecret(privateKey crypto.PrivateKey, publicKey crypto.PublicKey) ([]byte, error) {
	private := privateKey.(*ecdhPrivateKey)
	public := publicKey.(*ecdhPublicKey)

	x, _ := e.curve.ScalarMult(public.X, public.Y, private.D)
	return x.Bytes(), nil
}

// HexEncodePublicKeyToString 将公钥钥转换为hex字符串
func (e *ecdh) HexEncodePublicKeyToString(publicKey crypto.PublicKey) string {
	return hex.EncodeToString(e.Marshal(publicKey))
}

// HexEncodePrivateKeyToString 将私钥转换为hex字符串
func (e *ecdh) HexEncodePrivateKeyToString(privateKey crypto.PrivateKey) string {
	return hex.EncodeToString(privateKey.(*ecdhPrivateKey).D)
}

// GetECDHSecret 使用自己的私钥和别人的公钥(hex)，协商密钥
func GetECDHSecret(privateKey, publicKey string) (string, error) {
	e := NewECDH(elliptic.P256())

	privateKeyByte, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	servicePrivateKey := &ecdhPrivateKey{
		D: privateKeyByte,
	}

	publicKeyByte, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	clientPublicKey, ok := e.Unmarshal(publicKeyByte)
	if !ok {
		return "", errors.New("unmarshal public key error")
	}

	secret, err := e.GenerateSharedSecret(servicePrivateKey, clientPublicKey)
	if err != nil {
		return "", err
	}

	if len(secret) == 0 {
		return "", errors.New("secret is empty")
	}

	return fmt.Sprintf("%064x", secret), nil
}
