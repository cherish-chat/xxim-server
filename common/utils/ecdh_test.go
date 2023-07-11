package utils

import (
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func Test_1(t *testing.T) {
	e := NewECDH(elliptic.P256())
	testECDH(e, t)
}

func testECDH(e ECDH, t testing.TB) {
	var privKey1, privKey2 crypto.PrivateKey
	var pubKey1, pubKey2 crypto.PublicKey
	var pubKey1Buf, pubKey2Buf []byte
	var err error
	var ok bool
	var secret1, secret2 []byte

	privKey1, pubKey1, err = e.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	privKey2, pubKey2, err = e.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	pubKey1Buf = e.Marshal(pubKey1)
	pubKey2Buf = e.Marshal(pubKey2)

	pubKey1, ok = e.Unmarshal(pubKey1Buf)
	if !ok {
		t.Fatalf("Unmarshal does not work")
	}

	pubKey2, ok = e.Unmarshal(pubKey2Buf)
	if !ok {
		t.Fatalf("Unmarshal does not work")
	}

	secret1, err = e.GenerateSharedSecret(privKey1, pubKey2)
	if err != nil {
		t.Error(err)
	}
	secret2, err = e.GenerateSharedSecret(privKey2, pubKey1)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Fatalf("The two shared keys: %d, %d do not match", secret1, secret2)
	}
	t.Logf("The shared key is: %d, len: %d", secret1, len(secret1))

	// iv 取中间16位
	iv := secret1[8:24]
	// 模拟服务端加密
	encryptString := Aes.Encrypt(secret1, iv, []byte("hello world"))
	fmt.Println(string(encryptString))

	// 模拟客户端解密
	decryptString, err := Aes.Decrypt(secret1, iv, encryptString)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decryptString))
}

func TestGetECDHSecret(t *testing.T) {
	publicKey := "04421006ef7fdfa44baa1df88bdff6be915f3e20938041604deec9e3f8ba532cf1248dfb90f874a905464a1e06e7154bb0850586d2c3e996e5260d97efc805d34a"
	// 等待用户输入 publicKey
	fmt.Println("请输入公钥:")
	fmt.Scanln(&publicKey)
	e := NewECDH(elliptic.P256())
	privKey, pubKey, err := e.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	fmt.Println(e.HexEncodePublicKeyToString(pubKey))

	privateKey := e.HexEncodePrivateKeyToString(privKey)
	fmt.Println(privateKey)

	secret, err := GetECDHSecret(privateKey, publicKey)
	if err != nil {
		return
	}
	fmt.Println(secret)

	//keyByte, err := hex.DecodeString(secret)
	//encryptString := Aes.Encrypt("1", "hello world", keyByte)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(encryptString))
}
