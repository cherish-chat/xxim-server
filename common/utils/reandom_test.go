package utils

import "testing"

func Test_xRandom_Int(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				length: 6,
			},
		},
	}
	for _, tt := range tests {
		for i := 0; i < 10000; i++ {
			t.Run(tt.name, func(t *testing.T) {
				x := &xRandom{}
				if got := x.Int(tt.args.length); len(Number.Int64ToString(int64(got))) != tt.args.length {
					t.Fatalf("Int() = %v, wantlen %v", got, tt.args.length)
				}
			})
		}
	}
}
