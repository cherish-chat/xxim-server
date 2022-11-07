package xpwd

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
)

func GeneratePwd(original string, salt string) string {
	dk := pbkdf2.Key([]byte(original), []byte(salt), 1000, 32, sha256.New)
	return hex.EncodeToString(dk)
}

func VerifyPwd(input, db string, salt string) bool {
	return GeneratePwd(input, salt) == db
}
