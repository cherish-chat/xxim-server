package utils

import (
	"testing"
)

func Test_xAes(t *testing.T) {
	type args struct {
		key  string
		iv   string
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				key:  "test1",
				iv:   "test1test1test1test1test1",
				data: []byte("test1test1test1test1test1test1test1test1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aes.Encrypt(tt.args.key, tt.args.iv, tt.args.data)
			t.Logf("Encrypt() = %v", got)
			decrypted, err := Aes.Decrypt(tt.args.key, tt.args.iv, got)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}
			if string(decrypted) != string(tt.args.data) {
				t.Errorf("Decrypt() got = %v, want %v", string(decrypted), string(tt.args.data))
			}
			t.Logf("Decrypt() = %v", string(decrypted))
		})
	}
}
