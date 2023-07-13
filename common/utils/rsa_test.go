package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
	"time"
)

const testPublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAu2PcvRRtbgJteugUfH9cE7JXEZMkcvUTtW0hzguyOcrO5U7Z90iT
Wu+fimnoJBzkutmocUshyDmJUPYMq9R+GdeHza85ycLC72qkqTpwqG8xE6TmqSDD
K0PQLm7fB6hngX0xcrYCJItWv2lSGhismzvyd9scE7JDZg1s2wUMzh58rRYMNEY3
5Sz2FN5LTwPxFqj85+O254GZh6eY4gyWJMpWF0VKtiXlGXeDt+BL5d7bLjYICFHY
u2D6mG5U+tMuvM2wJRS6Q6xQtq96nF92NeS6MjyHNbplxYYOq9I8n7tY7Z/1as/4
Xnq93hf2PqgaV0/zH77rjW4F6ZCvgpVxLQIDAQAB
-----END RSA PUBLIC KEY-----
`

const testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAu2PcvRRtbgJteugUfH9cE7JXEZMkcvUTtW0hzguyOcrO5U7Z
90iTWu+fimnoJBzkutmocUshyDmJUPYMq9R+GdeHza85ycLC72qkqTpwqG8xE6Tm
qSDDK0PQLm7fB6hngX0xcrYCJItWv2lSGhismzvyd9scE7JDZg1s2wUMzh58rRYM
NEY35Sz2FN5LTwPxFqj85+O254GZh6eY4gyWJMpWF0VKtiXlGXeDt+BL5d7bLjYI
CFHYu2D6mG5U+tMuvM2wJRS6Q6xQtq96nF92NeS6MjyHNbplxYYOq9I8n7tY7Z/1
as/4Xnq93hf2PqgaV0/zH77rjW4F6ZCvgpVxLQIDAQABAoIBACpMM0I2vzCquZ2Z
jy4+7UDA66ha50pPiYBVPuEsgLFM1wCpmMeZiTFoj0GGAFFOeE643K2eAOUaH6W3
tEqA72nT3aKO3+Nr4+Z40uwj9dP/LTu66Bna/FLivrYMbqli2OJAqQ20ia1ICm+w
TUj4stVjZaqqOZ80iMQbWaviau6HTkfty2+xI6vpvbGYPmMx5IXVjlJsB8NbAoXH
7sVcHOgMgqrMnOWZi4kKyP235T/HOITBxFfWH98Mx0QL0BUlVOgLnPy4GAOBgHgE
ETpQi3czSTS2YsUM0p+fXpVB6aXsUPbASsg70LqfMKVJmUZWkJb417pgieAvW+OX
RrmHDSECgYEA2ALM9q0+6vmGHUPzdDi1MQ65hHA7yUeI56Z3b0q31VTOl/3YDkus
IJ9B4OLppv/QFYql0eQfhTpk5ClI4gbn/0KpHws7Q4+0Wlr+IiAWSI/YY50X+U5C
UymKD3R2ClELhvLnxxJIng1Zh20kVzAGVttiUnTlEYIOilcx41YVbQkCgYEA3hSp
NbCe9Z0LCQGL/TiIc5SBqiJwC62GjMwwTwSQtkmz8nFrlHZopmhYPxd4ayN76Z4O
JEn70ahlNoWFWpHx9+KosY9wVmaNYNjfEADm1caapsa2dz5mykgw+5arj6OAZ7j4
tDZA5x/CM3BpM8yCqrpb4tKLjq8PvN3ACBa70AUCgYBPvLjsVgdjtbhMFUlJHaXE
9iqFOOjY5A8lc82ix3IUzbl1Yb7fiA+B+0fWO+0EOGoXiZasZAk+pM+ZaaP9y47Y
K0NCsmKuDd4FfJFTB4UyQ+cc3mB7JuhUyoCsM9Fe/YvDxObKFXW44jSqSR+hD5lH
drRUu9HTJK85YfaIdL50AQKBgGuYhwLEN7+3/oi2fySIJ1QYN1o+pRqDUBUXOLCP
/azTuKNV4FFlrP4yv86RiH4gCwD82s0qKx9A/wiTWDCxVRJMdn7QiBTUStsJN8mB
JlWci4ER9YWAbjzDDThXn3dQN/4I2DY3supHsMdLRy0ZgJVHBQ24BHV0y6MtrMQ+
f3AhAoGACybBSjGCXQyrd+30wmTkjinSKFqeQDmQ1Lb1hyANgG1SfoTrmgqnY9Vp
DMl0ePDkP5Styp3/vmipti8XA43W609IKzrNlJkgibeJ8kfMTUkN4Je+wqJ3ClOy
e0634LXaHVRocRb6lF1ggqw5aXz+vNPKZz/BwwFctTpKWV2wXC4=
-----END RSA PRIVATE KEY-----
`

func TestParseKey(t *testing.T) {
	{
		block, rest := pem.Decode([]byte(testPrivateKey))
		if len(rest) > 0 {
			logx.Errorf("private key is invalid")
			panic("private key is invalid")
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			logx.Errorf("parse private key failed: %v", err)
			panic(err)
		}
		t.Logf("privateKey: %v", privateKey)
	}

	{
		block, rest := pem.Decode([]byte(testPublicKey))
		if len(rest) > 0 {
			logx.Errorf("public key is invalid")
			panic("public key is invalid")
		}
		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			logx.Errorf("parse public key failed: %v", err)
			panic(err)
		}
		t.Logf("publicKey: %v", publicKey)
	}

}

func TestRsa(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := &key.PublicKey
	privateKey := key

	// 转字符串 --BEGIN RSA PRIVATE KEY-- --END RSA PRIVATE KEY--
	privateKeyString := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	// 转字符串 --BEGIN PUBLIC KEY-- --END PUBLIC KEY--
	pubKeyString := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})
	t.Logf("privateKeyString: \n%s\n", privateKeyString)
	t.Logf("pubKeyString: \n%s\n", pubKeyString)

	rsa := NewRsa(pubKey, privateKey)
	origin := []byte("hello world")
	encrypt, err := rsa.Encrypt(origin)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encrypt: %s", encrypt)
	time.Sleep(time.Second * 20)
	decrypt, err := rsa.Decrypt(encrypt)
	if err != nil {
		t.Fatal(err)
	}
	if string(decrypt) != string(origin) {
		t.Fatal("decrypt not match origin")
	}
}
