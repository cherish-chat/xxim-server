package i18n

import (
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewI18N(t *testing.T) {
	type args struct {
		mysql *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *I18N
	}{
		{
			name: "init",
			args: args{nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewI18N(tt.args.mysql); !reflect.DeepEqual(got, tt.want) {
				t.Logf(got.T("en", "登录"))
			}
		})
	}
}
