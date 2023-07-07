package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
)

func main() {
	e := utils.NewECDH(elliptic.P256())
	privKey, pubKey, err := e.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	fmt.Println(e.HexEncodePublicKeyToString(pubKey))

	publicKey := "0475de22a325d0e73f1beebdbbb4e2219a9e5bbe1efcbfa83c51ae1142ed202a2d89b5692c63ca30d2ceb795c7b4744fca9df11c1977246dd4e5076c74ae5fc930"
	// 等待用户输入 publicKey
	fmt.Println("请输入公钥:")
	fmt.Scanln(&publicKey)
	privateKey := e.HexEncodePrivateKeyToString(privKey)
	fmt.Println(privateKey)

	secret, err := utils.GetECDHSecret(privateKey, publicKey)
	if err != nil {
		return
	}
	fmt.Println(secret)

}
