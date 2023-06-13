package messagemodel

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"reflect"
	"testing"
)

func TestGenerateMessageId(t *testing.T) {
	type args struct {
		senderUserId string
		targetId     string
		targetType   TargetType
		seq          int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				senderUserId: "1234567890abcdef1234567890abcdef",
				targetId:     "1234567890abcdef1234567890abcdef",
				targetType:   1,
				seq:          1234567890123456789,
			},
			want: utils.Md5("senderUserId=1234567890abcdef1234567890abcdef&seq=1234567890123456789&targetId=1234567890abcdef1234567890abcdef&targetType=1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMessageId(tt.args.senderUserId, tt.args.targetId, tt.args.targetType, tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateMessageId() = %v, want %v", got, tt.want)
			}
		})
	}
}
