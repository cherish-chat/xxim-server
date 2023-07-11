package utils

import "testing"

func Test_xString_FirstUtf8(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "test1",
		args: args{
			s: "test",
		},
		want: "t",
	}, {
		name: "test2",
		args: args{
			s: "你好",
		},
		want: "你",
	}, {
		name: "test3",
		args: args{
			s: "",
		},
		want: "",
	}, {
		name: "test4",
		args: args{
			//日语
			s: "こんにちは",
		},
		want: "こ",
	}, {
		name: "test5",
		args: args{
			//韩语
			s: "안녕하세요",
		},
		want: "안",
	}, {
		name: "test6",
		args: args{
			//数字
			s: "1办洒",
		},
		want: "1",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xString{}
			if got := x.FirstUtf8(tt.args.s); got != tt.want {
				t.Errorf("FirstUtf8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_xString_Utf8Split(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				s:      "赵钱孙李",
				length: 2,
			},
			want: "赵钱",
		},
		{
			name: "test1",
			args: args{
				s:      "赵钱孙李",
				length: 5,
			},
			want: "赵钱孙李",
		},
		{
			name: "test1",
			args: args{
				s:      "a赵钱孙李",
				length: 3,
			},
			want: "a赵钱",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xString{}
			if got := x.Utf8Split(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("Utf8Split() = %v, want %v", got, tt.want)
			}
		})
	}
}
