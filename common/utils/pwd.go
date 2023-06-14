package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
)

type xPwd struct {
}

var Pwd = &xPwd{}

func (x *xPwd) GeneratePwd(original string, salt string) string {
	dk := pbkdf2.Key([]byte(original), []byte(salt), 1000, 32, sha256.New)
	return hex.EncodeToString(dk)
}

func (x *xPwd) VerifyPwd(input, db string, salt string) bool {
	return x.GeneratePwd(input, salt) == db
}
