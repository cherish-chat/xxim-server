package pwd

import "testing"

func TestVerifyPwd(t *testing.T) {
	t.Logf("pwd: %v", VerifyPwd("123456", "98d2b307a143dafc1d15801425b957b712fba45d0d0edaf09e69335325fa33fb", []byte("869ca10190a31b1175f42543546eec6f")))
}
