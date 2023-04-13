package middleware

import (
	"encoding/base64"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xaes"
	"testing"
)

func TestDecode(t *testing.T) {
	//decodeString, err := base64.StdEncoding.DecodeString("Ick+T0qKs1TrPYxsRaeE6Vgd2VjJpaW1cShoAAxmtGh7xJCIh7uu3jcWtVEvcK64rSN0ghao+A6P+4 OApMjETcRE9LwbvkIyLZK7QF5QOAQDwsMNU2kJ1ToJ018dBu4vQ6INMmP+KOKxVYmLL0P1gVk5zfxAqcx6RvivAsK+YzCQZcy52hYP8s3oeBOMUYNMlcwo7rp31B/Lwk3SRZd3UfAlKojcp4bTd8Zr27sjVaWQ7GZqELUNuBKwvVk7/p9lwWz7YGQjGakVuSWDyO+I0oHQ8BEd0+RIeenQWEQXnTJ9+xTbTA2fjr4fv4W8Om4Tye4wsRtmjhjWzIWT73v4ZSFI+N4ll8JyE89ntwTn6Hld6f1VX08hHISgwSU/Q4UDsd1HZnnEWpeLPemwvPwwAUUtQnQQImww1Xm0dLfYJKx6D4IPklpELi8oETIIOOoM5wNeA8UbUtl/CTstqwmM6GtouKnfngvIzMLTmpNipJo=")
	decodeString, err := base64.StdEncoding.DecodeString("Al91cjdoHzbeXZenD0CcxWb8DftV76vCpM3SM6FJVnXvQxVYQRroN0YbJUn9dLJJFo08cxzR5jOIMevLHIhQYg==")
	//decodeString, err := base64.StdEncoding.DecodeString("k1+0/KyjL2XZCJWBLE2mAmo6a7kPaNLfIgkLmgdTeIw1NCpIZtohvY2tFutBreR+6iRhlnCK8ANmRY+QulL4sw==")
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	t.Logf("decodeString: %v", string(decodeString))
	decrypted, err := xaes.Decrypt([]byte("1234567890polkjhgfdsasdfghjkl"), []byte("34rtghjmhgfdsdert6y78"), decodeString)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	t.Logf("decrypted: %v", string(decrypted))
}

func TestEncode(t *testing.T) {
	iv := "asdfghjkl"
	key := "lkjhgfdsa"
	t.Logf(utils.Md516(iv))
	t.Logf(utils.Md5(key))
	//encrypt := xaes.Encrypt([]byte(iv), []byte(key), []byte("{}"))
	encrypt := xaes.Encrypt([]byte(iv), []byte(key), []byte("{id: 'test02', password: '123456', terminal: 1}"))
	t.Logf("encrypt: %v", string(encrypt))
	encodeString := base64.StdEncoding.EncodeToString(encrypt)
	t.Logf("encodeString: %v", encodeString)
}
