package pwd

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
)

func GeneratePwd(original string, salt []byte) string {
	dk := pbkdf2.Key([]byte(original), salt, 1000, 32, sha256.New)
	return hex.EncodeToString(dk)
}

func VerifyPwd(original, encode string, salt []byte) bool {
	pwd := GeneratePwd(original, salt)
	return pwd == encode
}
